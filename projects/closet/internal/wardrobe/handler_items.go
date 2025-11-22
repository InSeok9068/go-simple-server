package wardrobe

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"simple-server/pkg/util/authutil"
	"simple-server/projects/closet/db"
	"simple-server/projects/closet/views"

	"github.com/labstack/echo/v4"
	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
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
	var aiMetadata *ImageMetadata
	if metadata, metaErr := AnalyzeClosetImage(ctx, imageBytes, mimeType); metaErr != nil {
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

	go EnqueueEmbeddingJob(itemID, uid, embeddingText)

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
	if err := views.ItemDetailContent(detail).Render(c.Request().Context(), &builder); err != nil {
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

	ctx := c.Request().Context()
	id, row, err := fetchItemForUpdate(ctx, queries, uid, c.Param("id"))
	if err != nil {
		return err
	}

	meta, mergedTags := parseUpdateItemForm(c)
	if err := updateItemMetadataAndTags(ctx, id, uid, meta, mergedTags); err != nil {
		return err
	}

	kind := strings.ToLower(row.Kind)
	go EnqueueEmbeddingJob(id, uid, buildEmbeddingContext(kind, mergedTags, meta))

	return c.NoContent(http.StatusNoContent)
}
