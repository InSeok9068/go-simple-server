package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"unicode"

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

var kindOrder = []string{"top", "bottom", "shoes", "accessory"}

// UploadItem은 옷장 아이템을 업로드한다.
// nolint:cyclop // 업로드 파이프라인의 단계가 많아 임시로 순환 복잡도 검사를 무시한다.
func UploadItem(c echo.Context) error {
	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	kind := normalizeKind(strings.TrimSpace(c.FormValue("kind")))
	if kind == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "의상 종류를 선택해주세요.")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "이미지 파일을 선택해주세요.")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "이미지를 열 수 없습니다.")
	}
	defer src.Close()

	limited := io.LimitReader(src, maxUploadSize+1)
	buf, err := io.ReadAll(limited)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "이미지를 읽을 수 없습니다.")
	}
	if len(buf) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "빈 파일은 업로드할 수 없습니다.")
	}
	if len(buf) > maxUploadSize {
		return echo.NewHTTPError(http.StatusRequestEntityTooLarge, "이미지는 20MB 이하로 업로드해주세요.")
	}

	sha := sha256.Sum256(buf)
	shaString := hex.EncodeToString(sha[:])

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(buf)
	}
	if !strings.HasPrefix(mimeType, "image/") {
		return echo.NewHTTPError(http.StatusBadRequest, "이미지 파일만 업로드할 수 있습니다.")
	}

	width, height := imageDimensions(buf)

	tags := parseTags(c.FormValue("tags"))

	database, err := db.GetDB()
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
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
		Kind:       kind,
		Filename:   file.Filename,
		MimeType:   mimeType,
		Bytes:      buf,
		ThumbBytes: nil,
		Sha256:     sql.NullString{String: shaString, Valid: true},
		Width:      nullableInt(width),
		Height:     nullableInt(height),
	})
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == int(sqlite3.SQLITE_CONSTRAINT_UNIQUE) {
			existingID, lookupErr := qtx.GetItemIDBySha(ctx, sql.NullString{String: shaString, Valid: true})
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

	go services.EnqueueEmbeddingJob(itemID)

	html, err := renderItemsHTML(ctx, queries, "", nil)
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

	kind := normalizeKind(strings.TrimSpace(c.QueryParam("kind")))
	tags := parseTags(c.QueryParam("tags"))

	html, err := renderItemsHTML(c.Request().Context(), queries, kind, tags)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, html)
}

// ItemImage는 아이템 이미지를 반환한다.
func ItemImage(c echo.Context) error {
	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "잘못된 요청입니다.")
	}

	row, err := queries.GetItemContent(c.Request().Context(), id)
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

func renderItemsHTML(ctx context.Context, queries *db.Queries, kind string, tags []string) (string, error) {
	groups, err := loadGroupedItems(ctx, queries, kind, tags)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if err := views.ItemsSection(groups).Render(&builder); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func loadGroupedItems(ctx context.Context, queries *db.Queries, kind string, tags []string) (map[string][]views.ClosetItem, error) {
	targetKinds := kindOrder
	if kind != "" {
		targetKinds = []string{kind}
	}

	tagJSON := tagsToJSON(tags)
	groups := make(map[string][]views.ClosetItem, len(targetKinds))

	for _, k := range targetKinds {
		rows, err := queries.ListItems(ctx, db.ListItemsParams{
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

// DeleteItem은 업로드된 항목을 제거한다.
func DeleteItem(c echo.Context) error {
	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "잘못된 요청입니다.")
	}

	if err := queries.DeleteItem(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
