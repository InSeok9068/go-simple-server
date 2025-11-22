package main

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	resources "simple-server"
	"simple-server/projects/deario/ai"
	"simple-server/projects/deario/auth"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/diary"
	"simple-server/projects/deario/notification"
	"simple-server/projects/deario/privacy"
	"simple-server/projects/deario/settings"

	"github.com/robfig/cron/v3"

	"simple-server/internal/config"
	"simple-server/internal/debug"
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
	defer config.ShutdownTracer(context.Background())
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
	migrations, _ := fs.Sub(resources.EmbeddedFiles, "projects/deario/migrations")
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
		server := config.TransferEchoToGoVisualServerOnlyDev(e, "8002")
		slog.Info("[✅ GoVisual] http server started on [::]:8002")
		slog.Info("Browser Open : http://localhost:8080")
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
	if err := middleware.RegisterFirebaseAuthMiddleware(e, auth.EnsureUser); err != nil {
		slog.Error("Firebase 인증 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
	e.GET("/", diary.IndexPage)
	e.GET("/login", auth.LoginPage)
	e.GET("/privacy", privacy.PrivacyPage)
	e.POST("/logout", auth.Logout)
	e.GET("/diary", diary.GetDiary)
	e.GET("/diary/list", diary.ListDiaries)
	/* 공개 라우터 */

	/* 권한 라우터 */
	authGroup := e.Group("")
	if err := middleware.RegisterCasbinMiddleware(authGroup); err != nil {
		slog.Error("Casbin 권한 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
	authGroup.GET("/diary/month", diary.MonthlyDiaryDates)
	authGroup.GET("/diary/random", diary.RedirectToRandomDiary)
	authGroup.POST("/diary/save", diary.SaveDiary)
	authGroup.GET("/diary/search", diary.SearchDiaries)
	authGroup.GET("/ai-feedback", ai.GetAIFeedback)
	authGroup.POST("/ai-feedback", ai.GenerateAIFeedback)
	authGroup.POST("/ai-feedback/save", ai.SaveAIFeedback)
	authGroup.POST("/ai-report", ai.GenerateAIReport)
	authGroup.POST("/save-pushToken", notification.RegisterPushToken)
	authGroup.GET("/setting", settings.SettingsPage)
	authGroup.POST("/setting", settings.UpdateSettings)
	authGroup.POST("/diary/mood", diary.UpdateDiaryMood)
	authGroup.GET("/statistic", diary.StatsPage)
	authGroup.GET("/statistic/data", diary.GetStatsData)
	authGroup.GET("/diary/images", diary.DiaryImagesPage)
	authGroup.POST("/diary/image", diary.UploadDiaryImage)
	authGroup.DELETE("/diary/image", diary.DeleteDiaryImage)
	authGroup.POST("/diary/transcribe", diary.TranscribeDiaryVoice)
	/* 권한 라우터 */

	/* 큐 리시버 */
	go notification.PushSendJob()         // 알기 작성 알림 푸시 리시버
	go notification.GenerateAIReportJob() // AI 리포트 생성 리시버
	/* 큐 리시버 */

	/* 스케줄 */
	c := cron.New()
	notification.PushSendCron(c) // 일기 작성 알림 푸시
	c.Start()
	/* 스케줄 */

	return e
}
