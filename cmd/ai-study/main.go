package main

import (
	"os"

	"simple-server/internal"
	"simple-server/projects/ai-study/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	/* 환경 설정 */
	internal.LoadEnv()
	os.Setenv("SERVICE_NAME", "ai-study")
	os.Setenv("APP_TITLE", "🕵️‍♀️ AI 공부 길잡이")
	os.Setenv("APP_DATABASE_URL", "file:./projects/homepage/pb_data/data.db")
	os.Setenv("LOG_DATABASE_URL", "file:./projects/homepage/pb_data/auxiliary.db")
	/* 환경 설정 */

	/* 로깅 초기화 */
	internal.LoggerWithDatabaseInit()
	/* 로깅 초기화 */

	e := echo.New()

	/* 미들 웨어 */
	internal.RegisterCommonMiddleware(e, os.Getenv("SERVICE_NAME"))

	// 공개 그룹
	public := e.Group("")

	/* 라우터  */
	public.GET("/", handlers.IndexPageHandler)

	public.POST("/ai-study", func(c echo.Context) error {
		return handlers.AIStudy(c, false)
	})
	public.POST("/ai-study-ramdom", func(c echo.Context) error {
		return handlers.AIStudy(c, true)
	})
	/* 라우터  */

	e.Logger.Fatal(e.Start(":8001"))
}
