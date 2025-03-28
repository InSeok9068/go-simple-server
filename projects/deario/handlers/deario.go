package handlers

import (
	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os"
	"simple-server/internal/connection"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views"
	shared "simple-server/shared/views"
	"strings"
	"time"
)

func Index(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("20060102")
	} else {
		date = strings.ReplaceAll(date, "-", "")
	}
	return views.Index(os.Getenv("APP_TITLE"), date).Render(c.Response().Writer)
}

func Login(c echo.Context) error {
	return shared.Login().Render(c.Response().Writer)
}

func Diary(c echo.Context) error {
	user, ok := c.Get("user").(*auth.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
	}

	date := c.FormValue("date")
	if date == "" {
		date = time.Now().Format("20060102")
	} else {
		date = strings.ReplaceAll(date, "-", "")
	}

	dbCon, err := connection.AppDBOpen()
	if err != nil {
		slog.Error("Failed to open database", "error", err.Error())
	}
	queries := db.New(dbCon)

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  user.UID,
		Date: date,
	})

	if err != nil {
		return views.DiaryContentForm(date, "").Render(c.Response().Writer)
	} else {
		return views.DiaryContentForm(diary.Date, diary.Content).Render(c.Response().Writer)
	}
}

func Save(c echo.Context) error {
	user, ok := c.Get("user").(*auth.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
	}

	date := c.FormValue("date")
	content := c.FormValue("content")

	dbCon, err := connection.AppDBOpen()
	if err != nil {
		slog.Error("Failed to open database", "error", err.Error())
	}
	queries := db.New(dbCon)

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  user.UID,
		Date: date,
	})

	if err != nil {
		_, err = queries.CreateDiary(c.Request().Context(), db.CreateDiaryParams{
			Uid:     user.UID,
			Content: content,
			Date:    date,
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "등록 실패")
		}
	} else {
		_, err = queries.UpdateDiary(c.Request().Context(), db.UpdateDiaryParams{
			Content: content,
			ID:      diary.ID,
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "수정 실패")
		}
	}

	return nil
}
