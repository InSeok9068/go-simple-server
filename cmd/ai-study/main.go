package main

import (
	"log/slog"
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

	/* 라우터  */
	if err := middleware.RegisterCommonMiddleware(e); err != nil {
		slog.Error("공통 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
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
