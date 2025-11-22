package main

import (
	"context"
	"log/slog"
	"os"

	"simple-server/internal/config"
	"simple-server/internal/middleware"
	"simple-server/projects/ai-study/internal/study"

	"github.com/labstack/echo/v4"
)

func main() {
	/* 환경 설정 */
	config.LoadEnv()
	os.Setenv("SERVICE_NAME", "ai-study")
	os.Setenv("APP_TITLE", "AI 공부 길잡이")
	/* 환경 설정 */

	/* 로깅 및 트레이서 초기화 */
	config.InitLoggerWithDatabase()
	config.InitTracer()
	defer config.ShutdownTracer(context.Background())
	/* 로깅 및 트레이서 초기화 */

	e := setUpServer()

	/* 개발은 GoVisual, 운영은 Echo */
	if config.IsDevEnv() {
		server := config.TransferEchoToGoVisualServerOnlyDev(e, "8001")
		slog.Info("[✅ GoVisual] http server started on [::]:8001")
		slog.Info("Browser Open : http://localhost:8080")
		if err := server.ListenAndServe(); err != nil {
			e.Logger.Fatal("GoVisual 서버 시작 실패", "error", err)
		}
	} else {
		e.Logger.Fatal(e.Start(":8001"))
	}
}

func setUpServer() *echo.Echo {
	e := echo.New()

	/* 라우터  */
	if err := middleware.RegisterCommonMiddleware(e); err != nil {
		slog.Error("공통 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
	e.GET("/", study.IndexPageHandler)

	e.POST("/ai-study", func(c echo.Context) error {
		return study.AIStudy(c, false)
	})
	e.POST("/ai-study-random", func(c echo.Context) error {
		return study.AIStudy(c, true)
	})
	/* 라우터  */

	return e
}
