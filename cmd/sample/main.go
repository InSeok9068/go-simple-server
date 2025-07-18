package main

import (
	"log/slog"
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

	/* 로깅 및 트레이서 초기화 */
	config.InitLoggerWithDatabase()
	config.InitTracer()
	// defer config.ShutdownTracer(context.Background())
	/* 로깅 및 트레이서 초기화 */

	e := setUpServer()

	/* 개발은 GoVisual, 운영은 Echo */
	if config.IsDevEnv() {
		server := config.TransferEchoToGoVisualServerOnlyDev(e, "8999")
		slog.Info("[✅ GoVisual] http server started on [::]:8999")
		if err := server.ListenAndServe(); err != nil {
			e.Logger.Fatal("GoVisual 서버 시작 실패", "error", err)
		}
	} else {
		e.Logger.Fatal(e.Start(":8999"))
	}
}

func setUpServer() *echo.Echo {
	e := echo.New()

	/* 라우터  */
	if err := middleware.RegisterCommonMiddleware(e); err != nil {
		slog.Error("공통 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
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
