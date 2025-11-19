package handlers

import (
	"net/http"
	"os"

	"simple-server/pkg/util/authutil"
	"simple-server/projects/closet/db"
	"simple-server/projects/closet/views"

	"github.com/labstack/echo/v4"
)

// IndexPage는 메인 페이지를 렌더링한다.
func IndexPage(c echo.Context) error {
	uid, _ := authutil.SessionUID(c)

	queries, err := db.GetQueries()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "데이터베이스 연결에 실패했습니다.")
	}

	groups, err := loadGroupedItems(c.Request().Context(), queries, uid, "", nil)
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)

	return views.Index(os.Getenv("APP_TITLE"), groups).Render(c.Request().Context(), c.Response().Writer)
}
