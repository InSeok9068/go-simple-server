package handlers

import (
	"os"

	"simple-server/projects/portfolio/views"

	"github.com/labstack/echo/v4"
)

func IndexPage(c echo.Context) error {
	title := os.Getenv("APP_TITLE")
	return views.Index(title).Render(c.Response().Writer)
}
