package main

import (
	"os"

	"simple-server/internal"
	"simple-server/projects/homepage/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

func main() {
	/* 환경 설정 */
	internal.LoadEnv()
	os.Setenv("SERVICE_NAME", "homepage")
	os.Setenv("APP_TITLE", "홈페이지")
	// os.Setenv("LOG_DATABASE_URL", "file:./projects/homepage/pb_data/auxiliary.db")
	/* 환경 설정 */

	/* 로깅 초기화 */
	internal.LoggerWithDatabaseInit()
	/* 로깅 초기화 */

	e := echo.New()

	/* 미들 웨어 */
	internal.RegisterCommonMiddleware(e, os.Getenv("SERVICE_NAME"))
	// PWA 파일
	e.StaticFS("/manifest.json", os.DirFS("projects/homepage/static/manifest.json"))

	// Prometheus 미들웨어
	e.Use(echoprometheus.NewMiddleware("homepage"))
	e.GET("/metrics", echoprometheus.NewHandler())
	/* 미들 웨어 */

	/* 라우터  */
	e.GET("/", func(c echo.Context) error {
		return templ.Handler(views.Index(os.Getenv("APP_NAME"))).Component.Render(c.Request().Context(), c.Response().Writer)
	})
	/* 라우터  */

	e.Logger.Fatal(e.Start(":8000"))
}
