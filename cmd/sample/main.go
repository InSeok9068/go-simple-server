package main

import (
	"io/fs"
	resources "simple-server"
	"simple-server/internal"
	"simple-server/projects/sample/handlers"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	/* 환경 설정 */
	envData, _ := resources.EmbeddedFiles.ReadFile(".env")
	godotenv.Unmarshal(string(envData))
	/* 환경 설정 */

	/* 파이어베이스 초기화 */
	internal.FirebaseInit()
	/* 파이어베이스 초기화 */

	e := echo.New()

	/* 미들 웨어 */
	sharedStaticFS, _ := fs.Sub(resources.EmbeddedFiles, "shared/static")
	projectStaticFS, _ := fs.Sub(resources.EmbeddedFiles, "projects/sample/static")
	e.StaticFS("/shared/static", sharedStaticFS) // 공통 정적 파일
	e.StaticFS("/static", projectStaticFS)       // 프로젝트 정적 파일
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 공개 그룹
	public := e.Group("")

	// 인증 그룹
	private := e.Group("")

	private.Use(middleware.KeyAuthWithConfig(internal.FirebaseAuth()))
	/* 미들 웨어 */

	/* 라우터  */
	public.GET("/", handlers.IndexPageHandler)
	public.GET("/login", handlers.LoginPageHanlder)

	private.GET("/authors", handlers.GetAuthors)     // 저자 리스트 조회
	private.GET("/author", handlers.GetAuthor)       // 저자 조회
	private.POST("/author", handlers.CreateAuthor)   // 저자 등록
	private.PUT("/author", handlers.UpdateAuthor)    // 저자 수정
	private.DELETE("/author", handlers.DeleteAuthor) // 저자 삭제

	private.GET("/reset-form", handlers.ResetForm) // 저자 등록폼 리셋
	/* 라우터  */

	e.Logger.Fatal(e.Start(":8000"))
}
