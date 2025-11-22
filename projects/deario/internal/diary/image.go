package diary

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"simple-server/internal/middleware"
	"simple-server/pkg/util/authutil"
	"simple-server/pkg/util/firebaseutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views/components"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
)

// DiaryImagesPage는 일기 이미지 폼을 렌더링한다.
func DiaryImagesPage(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("20060102")
	}
	date = strings.ReplaceAll(date, "-", "")

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: date,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return components.DiaryImages(date, "", "", "").Render(c.Request().Context(), c.Response().Writer)
	}

	return components.DiaryImages(date, diary.ImageUrl1, diary.ImageUrl2, diary.ImageUrl3).Render(c.Request().Context(), c.Response().Writer)
}

// UploadDiaryImage는 새 이미지를 업로드하고 저장한다.
func UploadDiaryImage(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := normalizeDate(c.FormValue("date"))
	url := c.FormValue("url")
	if url == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "URL이 필요합니다.")
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	diary, slot, err := diaryAndSlot(c.Request().Context(), queries, uid, date)
	if err != nil {
		return err
	}

	diary = setDiaryImage(diary, slot, url)
	if err := queries.UpdateDiaryImages(c.Request().Context(), db.UpdateDiaryImagesParams{
		ImageUrl1: diary.ImageUrl1,
		ImageUrl2: diary.ImageUrl2,
		ImageUrl3: diary.ImageUrl3,
		ID:        diary.ID,
	}); err != nil {
		return err
	}

	return components.DiaryImages(date, diary.ImageUrl1, diary.ImageUrl2, diary.ImageUrl3).Render(c.Request().Context(), c.Response().Writer)
}

// DeleteDiaryImage는 지정된 일기 이미지를 삭제한다.
func DeleteDiaryImage(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := normalizeDate(c.FormValue("date"))
	slotStr := c.FormValue("slot")
	slot, err := strconv.Atoi(slotStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "잘못된 슬롯입니다.")
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	diary, err := getDiary(c.Request().Context(), queries, uid, date)
	if err != nil {
		return err
	}

	if err := removeDiaryFirebaseImage(c.Request().Context(), diary, slot); err != nil {
		slog.Error("파이어베이스 이미지 삭제 실패", "error", err)
		return err
	}

	diary, err = clearDiaryImage(diary, slot)
	if err != nil {
		return err
	}

	if err := queries.UpdateDiaryImages(c.Request().Context(), db.UpdateDiaryImagesParams{
		ImageUrl1: diary.ImageUrl1,
		ImageUrl2: diary.ImageUrl2,
		ImageUrl3: diary.ImageUrl3,
		ID:        diary.ID,
	}); err != nil {
		return err
	}

	return components.DiaryImages(date, diary.ImageUrl1, diary.ImageUrl2, diary.ImageUrl3).Render(c.Request().Context(), c.Response().Writer)
}

// firstEmptyImageSlot은 비어 있는 이미지 슬롯을 찾는다.
func firstEmptyImageSlot(d db.Diary) (int, bool) {
	slots := []string{d.ImageUrl1, d.ImageUrl2, d.ImageUrl3}
	for i, v := range slots {
		if v == "" {
			return i, true
		}
	}
	return -1, false
}

// normalizeDate는 날짜 문자열을 YYYYMMDD 형식으로 변환한다.
func normalizeDate(d string) string {
	if d == "" {
		return time.Now().Format("20060102")
	}
	return strings.ReplaceAll(d, "-", "")
}

// diaryAndSlot은 일기와 비어 있는 이미지 슬롯을 반환한다.
func diaryAndSlot(ctx context.Context, q *db.Queries, uid, date string) (db.Diary, int, error) {
	diary, err := q.GetDiary(ctx, db.GetDiaryParams{Uid: uid, Date: date})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Diary{}, -1, echo.NewHTTPError(http.StatusBadRequest, "일기를 먼저 작성해주세요.")
		}
		return db.Diary{}, -1, err
	}
	slot, ok := firstEmptyImageSlot(diary)
	if !ok {
		return db.Diary{}, -1, echo.NewHTTPError(http.StatusBadRequest, "이미지는 최대 3장까지 저장할 수 있습니다.")
	}
	return diary, slot, nil
}

// getDiary는 해당 날짜의 일기를 조회한다.
func getDiary(ctx context.Context, q *db.Queries, uid, date string) (db.Diary, error) {
	diary, err := q.GetDiary(ctx, db.GetDiaryParams{Uid: uid, Date: date})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Diary{}, echo.NewHTTPError(http.StatusBadRequest, "이미지가 없습니다.")
		}
		return db.Diary{}, err
	}
	return diary, nil
}

// removeDiaryFirebaseImage는 파이어베이스에서 이미지를 삭제한다.
func removeDiaryFirebaseImage(ctx context.Context, d db.Diary, slot int) error {
	imageURL, err := diaryImageURL(d, slot)
	if err != nil {
		return err
	}
	if imageURL == "" {
		return nil
	}
	return deleteFirebaseImage(ctx, imageURL)
}

// setDiaryImage는 지정된 슬롯에 이미지 URL을 설정한다.
func setDiaryImage(d db.Diary, slot int, url string) db.Diary {
	switch slot {
	case 0:
		d.ImageUrl1 = url
	case 1:
		d.ImageUrl2 = url
	case 2:
		d.ImageUrl3 = url
	}
	return d
}

// clearDiaryImage는 지정된 슬롯의 이미지 URL을 비운다.
func clearDiaryImage(d db.Diary, slot int) (db.Diary, error) {
	switch slot {
	case 1:
		d.ImageUrl1 = ""
	case 2:
		d.ImageUrl2 = ""
	case 3:
		d.ImageUrl3 = ""
	default:
		return d, echo.NewHTTPError(http.StatusBadRequest, "잘못된 슬롯입니다.")
	}
	return d, nil
}

// diaryImageURL은 지정된 슬롯의 이미지 URL을 반환한다.
func diaryImageURL(d db.Diary, slot int) (string, error) {
	switch slot {
	case 1:
		return d.ImageUrl1, nil
	case 2:
		return d.ImageUrl2, nil
	case 3:
		return d.ImageUrl3, nil
	default:
		return "", echo.NewHTTPError(http.StatusBadRequest, "잘못된 슬롯입니다.")
	}
}

// deleteFirebaseImage는 파이어베이스에서 이미지를 삭제한다.
func deleteFirebaseImage(ctx context.Context, raw string) error {
	bucket, object, err := firebaseutil.ParseFirebaseURL(raw)
	if err != nil {
		return fmt.Errorf("URL 파싱 실패: %w", err)
	}
	client, err := middleware.App.Storage(ctx)
	if err != nil {
		return fmt.Errorf("스토리지 클라이언트 생성 실패: %w", err)
	}
	b, err := client.Bucket(bucket)
	if err != nil {
		return fmt.Errorf("버킷 가져오기 실패: %w", err)
	}
	if err := b.Object(object).Delete(ctx); err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil
		}
		return fmt.Errorf("파일 삭제 실패: %w", err)
	}
	return nil
}
