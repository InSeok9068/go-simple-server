package main

import (
	"github.com/robfig/cron/v3"
	"io/fs"
	"os"
	resources "simple-server"
	"simple-server/projects/deario/handlers"
	"simple-server/projects/deario/tasks"

	"simple-server/internal/config"
	"simple-server/internal/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	/* 환경 설정 */
	config.LoadEnv()
	os.Setenv("SERVICE_NAME", "deario")
	os.Setenv("APP_TITLE", "Deario")
	os.Setenv("APP_DATABASE_URL", "file:./projects/deario/pb_data/data.db")
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
	middleware.RegisterFirebaseAuthMiddleware(e)

	// PWA 파일
	manifest, _ := fs.Sub(resources.EmbeddedFiles, "projects/deario/static/manifest.json")
	e.StaticFS("/manifest.json", manifest)

	// Web Push 서비스워커 파일
	firebaseMessagingSw, _ := fs.Sub(resources.EmbeddedFiles, "shared/static/firebase-messaging-sw.js")
	e.StaticFS("/firebase-messaging-sw.js", firebaseMessagingSw)
	/* 미들 웨어 */

	/* 라우터  */
	e.GET("/", handlers.Index)
	e.GET("/login", handlers.Login)
	e.GET("/diary", handlers.Diary)
	e.GET("/diary/list", handlers.DiaryList)
	e.GET("/diary/random", handlers.DiaryRandom)
	e.POST("/save", handlers.Save)
	e.GET("/ai-feedback", handlers.GetAiFeedback)
	e.POST("/ai-feedback", handlers.AiFeedback)
	e.POST("/ai-feedback/save", handlers.AiFeedbackSave)

	// 푸시키 갱신
	e.POST("/save-pushToken", handlers.SavePushKey)
	/* 라우터  */

	/* 스케줄 */
	c := cron.New()
	tasks.PushTask(c)
	c.Start()
	/* 스케줄 */

	return e
}
