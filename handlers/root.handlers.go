package handlers

import (
	"simple-server/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RootHandler(c echo.Context) error {
	return templ.Handler(views.Index()).Component.Render(c.Request().Context(), c.Response().Writer)
}
