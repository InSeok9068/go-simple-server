package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
	"time"

	"simple-server/internal/middleware"
	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
)

func DiaryImages(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("20060102")
	}
	date = strings.ReplaceAll(date, "-", "")

	queries, err := db.GetQueries(c.Request().Context())
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
		return views.DiaryImages(date, "", "", "").Render(c.Response().Writer)
	}

	return views.DiaryImages(date, diary.ImageUrl1, diary.ImageUrl2, diary.ImageUrl3).Render(c.Response().Writer)
}

func DiaryImageSave(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := normalizeDate(c.FormValue("date"))
	url := c.FormValue("url")
	if url == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "URL이 필요합니다.")
	}

	queries, err := db.GetQueries(c.Request().Context())
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

	return views.DiaryImages(date, diary.ImageUrl1, diary.ImageUrl2, diary.ImageUrl3).Render(c.Response().Writer)
}

func DiaryImageDelete(c echo.Context) error {
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

	queries, err := db.GetQueries(c.Request().Context())
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

	return views.DiaryImages(date, diary.ImageUrl1, diary.ImageUrl2, diary.ImageUrl3).Render(c.Response().Writer)
}

func firstEmptyImageSlot(d db.Diary) (int, bool) {
	slots := []string{d.ImageUrl1, d.ImageUrl2, d.ImageUrl3}
	for i, v := range slots {
		if v == "" {
			return i, true
		}
	}
	return -1, false
}

func normalizeDate(d string) string {
	if d == "" {
		return time.Now().Format("20060102")
	}
	return strings.ReplaceAll(d, "-", "")
}

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

func parseFirebaseURL(raw string) (string, string, error) {
	u, err := neturl.Parse(raw)
	if err != nil {
		return "", "", err
	}
	p := u.RawPath
	if p == "" {
		p = u.Path
	}
	parts := strings.Split(p, "/")
	if len(parts) < 6 {
		return "", "", fmt.Errorf("잘못된 URL")
	}
	bucket := parts[3]
	object, err := neturl.QueryUnescape(parts[5])
	if err != nil {
		return "", "", err
	}
	return bucket, object, nil
}

func deleteFirebaseImage(ctx context.Context, raw string) error {
	bucket, object, err := parseFirebaseURL(raw)
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
