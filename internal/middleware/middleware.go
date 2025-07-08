package middleware

import (
	"expvar"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	resources "simple-server"
	"simple-server/internal/config"

	ipfilter "github.com/crazy-max/echo-ipfilter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"golang.org/x/time/rate"
)

func RegisterCommonMiddleware(e *echo.Echo) error {
	RegisterErrorHandler(e)

	serviceName := os.Getenv("SERVICE_NAME")

	var sharedStaticFS fs.FS
	var projectStaticFS fs.FS
	projectStaticDir := fmt.Sprintf("projects/%s/static", serviceName)
	if config.IsProdEnv() {
		var err error
		if sharedStaticFS, err = fs.Sub(resources.EmbeddedFiles, "shared/static"); err != nil {
			return fmt.Errorf("정적 파일 시스템 초기화 실패: %w", err)
		}
		if projectStaticFS, err = fs.Sub(resources.EmbeddedFiles, projectStaticDir); err != nil {
			return fmt.Errorf("프로젝트 정적 파일 시스템 초기화 실패: %w", err)
		}
	} else {
		sharedStaticFS = os.DirFS("shared/static")
		projectStaticFS = os.DirFS(projectStaticDir)
	}

	e.StaticFS("/shared/static", sharedStaticFS) // 공통 정적 파일
	e.StaticFS("/static", projectStaticFS)       // 프로젝트 정적 파일

	// 개발환경은 GoVisual 확인을 위해서 Gzip 미적용
	if config.IsProdEnv() {
		e.Use(middleware.Gzip())
	}
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
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

	// Tracing
	tracer := otel.Tracer(serviceName)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, span := tracer.Start(c.Request().Context(), c.Request().Method+" "+c.Path())
			defer span.End()
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})

	// Debug
	debugGroup := e.Group("/debug")
	if config.IsProdEnv() {
		debugGroup.Use(ipfilter.MiddlewareWithConfig(ipfilter.Config{
			WhiteList: []string{
				"121.190.49.104/32",
			},
			BlockByDefault: true,
		}))
	}
	debugGroup.GET("/vars", echo.WrapHandler(expvar.Handler()))

	return nil
}
