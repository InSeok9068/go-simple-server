package handlers

import (
	"simple-server/projects/deario/views/pages"

	"github.com/labstack/echo/v4"
)

// PrivacyPage는 개인정보 처리방침 페이지를 렌더링한다.
func PrivacyPage(c echo.Context) error {
	return pages.Privacy().Render(c.Request().Context(), c.Response().Writer)
}
