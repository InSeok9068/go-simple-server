package main

import (
	"io/fs"
	"log/slog"
	"os"
	resources "simple-server"
	"simple-server/projects/portfolio/db"
	"simple-server/projects/portfolio/handlers"
	"simple-server/projects/portfolio/services"

	"simple-server/internal/config"
	"simple-server/internal/debug"
	"simple-server/internal/middleware"
	"simple-server/internal/migration"

	"github.com/labstack/echo/v4"
)

func main() {
	/* 환경 설정 */
	config.LoadEnv()
	os.Setenv("SERVICE_NAME", "portfolio")
	os.Setenv("APP_TITLE", "Portfolio")
	os.Setenv("APP_DATABASE_URL", config.AppDatabaseURL(os.Getenv("SERVICE_NAME")))
	/* 환경 설정 */

	/* 로깅 및 트레이서 초기화 */
	config.InitLoggerWithDatabase()
	config.InitTracer()
	// defer config.ShutdownTracer(context.Background())
	/* 로깅 및 트레이서 초기화 */

	/* DB 초기화 */
	database, err := db.GetDB()
	if err != nil {
		slog.Error("데이터베이스 연결 실패", "error", err)
		return
	}
	defer database.Close()
	/* DB 초기화 */

	/* DB 마이그레이션 */
	migrations, _ := fs.Sub(resources.EmbeddedFiles, "projects/portfolio/migrations")
	if err := migration.Up(database, migrations); err != nil {
		slog.Error("마이그레이션 실패", "error", err)
		return
	}
	/* DB 마이그레이션 */

	/* 디버그 지표 노출 */
	debug.Init(os.Getenv("SERVICE_NAME"), database)
	/* 디버그 지표 노출 */

	e := setUpServer()

	/* 개발은 GoVisual, 운영은 Echo */
	if config.IsDevEnv() {
		server := config.TransferEchoToGoVisualServerOnlyDev(e, "8003")
		slog.Info("[✅ GoVisual] http server started on [::]:8003")
		if err := server.ListenAndServe(); err != nil {
			e.Logger.Fatal("GoVisual 서버 시작 실패", "error", err)
		}
	} else {
		e.Logger.Fatal(e.Start(":8003"))
	}
}

func setUpServer() *echo.Echo {
	e := echo.New()

	// PWA 파일
	manifest, _ := fs.Sub(resources.EmbeddedFiles, "projects/portfolio/static/manifest.json")
	e.StaticFS("/manifest.json", manifest)

	// Web Push 서비스워커 파일
	firebaseMessagingSw, _ := fs.Sub(resources.EmbeddedFiles, "shared/static/firebase-messaging-sw.js")
	e.StaticFS("/firebase-messaging-sw.js", firebaseMessagingSw)

	/* 공개 라우터 */
	if err := middleware.RegisterCommonMiddleware(e); err != nil {
		slog.Error("공통 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
	if err := middleware.RegisterFirebaseAuthMiddleware(e, services.EnsureUser); err != nil {
		slog.Error("Firebase 인증 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
	e.GET("/", handlers.IndexPage)
	e.GET("/login", handlers.LoginPage)
	/* 공개 라우터 */

	/* 권한 라우터 */
	authGroup := e.Group("")
	if err := middleware.RegisterCasbinMiddleware(authGroup); err != nil {
		slog.Error("Casbin 권한 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
	/* 권한 라우터 */

	return e
}
