package main

import (
	"os"
	"simple-server/internal/config"
	"simple-server/internal/middleware"
	"simple-server/projects/sample/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func main() {
	/* 환경 설정 */
	config.LoadEnv()
	os.Setenv("SERVICE_NAME", "sample")
	os.Setenv("APP_TITLE", "샘플")
	/* 환경 설정 */

	e := setUpServer()

	e.Logger.Fatal(e.Start(":8002"))
}

func setUpServer() *echo.Echo {
	e := echo.New()

	/* 미들 웨어 */
	middleware.RegisterCommonMiddleware(e, os.Getenv("SERVICE_NAME"))
	/* 미들 웨어 */

	/* 라우터  */
	e.GET("/", func(c echo.Context) error {
		return templ.Handler(views.Index("Sample")).Component.Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/radio", func(c echo.Context) error {
		return templ.Handler(views.Radio()).Component.Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/radio2", func(c echo.Context) error {
		return templ.Handler(views.Radio2()).Component.Render(c.Request().Context(), c.Response().Writer)
	})
	/* 라우터  */

	return e
}
