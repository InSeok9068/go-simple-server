package wardrobe

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode"

	"simple-server/projects/closet/db"
	"simple-server/projects/closet/views"

	"github.com/labstack/echo/v4"
)

const (
	maxUploadSize = 20 << 20 // 20MB
	itemsPerKind  = 12
)

var (
	kindOrder = []string{"top", "bottom", "shoes", "accessory"}

	remoteFetchClient = &http.Client{
		Timeout: 15 * time.Second,
	}
)

func renderItemsHTML(ctx context.Context, queries *db.Queries, uid string, kind string, tags []string) (string, error) {
	groups, err := loadGroupedItems(ctx, queries, uid, kind, tags)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if err := views.ItemsSection(groups).Render(ctx, &builder); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func loadGroupedItems(ctx context.Context, queries *db.Queries, uid string, kind string, tags []string) (map[string][]views.ClosetItem, error) {
	targetKinds := determineTargetKinds(kind)
	groups := initializeGroups(targetKinds)

	if uid == "" {
		ensureAllKinds(groups, kind)
		return groups, nil
	}

	tagJSON := tagsToJSON(tags)
	for _, k := range targetKinds {
		items, err := fetchClosetItems(ctx, queries, uid, k, tagJSON)
		if err != nil {
			return nil, err
		}
		groups[k] = items
	}

	ensureAllKinds(groups, kind)
	return groups, nil
}

func determineTargetKinds(kind string) []string {
	if kind == "" {
		return kindOrder
	}
	return []string{kind}
}

func initializeGroups(kinds []string) map[string][]views.ClosetItem {
	groups := make(map[string][]views.ClosetItem, len(kinds))
	for _, k := range kinds {
		groups[k] = []views.ClosetItem{}
	}
	return groups
}

func ensureAllKinds(groups map[string][]views.ClosetItem, requestedKind string) {
	if requestedKind != "" {
		return
	}
	for _, k := range kindOrder {
		if _, ok := groups[k]; !ok {
			groups[k] = []views.ClosetItem{}
		}
	}
}

func fetchClosetItems(ctx context.Context, queries *db.Queries, uid, kind, tagJSON string) ([]views.ClosetItem, error) {
	rows, err := queries.ListItems(ctx, db.ListItemsParams{
		UserUid:    uid,
		KindFilter: kind,
		TagJson:    tagJSON,
		Offset:     0,
		Limit:      int64(itemsPerKind),
	})
	if err != nil {
		return nil, err
	}

	items := make([]views.ClosetItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, views.NewClosetItem(row))
	}
	return items, nil
}

func buildMetadataTags(meta *ImageMetadata) []string {
	if meta == nil {
		return nil
	}
	extras := make([]string, 0, len(meta.Tags)+4)
	extras = append(extras, meta.Tags...)
	if meta.Season != "" {
		extras = append(extras, fmt.Sprintf("season:%s", meta.Season))
	}
	if meta.Style != "" {
		extras = append(extras, fmt.Sprintf("style:%s", meta.Style))
	}
	for _, color := range meta.Colors {
		if color == "" {
			continue
		}
		extras = append(extras, fmt.Sprintf("color:%s", color))
	}
	return extras
}

func mergeTags(base []string, extras ...string) []string {
	if len(extras) == 0 {
		return base
	}
	seen := make(map[string]struct{}, len(base)+len(extras))
	merged := make([]string, 0, len(base)+len(extras))
	for _, tag := range base {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		merged = append(merged, tag)
	}
	for _, tag := range extras {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		merged = append(merged, tag)
	}
	return merged
}

func buildEmbeddingContext(kind string, tags []string, meta *ImageMetadata) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("종류: %s\n", kind))
	if meta != nil {
		if meta.Summary != "" {
			builder.WriteString("요약: ")
			builder.WriteString(meta.Summary)
			builder.WriteString("\n")
		}
		if meta.Season != "" {
			builder.WriteString("계절: ")
			builder.WriteString(meta.Season)
			builder.WriteString("\n")
		}
		if meta.Style != "" {
			builder.WriteString("스타일: ")
			builder.WriteString(meta.Style)
			builder.WriteString("\n")
		}
		if len(meta.Colors) > 0 {
			builder.WriteString("색상: ")
			builder.WriteString(strings.Join(meta.Colors, ", "))
			builder.WriteString("\n")
		}
	}
	if len(tags) > 0 {
		builder.WriteString("태그: ")
		builder.WriteString(strings.Join(tags, ", "))
	}
	return builder.String()
}

func toNullString(value string) sql.NullString {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: trimmed, Valid: true}
}

func parseTags(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	splitter := func(r rune) bool {
		if r == ',' {
			return true
		}
		if r == '\n' || r == '\r' || r == '\t' {
			return true
		}
		return unicode.IsSpace(r)
	}
	raw := strings.FieldsFunc(input, splitter)
	seen := make(map[string]struct{}, len(raw))
	tags := make([]string, 0, len(raw))
	for _, part := range raw {
		tag := strings.TrimSpace(strings.TrimPrefix(part, "#"))
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		tags = append(tags, tag)
	}
	return tags
}

func tagsToJSON(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	b, err := json.Marshal(tags)
	if err != nil {
		slog.Warn("태그 JSON 직렬화 실패", "error", err)
		return "[]"
	}
	return string(b)
}

func normalizeKind(kind string) string {
	lower := strings.ToLower(kind)
	for _, k := range kindOrder {
		if lower == k {
			return k
		}
	}
	return ""
}

func imageDimensions(data []byte) (int, int) {
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0
	}
	return config.Width, config.Height
}

func nullableInt(v int) sql.NullInt64 {
	if v == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(v), Valid: true}
}

func convertIDsRow(row db.ListItemsByIDsRow) db.ListItemsRow {
	return db.ListItemsRow(row)
}

func fetchItemForUpdate(ctx context.Context, queries *db.Queries, uid, rawID string) (int64, db.GetItemDetailRow, error) {
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		return 0, db.GetItemDetailRow{}, echo.NewHTTPError(http.StatusBadRequest, "잘못된 요청입니다.")
	}

	row, err := queries.GetItemDetail(ctx, db.GetItemDetailParams{
		ID:      id,
		UserUid: uid,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return 0, db.GetItemDetailRow{}, echo.NewHTTPError(http.StatusNotFound, "아이템을 찾지 못했어요.")
	}
	if err != nil {
		return 0, db.GetItemDetailRow{}, err
	}

	return id, row, nil
}

func updateItemMetadataAndTags(ctx context.Context, itemID int64, uid string, meta *ImageMetadata, tags []string) error {
	colorText := strings.Join(meta.Colors, ",")
	database, err := db.GetDB()
	if err != nil {
		return err
	}

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			slog.Warn("트랜잭션 롤백 실패", "error", rollbackErr)
		}
	}()

	qtx := db.New(tx)
	if err := qtx.UpdateItemMetadata(ctx, db.UpdateItemMetadataParams{
		MetaSummary: toNullString(meta.Summary),
		MetaSeason:  toNullString(meta.Season),
		MetaStyle:   toNullString(meta.Style),
		MetaColors:  toNullString(colorText),
		ID:          itemID,
		UserUid:     uid,
	}); err != nil {
		return err
	}

	if err := persistItemTags(ctx, qtx, itemID, tags); err != nil {
		return err
	}

	return tx.Commit()
}

func parseUpdateItemForm(c echo.Context) (*ImageMetadata, []string) {
	summary := strings.TrimSpace(c.FormValue("meta_summary"))
	season := strings.TrimSpace(c.FormValue("meta_season"))
	style := strings.TrimSpace(c.FormValue("meta_style"))
	colors := parseTags(c.FormValue("meta_colors"))
	tags := parseTags(c.FormValue("tags"))

	meta := &ImageMetadata{
		Summary: summary,
		Season:  season,
		Style:   style,
		Colors:  colors,
	}

	return meta, mergeTags(tags, buildMetadataTags(meta)...)
}

func persistItemTags(ctx context.Context, qtx *db.Queries, itemID int64, tags []string) error {
	if err := qtx.DeleteItemTags(ctx, itemID); err != nil {
		return err
	}

	for _, tag := range tags {
		tagID, err := qtx.UpsertTag(ctx, tag)
		if err != nil {
			return err
		}
		if err := qtx.AttachTag(ctx, db.AttachTagParams{ItemID: itemID, TagID: tagID}); err != nil {
			return err
		}
	}
	return nil
}

func readUploadedFile(fileHeader *multipart.FileHeader) ([]byte, string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, "", echo.NewHTTPError(http.StatusBadRequest, "이미지를 열 수 없어요.")
	}
	defer src.Close()

	limited := io.LimitReader(src, maxUploadSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, "", echo.NewHTTPError(http.StatusInternalServerError, "이미지를 읽는 중 문제가 발생했어요.")
	}
	if len(data) == 0 {
		return nil, "", echo.NewHTTPError(http.StatusBadRequest, "빈 파일은 업로드할 수 없어요.")
	}
	if len(data) > maxUploadSize {
		return nil, "", echo.NewHTTPError(http.StatusRequestEntityTooLarge, "이미지는 20MB 이하로 올려 주세요.")
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}
	mimeType = normalizeMimeType(mimeType)
	if !strings.HasPrefix(mimeType, "image/") {
		return nil, "", echo.NewHTTPError(http.StatusBadRequest, "이미지 파일만 업로드할 수 있어요.")
	}

	return data, mimeType, nil
}

func downloadImageFromURL(ctx context.Context, rawURL string) (string, string, []byte, error) {
	parsed, err := validateImageURL(rawURL)
	if err != nil {
		return "", "", nil, err
	}

	mimeType, data, err := fetchRemoteImageData(ctx, parsed)
	if err != nil {
		return "", "", nil, err
	}

	filename := sanitizeRemoteFilename(parsed)
	return filename, mimeType, data, nil
}

func validateImageURL(rawURL string) (*url.URL, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "올바른 이미지 URL을 입력해 주세요.")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "이미지 URL은 http 또는 https만 허용돼요.")
	}
	return parsed, nil
}

func fetchRemoteImageData(ctx context.Context, parsed *url.URL) (string, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := remoteFetchClient.Do(req)
	if err != nil {
		return "", nil, echo.NewHTTPError(http.StatusBadGateway, "이미지 URL에 접근하지 못했어요.")
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", nil, echo.NewHTTPError(http.StatusBadRequest, "이미지 URL 응답이 정상이 아니에요.")
	}

	limited := io.LimitReader(resp.Body, maxUploadSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return "", nil, echo.NewHTTPError(http.StatusBadGateway, "이미지를 다운로드하지 못했어요.")
	}
	if len(data) == 0 {
		return "", nil, echo.NewHTTPError(http.StatusBadRequest, "이미지 URL에서 데이터를 찾지 못했어요.")
	}
	if len(data) > maxUploadSize {
		return "", nil, echo.NewHTTPError(http.StatusRequestEntityTooLarge, "이미지는 20MB 이하여야 해요.")
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}
	mimeType = normalizeMimeType(mimeType)
	if !strings.HasPrefix(mimeType, "image/") {
		return "", nil, echo.NewHTTPError(http.StatusBadRequest, "이미지 URL만 허용돼요.")
	}

	return mimeType, data, nil
}

func sanitizeRemoteFilename(u *url.URL) string {
	if u == nil {
		return "remote-image"
	}
	name := strings.TrimSpace(path.Base(u.Path))
	if name == "" || name == "." || name == "/" {
		name = strings.TrimSpace(u.Host)
	}
	if name == "" {
		name = "remote-image"
	}
	return name
}

func normalizeMimeType(mimeType string) string {
	if mimeType == "" {
		return ""
	}
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpg":
		return "image/jpeg"
	case "image/x-png":
		return "image/png"
	}
	return mimeType
}
