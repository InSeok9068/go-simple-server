package study

import (
	"github.com/labstack/echo/v4"
	"os"
	"simple-server/projects/ai-study/views"
)

func IndexPageHandler(c echo.Context) error {
	return views.Index(os.Getenv("APP_TITLE")).Render(c.Request().Context(), c.Response().Writer)
}
