package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"net/http"
	"simple-server/internal/main/db"
	"simple-server/views"
)

func dbConnection() (*db.Queries, context.Context) {
	ctx := context.Background()
	dbCon, err := sql.Open("sqlite3", "file:./pb_data/data.db")
	if err != nil {
		slog.Error(err.Error())
	}
	queries := db.New(dbCon)
	return queries, ctx
}

func GetAuthors(c echo.Context) error {
	queries, ctx := dbConnection()
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		slog.Error(err.Error())
	}
	return templ.Handler(views.Authors(authors)).Component.Render(c.Request().Context(), c.Response().Writer)
	// return echo.NewHTTPError(http.StatusInternalServerError, "오류 입니다.")
}

func GetAuthor(c echo.Context) error {
	id := c.QueryParam("id")

	queries, ctx := dbConnection()
	author, err := queries.GetAuthor(ctx, id)
	if err != nil {
		slog.Error(err.Error())
	}

	return templ.Handler(views.AuthorUpdateForm(author)).Component.Render(c.Request().Context(), c.Response().Writer)
}

func CreateAuthor(c echo.Context) error {
	name := c.FormValue("name")
	bio := c.FormValue("bio")

	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "이름을 입력해주세요.")
	}

	queries, ctx := dbConnection()
	author, err := queries.CreateAuthor(ctx, db.CreateAuthorParams{
		Name: name,
		Bio:  bio,
	})
	if err != nil {
		slog.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "등록 오류")
	}

	slog.Info(fmt.Sprintf("Author: %+v", author))

	return GetAuthors(c)
}

func UpdateAuthor(c echo.Context) error {
	id := c.FormValue("id")

	name := c.FormValue("name")
	bio := c.FormValue("bio")

	queries, ctx := dbConnection()
	author, err := queries.UpdateAuthor(ctx, db.UpdateAuthorParams{
		ID:   id,
		Name: name,
		Bio:  bio,
	})
	if err != nil {
		slog.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "수정 오류")
	}

	slog.Info(fmt.Sprintf("Author: %+v", author))

	return GetAuthors(c)
}

func DeleteAuthor(c echo.Context) error {
	id := c.QueryParam("id")

	queries, ctx := dbConnection()
	err := queries.DeleteAuthor(ctx, id)
	if err != nil {
		slog.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "삭제 오류")
	}

	return GetAuthors(c)
}

func ResetForm(c echo.Context) error {
	return templ.Handler(views.AuthorInsertForm()).Component.Render(c.Request().Context(), c.Response().Writer)
}
