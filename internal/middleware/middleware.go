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
	// 공통 오류 핸들러 등록
	RegisterErrorHandler(e)

	// 전역 검증기 등록 (go-playground/validator, 한국어 번역)
	e.Validator = validate.NewEchoValidator()

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

	// 1) 트레이싱: 가장 바깥에서 전체 구간을 감싸기
	e.Use(otelecho.Middleware(serviceName, otelecho.WithSkipper(func(c echo.Context) bool {
		return isSkippedPath(c.Path())
	})))

	// 2) 패닉 방지: 환경과 무관하게 회복해 장애 전파를 차단
	e.Use(middleware.Recover())

	// 3) 요청 식별은 초기에 부여 (로그/트레이스 속성에 활용)
	e.Use(middleware.RequestID())

	// 4) 보안/프리플라이트 등 가벼운 것들을 먼저 처리
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	// 5) QoS/용량 제한: 비용 큰 작업 전에 컷
	if config.IsProdEnv() {
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20)))) // 1초당 20회
		e.Use(middleware.Gzip())                                                            // 응답 압축 (거부될 요청은 레이트리미터가 먼저 컷)
	}
	e.Use(middleware.BodyLimit("5M"))

	// 6) CSRF는 본문 파싱/쿠키 세팅 이후
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieHTTPOnly: false,
		CookieSecure:   config.IsProdEnv(),
		CookieSameSite: http.SameSiteLaxMode,
	}))

	// 7) 요청 로깅 (스킵 경로 동일 적용 권장)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID:  true,
		LogLatency:    true,
		LogError:      true,
		LogRemoteIP:   true,
		LogValuesFunc: config.CustomLogValuesFunc,
		Skipper:       func(c echo.Context) bool { return isSkippedPath(c.Path()) },
	}))

	// 8) 커스텀 메트릭 (트레이싱/로깅 이후~타임아웃 이전)
	e.Use(debug.MetricsMiddleware)

	// 9) 타임아웃: 항상 가장 안쪽 (Writer 바꿔치기 → 200 이슈 방지)
	e.Use(middleware.ContextTimeout(1 * time.Minute))

	// Debug
	// https://{서버주소}/debug/vars/ui?auth={인증값}
	debugGroup := e.Group("/debug")
	authParam := os.Getenv("DEBUG_AUTH_PARAM")
	if config.IsProdEnv() {
		debugGroup.Use(ipfilter.MiddlewareWithConfig(ipfilter.Config{
			WhiteList:      []string{"121.190.49.104/32"},
			BlockByDefault: true,
			Skipper:        func(c echo.Context) bool { return authParam != "" && c.QueryParam("auth") == authParam },
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
