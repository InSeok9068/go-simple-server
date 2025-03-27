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
	"time"
)

func Index(c echo.Context) error {
	return views.Index(os.Getenv("APP_TITLE")).Render(c.Response().Writer)
}

func Login(c echo.Context) error {
	return shared.Login().Render(c.Response().Writer)
}

func Diary(c echo.Context) error {
	user, ok := c.Get("user").(*auth.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
	}

	dbCon, err := connection.AppDBOpen()
	if err != nil {
		slog.Error("Failed to open database", "error", err.Error())
	}
	queries := db.New(dbCon)

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  user.UID,
		Date: time.Now().Format("20060102"),
	})

	if err != nil {
		return views.NewDiary().Render(c.Response().Writer)
	} else {
		return views.GetDiary(diary).Render(c.Response().Writer)
	}
}

func Save(c echo.Context) error {
	user, ok := c.Get("user").(*auth.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
	}

	id := c.FormValue("id")
	content := c.FormValue("content")

	dbCon, err := connection.AppDBOpen()
	if err != nil {
		slog.Error("Failed to open database", "error", err.Error())
	}
	queries := db.New(dbCon)

	var diary db.Diary
	if id == "" {
		diary, err = queries.CreateDiary(c.Request().Context(), db.CreateDiaryParams{
			Uid:     user.UID,
			Content: content,
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "등록 실패")
		}
	} else {
		diary, err = queries.UpdateDiary(c.Request().Context(), db.UpdateDiaryParams{
			Content: content,
			ID:      id,
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "수정 실패")
		}
	}

	return views.DiaryID(diary.ID).Render(c.Response().Writer)
}
