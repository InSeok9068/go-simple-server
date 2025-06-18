package main

import (
	"os"

	"simple-server/internal/config"
	"simple-server/internal/middleware"
	"simple-server/projects/ai-study/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	/* 환경 설정 */
	config.LoadEnv()
	os.Setenv("SERVICE_NAME", "ai-study")
	os.Setenv("APP_TITLE", "AI 공부 길잡이")
	/* 환경 설정 */

	/* 로깅 초기화 */
	config.LoggerWithDatabaseInit()
	/* 로깅 초기화 */

	e := setUpServer()

	e.Logger.Fatal(e.Start(":8001"))
}

func setUpServer() *echo.Echo {
	e := echo.New()

	/* 미들 웨어 */
	middleware.RegisterCommonMiddleware(e)
	/* 미들 웨어 */

	/* 라우터  */
	e.GET("/", handlers.IndexPageHandler)

	e.POST("/ai-study", func(c echo.Context) error {
		return handlers.AIStudy(c, false)
	})
	e.POST("/ai-study-random", func(c echo.Context) error {
		return handlers.AIStudy(c, true)
	})
	/* 라우터  */

	return e
}
