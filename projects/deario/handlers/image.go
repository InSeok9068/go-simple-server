package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views"

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
		return views.DiaryImages("", "", "").Render(c.Response().Writer)
	}

	return views.DiaryImages(diary.ImageUrl1, diary.ImageUrl2, diary.ImageUrl3).Render(c.Response().Writer)
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

	return views.DiaryImages(diary.ImageUrl1, diary.ImageUrl2, diary.ImageUrl3).Render(c.Response().Writer)
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
