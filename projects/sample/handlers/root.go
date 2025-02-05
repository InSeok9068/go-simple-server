package handlers

import (
	"os"
	"simple-server/projects/homepage/views"
	shared "simple-server/shared/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func IndexPageHandler(c echo.Context) error {
	return templ.Handler(views.Index(os.Getenv("APP_NAME"))).Component.Render(c.Request().Context(), c.Response().Writer)
}

func LoginPageHanlder(c echo.Context) error {
	return templ.Handler(shared.Login()).Component.Render(c.Request().Context(), c.Response().Writer)
}
