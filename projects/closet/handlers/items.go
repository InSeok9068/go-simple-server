package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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

	"simple-server/pkg/util/authutil"
	"simple-server/projects/closet/db"
	"simple-server/projects/closet/services"
	"simple-server/projects/closet/views"

	"github.com/labstack/echo/v4"
	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
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

// UploadItem은 옷장 아이템을 업로드한다.
// nolint:cyclop // 업로드 파이프라인의 단계가 많아 임시로 순환 복잡도 검사를 무시한다.
func UploadItem(c echo.Context) error {
	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	kind := normalizeKind(strings.TrimSpace(c.FormValue("kind")))
	if kind == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "의상 종류를 선택해주세요.")
	}

	fileHeader, fileErr := c.FormFile("image")
	if fileErr != nil && !errors.Is(fileErr, http.ErrMissingFile) {
		return echo.NewHTTPError(http.StatusBadRequest, "이미지를 불러오지 못했어요.")
	}

	imageURL := strings.TrimSpace(c.FormValue("image_url"))

	var (
		imageBytes []byte
		filename   string
		mimeType   string
	)

	switch {
	case fileHeader != nil:
		imageBytes, mimeType, err = readUploadedFile(fileHeader)
		if err != nil {
			return err
		}
		filename = fileHeader.Filename
	case imageURL != "":
		filename, mimeType, imageBytes, err = downloadImageFromURL(c.Request().Context(), imageURL)
		if err != nil {
			return err
		}
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "이미지 파일이나 URL 중 하나를 선택해 주세요.")
	}
	mimeType = normalizeMimeType(mimeType)

	sha := sha256.Sum256(imageBytes)
	shaString := hex.EncodeToString(sha[:])

	width, height := imageDimensions(imageBytes)
	ctx := c.Request().Context()

	tags := parseTags(c.FormValue("tags"))
	var aiMetadata *services.ImageMetadata
	if metadata, metaErr := services.AnalyzeClosetImage(ctx, imageBytes, mimeType); metaErr != nil {
		slog.Warn("이미지 메타데이터 생성 실패", "error", metaErr)
	} else {
		aiMetadata = metadata
		tags = mergeTags(tags, buildMetadataTags(metadata)...)
	}

	metaSummary := sql.NullString{}
	metaSeason := sql.NullString{}
	metaStyle := sql.NullString{}
	metaColors := sql.NullString{}
	if aiMetadata != nil {
		metaSummary = toNullString(aiMetadata.Summary)
		metaSeason = toNullString(aiMetadata.Season)
		metaStyle = toNullString(aiMetadata.Style)
		metaColors = toNullString(strings.Join(aiMetadata.Colors, ","))
	}

	embeddingText := buildEmbeddingContext(kind, tags, aiMetadata)

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

	itemID, err := qtx.InsertItem(ctx, db.InsertItemParams{
		UserUid:     uid,
		Kind:        kind,
		Filename:    filename,
		MimeType:    mimeType,
		Bytes:       imageBytes,
		ThumbBytes:  nil,
		Sha256:      sql.NullString{String: shaString, Valid: true},
		Width:       nullableInt(width),
		Height:      nullableInt(height),
		MetaSummary: metaSummary,
		MetaSeason:  metaSeason,
		MetaStyle:   metaStyle,
		MetaColors:  metaColors,
	})
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == int(sqlite3.SQLITE_CONSTRAINT_UNIQUE) {
			existingID, lookupErr := qtx.GetItemIDBySha(ctx, db.GetItemIDByShaParams{
				UserUid: uid,
				Sha256:  sql.NullString{String: shaString, Valid: true},
			})
			if lookupErr != nil {
				return lookupErr
			}
			itemID = existingID
		} else {
			return err
		}
	}

	for _, tag := range tags {
		tagID, tagErr := qtx.UpsertTag(ctx, tag)
		if tagErr != nil {
			return tagErr
		}
		if attachErr := qtx.AttachTag(ctx, db.AttachTagParams{ItemID: itemID, TagID: tagID}); attachErr != nil {
			return attachErr
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go services.EnqueueEmbeddingJob(itemID, uid, embeddingText)

	html, err := renderItemsHTML(ctx, queries, uid, "", nil)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusCreated, html)
}

// ListItems는 필터에 맞는 아이템 목록을 반환한다.
func ListItems(c echo.Context) error {
	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	kind := normalizeKind(strings.TrimSpace(c.QueryParam("kind")))
	tags := parseTags(c.QueryParam("tags"))

	html, err := renderItemsHTML(c.Request().Context(), queries, uid, kind, tags)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, html)
}

// GetItemDetail은 아이템 상세 정보를 반환한다.
func GetItemDetail(c echo.Context) error {
	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "잘못된 요청입니다.")
	}

	row, err := queries.GetItemDetail(c.Request().Context(), db.GetItemDetailParams{
		ID:      id,
		UserUid: uid,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "아이템을 찾지 못했어요.")
	}
	if err != nil {
		return err
	}

	detail := views.NewClosetItemDetail(row)

	var builder strings.Builder
	if err := views.ItemDetailContent(detail).Render(&builder); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}

// ItemImage는 아이템 이미지를 반환한다.
func ItemImage(c echo.Context) error {
	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "잘못된 요청입니다.")
	}

	row, err := queries.GetItemContent(c.Request().Context(), db.GetItemContentParams{
		ID:      id,
		UserUid: uid,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "이미지를 찾을 수 없습니다.")
	}
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentType, row.MimeType)
	if _, err := c.Response().Writer.Write(row.Bytes); err != nil {
		return err
	}

	return nil
}

func renderItemsHTML(ctx context.Context, queries *db.Queries, uid string, kind string, tags []string) (string, error) {
	groups, err := loadGroupedItems(ctx, queries, uid, kind, tags)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if err := views.ItemsSection(groups).Render(&builder); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func loadGroupedItems(ctx context.Context, queries *db.Queries, uid string, kind string, tags []string) (map[string][]views.ClosetItem, error) {
	targetKinds := kindOrder
	if kind != "" {
		targetKinds = []string{kind}
	}

	tagJSON := tagsToJSON(tags)
	groups := make(map[string][]views.ClosetItem, len(targetKinds))

	if uid == "" {
		for _, k := range targetKinds {
			groups[k] = []views.ClosetItem{}
		}
		if kind == "" {
			for _, k := range kindOrder {
				if _, ok := groups[k]; !ok {
					groups[k] = []views.ClosetItem{}
				}
			}
		}
		return groups, nil
	}

	for _, k := range targetKinds {
		rows, err := queries.ListItems(ctx, db.ListItemsParams{
			UserUid:    uid,
			KindFilter: k,
			TagJson:    tagJSON,
			Offset:     0,
			Limit:      int64(itemsPerKind),
		})
		if err != nil {
			return nil, err
		}

		items := make([]views.ClosetItem, 0, len(rows))
		for _, row := range rows {
			item := views.NewClosetItem(row)
			items = append(items, item)
		}
		groups[k] = items
	}

	if kind == "" {
		for _, k := range kindOrder {
			if _, ok := groups[k]; !ok {
				groups[k] = []views.ClosetItem{}
			}
		}
	}

	return groups, nil
}

func buildMetadataTags(meta *services.ImageMetadata) []string {
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

func buildEmbeddingContext(kind string, tags []string, meta *services.ImageMetadata) string {
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

// DeleteItem은 업로드된 항목을 제거한다.
func DeleteItem(c echo.Context) error {
	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "잘못된 요청입니다.")
	}

	if err := queries.DeleteItem(c.Request().Context(), db.DeleteItemParams{
		ID:      id,
		UserUid: uid,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateItem은 아이템의 메타데이터와 태그를 수정한다.
func UpdateItem(c echo.Context) error {
	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "잘못된 요청입니다.")
	}

	ctx := c.Request().Context()
	row, err := queries.GetItemDetail(ctx, db.GetItemDetailParams{
		ID:      id,
		UserUid: uid,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "아이템을 찾지 못했어요.")
	}
	if err != nil {
		return err
	}

	summary := strings.TrimSpace(c.FormValue("meta_summary"))
	season := strings.TrimSpace(c.FormValue("meta_season"))
	style := strings.TrimSpace(c.FormValue("meta_style"))
	colors := parseTags(c.FormValue("meta_colors"))
	tags := parseTags(c.FormValue("tags"))

	colorText := strings.Join(colors, ",")

	meta := &services.ImageMetadata{
		Summary: summary,
		Season:  season,
		Style:   style,
		Colors:  colors,
	}

	mergedTags := mergeTags(tags, buildMetadataTags(meta)...)

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
		MetaSummary: toNullString(summary),
		MetaSeason:  toNullString(season),
		MetaStyle:   toNullString(style),
		MetaColors:  toNullString(colorText),
		ID:          id,
		UserUid:     uid,
	}); err != nil {
		return err
	}

	if err := qtx.DeleteItemTags(ctx, id); err != nil {
		return err
	}

	for _, tag := range mergedTags {
		tagID, tagErr := qtx.UpsertTag(ctx, tag)
		if tagErr != nil {
			return tagErr
		}
		if attachErr := qtx.AttachTag(ctx, db.AttachTagParams{ItemID: id, TagID: tagID}); attachErr != nil {
			return attachErr
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	kind := strings.ToLower(row.Kind)
	go services.EnqueueEmbeddingJob(id, uid, buildEmbeddingContext(kind, mergedTags, meta))

	return c.NoContent(http.StatusNoContent)
}

// RecommendOutfit은 날씨와 스타일 조건에 맞는 아이템을 추천한다.
func RecommendOutfit(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	if err := c.Request().ParseForm(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "입력값을 확인해주세요.")
	}
	weather := strings.TrimSpace(c.FormValue("weather"))
	style := strings.TrimSpace(c.FormValue("style"))
	skipIDs := strings.TrimSpace(c.FormValue("skip_ids"))
	locks := parseLockSelections(c)

	results, cacheToken, hasMore, err := services.RecommendOutfit(c.Request().Context(), uid, weather, style, skipIDs, locks)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	viewResults := make([]views.RecommendationItem, 0, len(results))
	for _, result := range results {
		viewResults = append(viewResults, views.RecommendationItem{
			Kind: result.Kind,
			Item: views.NewClosetItem(convertIDsRow(result.Item)),
		})
	}

	var builder strings.Builder
	if err := views.RecommendationDialog(viewResults, weather, style, cacheToken, hasMore, locks).Render(&builder); err != nil {
		return err
	}
	return c.HTML(http.StatusOK, builder.String())
}

func parseLockSelections(c echo.Context) map[string]int64 {
	locks := make(map[string]int64)
	for _, kind := range kindOrder {
		field := fmt.Sprintf("lock_%s", kind)
		value := strings.TrimSpace(c.FormValue(field))
		if value == "" {
			continue
		}
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			continue
		}
		locks[kind] = id
	}
	if len(locks) == 0 {
		return nil
	}
	return locks
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
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", "", nil, echo.NewHTTPError(http.StatusBadRequest, "올바른 이미지 URL을 입력해 주세요.")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", "", nil, echo.NewHTTPError(http.StatusBadRequest, "이미지 URL은 http 또는 https만 허용돼요.")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return "", "", nil, err
	}

	resp, err := remoteFetchClient.Do(req)
	if err != nil {
		return "", "", nil, echo.NewHTTPError(http.StatusBadGateway, "이미지 URL에 접근하지 못했어요.")
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", "", nil, echo.NewHTTPError(http.StatusBadRequest, "이미지 URL 응답이 정상이 아니에요.")
	}

	limited := io.LimitReader(resp.Body, maxUploadSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return "", "", nil, echo.NewHTTPError(http.StatusBadGateway, "이미지를 다운로드하지 못했어요.")
	}
	if len(data) == 0 {
		return "", "", nil, echo.NewHTTPError(http.StatusBadRequest, "이미지 URL에서 데이터를 찾지 못했어요.")
	}
	if len(data) > maxUploadSize {
		return "", "", nil, echo.NewHTTPError(http.StatusRequestEntityTooLarge, "이미지는 20MB 이하여야 해요.")
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}
	mimeType = normalizeMimeType(mimeType)
	if !strings.HasPrefix(mimeType, "image/") {
		return "", "", nil, echo.NewHTTPError(http.StatusBadRequest, "이미지 URL만 허용돼요.")
	}

	filename := sanitizeRemoteFilename(parsed)

	return filename, mimeType, data, nil
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
