package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"simple-server/cmd"
	"simple-server/cmd/handlers"
)

func main() {
	/* 환경 설정 */
	err := godotenv.Load()
	if err != nil {
		slog.Error("Failed to load .env file", "err", err)
	}
	/* 환경 설정 */

	/* 파이어베이스 초기화 */
	cmd.FirebaseInit()
	/* 파이어베이스 초기화 */

	e := echo.New()

	/* 미들 웨어 */
	e.Use(middleware.Static("static"))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 공개 그룹
	public := e.Group("public")

	// 인증 그룹
	private := e.Group("private")

	private.Use(middleware.KeyAuthWithConfig(cmd.FirebaseAuth()))
	/* 미들 웨어 */

	/* 라우터  */
	public.GET("/", handlers.IndexPageHandler)
	public.GET("/login", handlers.LoginPageHanlder)

	e.GET("/authors", handlers.GetAuthors)     // 저자 리스트 조회
	e.GET("/author", handlers.GetAuthor)       // 저자 조회
	e.POST("/author", handlers.CreateAuthor)   // 저자 등록
	e.PUT("/author", handlers.UpdateAuthor)    // 저자 수정
	e.DELETE("/author", handlers.DeleteAuthor) // 저자 삭제

	e.GET("/reset-form", handlers.ResetForm) // 저자 등록폼 리셋
	/* 라우터  */

	e.Logger.Fatal(e.Start(":8000"))
}
