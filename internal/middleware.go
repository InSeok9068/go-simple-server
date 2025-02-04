package internal

import (
	"fmt"
	"io/fs"
	"os"
	resources "simple-server"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterCommonMiddleware(e *echo.Echo, serviceName string) {
	var sharedStaticFS fs.FS
	var projectStaticFS fs.FS
	projectStaticDir := fmt.Sprintf("projects/%s/static", serviceName)
	if IsProdEnv() {
		sharedStaticFS, _ = fs.Sub(resources.EmbeddedFiles, "shared/static")
		projectStaticFS, _ = fs.Sub(resources.EmbeddedFiles, projectStaticDir)
	} else {
		sharedStaticFS = os.DirFS("shared/static")
		projectStaticFS = os.DirFS(projectStaticDir)
	}

	e.StaticFS("/shared/static", sharedStaticFS) // 공통 정적 파일
	e.StaticFS("/static", projectStaticFS)       // 프로젝트 정적 파일
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID:  true,
		LogLatency:    true,
		LogError:      true,
		LogRemoteIP:   true,
		LogValuesFunc: CustomLogValuesFunc,
	}))
}
