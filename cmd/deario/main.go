package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"os"
	resources "simple-server"
	"simple-server/projects/deario/handlers"
	"simple-server/projects/deario/services"
	"simple-server/projects/deario/tasks"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/robfig/cron/v3"

	"simple-server/internal/config"
	"simple-server/internal/connection"
	"simple-server/internal/middleware"

	"github.com/labstack/echo/v4"
)

var migrations embed.FS

func main() {
	/* 환경 설정 */
	config.LoadEnv()
	os.Setenv("SERVICE_NAME", "deario")
	os.Setenv("APP_TITLE", "Deario")
	os.Setenv("APP_DATABASE_URL", config.AppDatabaseURL(os.Getenv("SERVICE_NAME")))
	/* 환경 설정 */

	/* 로깅 초기화 */
	config.LoggerWithDatabaseInit()
	/* 로깅 초기화 */

	/* DB 마이그레이션 */
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, _ := connection.AppDBOpen()
	migrations, _ := fs.Sub(resources.EmbeddedFiles, "projects/deario/migrations")
	provider, _ := goose.NewProvider(
		goose.DialectSQLite3, db, migrations,
		goose.WithVerbose(config.IsDevEnv()),
	)
	if _, err := provider.Up(ctx); err != nil {
		slog.Error("마이그레이션 실패", "error", err)
		os.Exit(1)
	}
	/* DB 마이그레이션 */

	e := setUpServer()

	e.Logger.Fatal(e.Start(":8002"))
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
	e.GET("/", handlers.Index)
	e.GET("/login", handlers.Login)
	e.POST("/logout", handlers.Logout)
	e.GET("/diary", handlers.Diary)
	e.GET("/diary/list", handlers.DiaryList)
	/* 공개 라우터 */

	/* 권한 라우터 */
	authGroup := e.Group("")
	if err := middleware.RegisterCasbinMiddleware(authGroup); err != nil {
		slog.Error("Casbin 권한 미들웨어 등록 실패", "error", err)
		os.Exit(1)
	}
	authGroup.GET("/diary/random", handlers.DiaryRandom)
	authGroup.POST("/diary/save", handlers.Save)
	authGroup.GET("/ai-feedback", handlers.GetAiFeedback)
	authGroup.POST("/ai-feedback", handlers.AiFeedback)
	authGroup.POST("/ai-feedback/save", handlers.AiFeedbackSave)
	authGroup.POST("/save-pushToken", handlers.SavePushKey)
	/* 권한 라우터 */

	/* 스케줄 */
	c := cron.New()
	tasks.PushTask(c)
	c.Start()
	/* 스케줄 */

	return e
}
