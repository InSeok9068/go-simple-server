package middleware

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	resources "simple-server"
	"simple-server/internal/config"

	"github.com/doganarif/govisual"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func RegisterCommonMiddleware(e *echo.Echo) error {
	serviceName := os.Getenv("SERVICE_NAME")

	var sharedStaticFS fs.FS
	var projectStaticFS fs.FS
	projectStaticDir := fmt.Sprintf("projects/%s/static", serviceName)
	if config.IsProdEnv() {
		sharedStaticFS, _ = fs.Sub(resources.EmbeddedFiles, "shared/static")
		projectStaticFS, _ = fs.Sub(resources.EmbeddedFiles, projectStaticDir)
	} else {
		sharedStaticFS = os.DirFS("shared/static")
		projectStaticFS = os.DirFS(projectStaticDir)
	}

	e.StaticFS("/shared/static", sharedStaticFS) // 공통 정적 파일
	e.StaticFS("/static", projectStaticFS)       // 프로젝트 정적 파일
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.Timeout())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("5M"))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20)))) // 1초당 20회 제한
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieHTTPOnly: false,
		CookieSecure:   config.IsProdEnv(),
		CookieSameSite: http.SameSiteLaxMode,
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID:  true,
		LogLatency:    true,
		LogError:      true,
		LogRemoteIP:   true,
		LogValuesFunc: config.CustomLogValuesFunc,
	}))

	return nil
}

func RegisterGoVisualMiddleware(e *echo.Echo) {
	// govisual.Handler 설정
	visualHandler := govisual.Wrap(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 이 핸들러는 실제 호출되지 않음 — Echo에서 직접 라우팅 처리됨
		}),
		govisual.WithRequestBodyLogging(true),
		govisual.WithResponseBodyLogging(true),
	)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Echo context에서 http.Request, http.ResponseWriter 가져오기
			req := c.Request()
			res := c.Response()

			// govisual 내부 미들웨어 호출 (라우팅은 건너뜀)
			visualHandler.ServeHTTP(res, req)

			// 실제 Echo 핸들러 실행
			return next(c)
		}
	})
}
