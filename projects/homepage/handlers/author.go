package handlers

import (
	"log/slog"
	"net/http"
	"simple-server/projects/homepage/services"
	"simple-server/projects/homepage/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func GetAuthors(c echo.Context) error {
	authors, err := services.GetAuthors()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "목록 조회 오류")
	}
	return templ.Handler(views.Authors(authors)).Component.Render(c.Request().Context(), c.Response().Writer)
}

func GetAuthor(c echo.Context) error {
	id := c.QueryParam("id")

	author, err := services.GetAuthor(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "조회 오류")
	}

	return templ.Handler(views.AuthorUpdateForm(author)).Component.Render(c.Request().Context(), c.Response().Writer)
}

func CreateAuthor(c echo.Context) error {
	name := c.FormValue("name")
	bio := c.FormValue("bio")

	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "이름을 입력해주세요.")
	}

	_, err := services.CreateAuthor(name, bio)
	if err != nil {
		slog.Error("저자 등록 오류", "error", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "등록 오류")
	}

	return GetAuthors(c)
}

func UpdateAuthor(c echo.Context) error {
	id := c.FormValue("id")
	name := c.FormValue("name")
	bio := c.FormValue("bio")

	_, err := services.UpdateAuthor(id, name, bio)
	if err != nil {
		slog.Error("저자 수정 오류", "error", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "수정 오류")
	}

	return GetAuthors(c)
}

func DeleteAuthor(c echo.Context) error {
	id := c.QueryParam("id")

	err := services.DeleteAuthor(id)
	if err != nil {
		slog.Error("저자 삭제 오류", "error", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "삭제 오류")
	}

	return GetAuthors(c)
}

func ResetForm(c echo.Context) error {
	return templ.Handler(views.AuthorInsertForm()).Component.Render(c.Request().Context(), c.Response().Writer)
}
