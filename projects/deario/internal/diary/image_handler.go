package diary

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views/components"

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
