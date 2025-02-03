package main

import (
	"io/fs"
	"os"

	resources "simple-server"
	"simple-server/internal"
	"simple-server/projects/ai-study/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	/* 환경 설정 */
	internal.LoadEnv()
	os.Setenv("APP_NAME", "🕵️‍♀️ AI 공부 길잡이")
	os.Setenv("APP_DATABASE_URL", "file:./projects/homepage/pb_data/data.db")
	os.Setenv("LOG_DATABASE_URL", "file:./projects/homepage/pb_data/auxiliary.db")
	/* 환경 설정 */

	/* 로깅 초기화 */
	internal.LoggerWithDatabaseInit()
	/* 로깅 초기화 */

	e := echo.New()

	/* 미들 웨어 */
	var sharedStaticFS fs.FS
	var projectStaticFS fs.FS
	if internal.IsProdEnv() {
		sharedStaticFS, _ = fs.Sub(resources.EmbeddedFiles, "shared/static")
		projectStaticFS, _ = fs.Sub(resources.EmbeddedFiles, "projects/ai-study/static")
	} else {
		sharedStaticFS = os.DirFS("./shared/static")
		projectStaticFS = os.DirFS("./projects/ai-study/static")
	}

	e.StaticFS("/shared/static", sharedStaticFS) // 공통 정적 파일
	e.StaticFS("/static", projectStaticFS)       // 프로젝트 정적 파일
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID:  true,
		LogLatency:    true,
		LogError:      true,
		LogRemoteIP:   true,
		LogValuesFunc: internal.CustomLogValuesFunc,
	}))

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
