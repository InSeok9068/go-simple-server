package handlers

import (
	"net/http"
	"os"

	"simple-server/pkg/util/authutil"
	"simple-server/projects/portfolio/services"
	"simple-server/projects/portfolio/views"

	"github.com/labstack/echo/v4"
)

func IndexPage(c echo.Context) error {
	title := os.Getenv("APP_TITLE")

	uid, err := authutil.SessionUID(c)
	if err != nil {
		uid = services.DemoUID
	}

	snapshot, err := services.LoadPortfolio(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "포트폴리오 데이터를 불러오지 못했습니다.")
	}

	if uid == services.DemoUID {
		snapshot.IsDemo = true
	}

	return views.Index(title, snapshot).Render(c.Response().Writer)
}
