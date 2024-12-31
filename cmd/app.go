package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"simple-server/handlers"
)

func AppServer() {
	// 앱 서버
	e := echo.New()

	e.Use(middleware.Static("static"))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handlers.RootHandler) // 페이지 렌더링

	e.GET("/authors", handlers.GetAuthors)     // 저자 리스트 조회
	e.GET("/author", handlers.GetAuthor)       // 저자 조회
	e.POST("/author", handlers.CreateAuthor)   // 저자 등록
	e.PUT("/author", handlers.UpdateAuthor)    // 저자 수정
	e.DELETE("/author", handlers.DeleteAuthor) // 저자 삭제

	e.GET("/reset-form", handlers.ResetForm) // 저자 등록폼 리셋

	e.Logger.Fatal(e.Start(":8000"))
}
