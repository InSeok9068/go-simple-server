package main

import (
	"io/fs"
	"log/slog"
	"os"
	resources "simple-server"
	"simple-server/projects/deario/handlers"
	"simple-server/projects/deario/services"
	"simple-server/projects/deario/tasks"

	"github.com/robfig/cron/v3"

	"simple-server/internal/config"
	"simple-server/internal/connection"
	"simple-server/internal/middleware"
	"simple-server/internal/migration"

	"github.com/labstack/echo/v4"
)

func main() {
	/* 환경 설정 */
	config.LoadEnv()
	os.Setenv("SERVICE_NAME", "deario")
	os.Setenv("APP_TITLE", "Deario")
	os.Setenv("APP_DATABASE_URL", config.AppDatabaseURL(os.Getenv("SERVICE_NAME")))
	/* 환경 설정 */

	/* 로깅 및 트레이서 초기화 */
	config.InitLoggerWithDatabase()
	config.InitTracer()
	// defer config.ShutdownTracer(context.Background())
	/* 로깅 및 트레이서 초기화 */

	/* DB 마이그레이션 */
	db, _ := connection.AppDBOpen()
	migrations, _ := fs.Sub(resources.EmbeddedFiles, "projects/deario/migrations")
	if err := migration.Up(db, migrations); err != nil {
		slog.Error("마이그레이션 실패", "error", err)
		os.Exit(1)
	}
	/* DB 마이그레이션 */

	e := setUpServer()

	/* 개발은 GoVisual, 운영은 Echo */
	if config.IsDevEnv() {
		server := config.TransferEchoToGoVisualServerOnlyDev(e, "8002")
		slog.Info("[✅ GoVisual] http server started on [::]:8002")
		if err := server.ListenAndServe(); err != nil {
			e.Logger.Fatal("GoVisual 서버 시작 실패", "error", err)
		}
	} else {
		e.Logger.Fatal(e.Start(":8002"))
	}
}

func setUpServer() *echo.Echo {
	e := echo.New()

	// PWA 파일
	manifest, _ := fs.Sub(resources.EmbeddedFiles, "projects/deario/static/manifest.json")
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
	e.POST("/logout", handlers.Logout)
	e.GET("/diary", handlers.GetDiary)
	e.GET("/diary/list", handlers.ListDiaries)
	/* 공개 라우터 */

	/* 권한 라우터 */
	authGroup := e.Group("")
	if err := middleware.RegisterCasbinMiddleware(authGroup); err != nil {
		slog.Error("Casbin 권한 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
	authGroup.GET("/diary/random", handlers.RedirectToRandomDiary)
	authGroup.POST("/diary/save", handlers.SaveDiary)
	authGroup.GET("/diary/search", handlers.SearchDiaries)
	authGroup.GET("/ai-feedback", handlers.GetAIFeedback)
	authGroup.POST("/ai-feedback", handlers.GenerateAIFeedback)
	authGroup.POST("/ai-feedback/save", handlers.SaveAIFeedback)
	authGroup.POST("/save-pushToken", handlers.RegisterPushToken)
	authGroup.GET("/setting", handlers.SettingsPage)
	authGroup.POST("/setting", handlers.UpdateSettings)
	authGroup.POST("/diary/mood", handlers.UpdateDiaryMood)
	authGroup.GET("/statistic", handlers.StatsPage)
	authGroup.GET("/statistic/data", handlers.GetStatsData)
	authGroup.GET("/diary/images", handlers.DiaryImagesPage)
	authGroup.POST("/diary/image", handlers.UploadDiaryImage)
	authGroup.DELETE("/diary/image", handlers.DeleteDiaryImage)
	/* 권한 라우터 */

	/* 큐 리시버 */
	go tasks.PushSendJob()         // 알기 작성 알림 푸시 리시버
	go tasks.GenerateAIReportJob() // AI 리포트 생성 리시버
	/* 큐 리시버 */

	/* 스케줄 */
	c := cron.New()
	tasks.PushSendCron(c) // 일기 작성 알림 푸시
	c.Start()
	/* 스케줄 */

	return e
}
