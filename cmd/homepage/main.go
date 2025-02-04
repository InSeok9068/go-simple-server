package main

import (
	"io/fs"
	"os"

	resources "simple-server"
	"simple-server/internal"
	"simple-server/projects/homepage/handlers"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	/* 환경 설정 */
	internal.LoadEnv()
	os.Setenv("APP_NAME", "홈페이지")
	os.Setenv("LOG_DATABASE_URL", "file:./projects/homepage/pb_data/auxiliary.db")
	/* 환경 설정 */

	/* 로깅 초기화 */
	internal.LoggerWithDatabaseInit()
	/* 로깅 초기화 */

	/* 파이어베이스 초기화 */
	internal.FirebaseInit()
	/* 파이어베이스 초기화 */

	e := echo.New()

	/* 미들 웨어 */
	var sharedStaticFS fs.FS
	var projectStaticFS fs.FS
	if internal.IsProdEnv() {
		sharedStaticFS, _ = fs.Sub(resources.EmbeddedFiles, "shared/static")
		projectStaticFS, _ = fs.Sub(resources.EmbeddedFiles, "projects/homepage/static")
	} else {
		sharedStaticFS = os.DirFS("./shared/static")
		projectStaticFS = os.DirFS("./projects/homepage/static")
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

	// Prometheus 미들웨어
	e.Use(echoprometheus.NewMiddleware("homepage"))
	e.GET("/metrics", echoprometheus.NewHandler())

	// 공개 그룹
	public := e.Group("")

	// 인증 그룹
	private := e.Group("")
	private.Use(middleware.KeyAuthWithConfig(internal.FirebaseAuth()))
	/* 미들 웨어 */

	/* 라우터  */
	public.GET("/", handlers.IndexPageHandler)
	/* 라우터  */

	e.Logger.Fatal(e.Start(":8000"))
}
