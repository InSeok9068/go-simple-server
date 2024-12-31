package main

import (
	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"log"
	"log/slog"
	"os"
	"simple-server/views"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Failed to load .env file", "err", err)
	}
	slog.Info(os.Getenv("TEST_ENV"))

	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		//serves static files from the provided static dir (if exists)
		se.Router.GET("/{path...}", func(e *core.RequestEvent) error {
			apis.Static(os.DirFS("./static"), false)
			return se.Next()
		})

		// index page
		se.Router.GET("/app", func(e *core.RequestEvent) error {
			return templ.Handler(views.Index()).Component.Render(e.Request.Context(), e.Response)
		})

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

	//e := echo.New()
	//
	//e.Use(middleware.Static("static"))
	//
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	//
	//e.GET("/", handlers.RootHandler) // 페이지 렌더링
	//
	//e.GET("/authors", handlers.GetAuthors)     // 저자 리스트 조회
	//e.GET("/author", handlers.GetAuthor)       // 저자 조회
	//e.POST("/author", handlers.CreateAuthor)   // 저자 등록
	//e.PUT("/author", handlers.UpdateAuthor)    // 저자 수정
	//e.DELETE("/author", handlers.DeleteAuthor) // 저자 삭제
	//
	//e.GET("/reset-form", handlers.ResetForm) // 저자 등록폼 리셋
	//
	//e.Logger.Fatal(e.Start(":8000"))
}
