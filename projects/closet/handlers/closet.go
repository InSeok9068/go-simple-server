package handlers

import (
	"os"

	"simple-server/projects/closet/views"

	"github.com/labstack/echo/v4"
)

// IndexPage는 메인 페이지를 렌더링한다.
func IndexPage(c echo.Context) error {
	return views.Index(os.Getenv("APP_TITLE")).Render(c.Response().Writer)
}
