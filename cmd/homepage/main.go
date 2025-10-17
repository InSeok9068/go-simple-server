package main

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	resources "simple-server"

	"simple-server/internal/config"
	"simple-server/internal/middleware"
	"simple-server/projects/homepage/views"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

func main() {
	/* 환경 설정 */
	config.LoadEnv()
	os.Setenv("SERVICE_NAME", "homepage")
	os.Setenv("APP_TITLE", "홈페이지")
	/* 환경 설정 */

	/* 로깅 및 트레이서 초기화 */
	config.InitLoggerWithDatabase()
	config.InitTracer()
	defer config.ShutdownTracer(context.Background())
	/* 로깅 및 트레이서 초기화 */

	e := setUpServer()

	/* 개발은 GoVisual, 운영은 Echo */
	if config.IsDevEnv() {
		server := config.TransferEchoToGoVisualServerOnlyDev(e, "8000")
		slog.Info("[✅ GoVisual] http server started on [::]:8000")
		slog.Info("Browser Open : http://localhost:8080")
		if err := server.ListenAndServe(); err != nil {
			e.Logger.Fatal("GoVisual 서버 시작 실패", "error", err)
		}
	} else {
		e.Logger.Fatal(e.Start(":8000"))
	}
}

func setUpServer() *echo.Echo {
	e := echo.New()

	// PWA 파일
	manifest, _ := fs.Sub(resources.EmbeddedFiles, "projects/homepage/static/manifest.json")
	e.StaticFS("/manifest.json", manifest)

	// Prometheus 미들웨어
	e.Use(echoprometheus.NewMiddleware("homepage"))
	e.GET("/metrics", echoprometheus.NewHandler())

	/* 라우터  */
	if err := middleware.RegisterCommonMiddleware(e); err != nil {
		slog.Error("공통 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
	e.GET("/", func(c echo.Context) error {
		return views.Index(os.Getenv("APP_TITLE")).Render(c.Response().Writer)
	})
	/* 라우터  */

	return e
}
