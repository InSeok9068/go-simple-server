package main

import (
	"os"
	"simple-server/projects/deario/handlers"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"simple-server/internal/config"
	"simple-server/internal/middleware"
)

func main() {
	/* 환경 설정 */
	config.LoadEnv()
	os.Setenv("SERVICE_NAME", "deario")
	os.Setenv("APP_TITLE", "Deario")
	/* 환경 설정 */

	/* 로깅 초기화 */
	config.LoggerWithDatabaseInit()
	/* 로깅 초기화 */

	e := setUpServer()

	e.Logger.Fatal(e.Start(":8002"))
}

func setUpServer() *echo.Echo {
	e := echo.New()

	/* 미들 웨어 */
	middleware.RegisterCommonMiddleware(e, os.Getenv("SERVICE_NAME"))

	middleware.FirebaseInit()
	e.Use(echoMiddleware.KeyAuthWithConfig(middleware.FirebaseAuth()))

	// PWA 파일
	// manifest, _ := fs.Sub(resources.EmbeddedFiles, "projects/homepage/static/manifest.json")
	// e.StaticFS("/manifest.json", manifest)
	/* 미들 웨어 */

	/* 라우터  */
	e.GET("/", handlers.Index)
	e.GET("/login", handlers.Login)
	e.GET("/save", handlers.Save)
	/* 라우터  */

	return e
}
