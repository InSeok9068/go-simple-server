package handlers

import (
	"os"

	"simple-server/projects/portfolio/views"

	"github.com/labstack/echo/v4"
)

func IndexPage(c echo.Context) error {
	return views.Index(os.Getenv("APP_TITLE")).Render(c.Response().Writer)
}
