package main

import (
	"io/fs"
	"os"
	resources "simple-server"

	"simple-server/internal/config"
	"simple-server/internal/middleware"
	"simple-server/projects/homepage/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

func main() {
	/* 환경 설정 */
	config.LoadEnv()
	os.Setenv("SERVICE_NAME", "homepage")
	os.Setenv("APP_TITLE", "홈페이지")
	/* 환경 설정 */

	/* 로깅 초기화 */
	config.LoggerWithDatabaseInit()
	/* 로깅 초기화 */

	e := setUpServer()

	e.Logger.Fatal(e.Start(":8000"))
}

func setUpServer() *echo.Echo {
	e := echo.New()

	/* 미들 웨어 */
	middleware.RegisterCommonMiddleware(e, os.Getenv("SERVICE_NAME"))
	// PWA 파일
	manifest, _ := fs.Sub(resources.EmbeddedFiles, "projects/homepage/static/manifest.json")
	e.StaticFS("/manifest.json", manifest)

	// Prometheus 미들웨어
	e.Use(echoprometheus.NewMiddleware("homepage"))
	e.GET("/metrics", echoprometheus.NewHandler())
	/* 미들 웨어 */

	/* 라우터  */
	e.GET("/", func(c echo.Context) error {
		return templ.Handler(views.Index(os.Getenv("APP_NAME"))).Component.Render(c.Request().Context(), c.Response().Writer)
	})
	/* 라우터  */

	return e
}
