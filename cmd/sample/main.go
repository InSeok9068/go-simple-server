package main

import (
	"os"
	"simple-server/internal/config"
	"simple-server/internal/middleware"
	"simple-server/projects/sample/views"

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
		// 코드 리뷰 한번 해주세요
		return views.Index("Sample").Render(c.Response().Writer)
	})

	e.GET("/radio", func(c echo.Context) error {
		return views.Radio().Render(c.Response().Writer)
	})

	e.GET("/radio2", func(c echo.Context) error {
		return views.Radio2().Render(c.Response().Writer)
	})
	/* 라우터  */

	return e
}
