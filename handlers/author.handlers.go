package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"simple-server/database"
	"simple-server/views"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func dbConnection() (*database.Queries, context.Context) {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "file:./database/data.db")
	if err != nil {
		slog.Error(err.Error())
	}
	queries := database.New(db)
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
	idInt, _ := strconv.ParseInt(id, 10, 64)

	queries, ctx := dbConnection()
	author, err := queries.GetAuthor(ctx, idInt)
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
	author, err := queries.CreateAuthor(ctx, database.CreateAuthorParams{
		Name: name,
		Bio: sql.NullString{
			String: bio,
			Valid:  true, // bio가 유효하므로 true 설정
		},
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
	idInt, _ := strconv.ParseInt(id, 10, 64)

	name := c.FormValue("name")
	bio := c.FormValue("bio")

	queries, ctx := dbConnection()
	author, err := queries.UpdateAuthor(ctx, database.UpdateAuthorParams{
		ID:   idInt,
		Name: name,
		Bio: sql.NullString{
			String: bio,
			Valid:  true, // bio가 유효하므로 true 설정
		},
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
	idInt, _ := strconv.ParseInt(id, 10, 64)

	queries, ctx := dbConnection()
	err := queries.DeleteAuthor(ctx, idInt)
	if err != nil {
		slog.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "삭제 오류")
	}

	return GetAuthors(c)
}

func ResetForm(c echo.Context) error {
	return templ.Handler(views.AuthorInsertForm()).Component.Render(c.Request().Context(), c.Response().Writer)
}
