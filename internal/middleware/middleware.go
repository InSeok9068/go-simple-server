package middleware

import (
	"expvar"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	resources "simple-server"
	"simple-server/internal/config"
	"simple-server/internal/debug"
	"simple-server/internal/validate"
	"strings"
	"time"

	ipfilter "github.com/crazy-max/echo-ipfilter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
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
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20)))) // 1초당 20회 제한
	}
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      1 * time.Minute,
		ErrorMessage: "요청 처리 시간이 지연되었습니다. 다시 시도해주세요.",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("5M"))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieHTTPOnly: false,
		CookieSecure:   config.IsProdEnv(),
		CookieSameSite: http.SameSiteLaxMode,
	}))
	e.Use(debug.MetricsMiddleware)
	e.Use(otelecho.Middleware(serviceName, otelecho.WithSkipper(func(c echo.Context) bool {
		return isSkippedPath(c.Path())
	})))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID:  true,
		LogLatency:    true,
		LogError:      true,
		LogRemoteIP:   true,
		LogValuesFunc: config.CustomLogValuesFunc,
		Skipper: func(c echo.Context) bool {
			return isSkippedPath(c.Path())
		},
	}))
	// 전역 검증기 등록 (go-playground/validator, 한국어 번역)
	e.Validator = validate.NewEchoValidator()

	// Debug
	// https://{서버주소}/debug/vars/ui?auth={인증값}
	debugGroup := e.Group("/debug")
	authParam := os.Getenv("DEBUG_AUTH_PARAM")
	if config.IsProdEnv() {
		debugGroup.Use(ipfilter.MiddlewareWithConfig(ipfilter.Config{
			WhiteList: []string{
				"121.190.49.104/32",
			},
			BlockByDefault: true,
			Skipper: func(c echo.Context) bool {
				if authParam == "" {
					return false
				}
				return c.Request().URL.Query().Get("auth") == authParam
			},
		}))
	}
	// expvar 핸들러
	debugGroup.GET("/vars", echo.WrapHandler(expvar.Handler()))
	debugGroup.GET("/vars/ui", echo.WrapHandler(http.HandlerFunc(debug.VarsUI)))

	return nil
}

func isSkippedPath(path string) bool {
	// 추적에서 제외할 경로/패턴의 접두어 목록
	// c.Path()가 라우트 패턴("/static*")을 반환하는 경우도 고려해 '*'가 포함된 접두어도 함께 둔다.
	prefixes := []string{
		"/metrics",
		"/metrics*",
		"/static/",
		"/static*",
		"/shared/static/",
		"/shared/static*",
		"/manifest.json",
		"/manifest.json*",
		"/firebase-messaging-sw.js",
		"/firebase-messaging-sw.js*",
		"/service-worker.js",
		"/service-worker.js*",
		"/favicon.ico",
		"/favicon.ico*",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
