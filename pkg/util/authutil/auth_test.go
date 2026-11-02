package authutil

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func TestSessionUIDInvalidCookiePreservesCause(t *testing.T) {
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("test-session-signing-key"))))
	var sessionErr error
	e.GET("/", func(c echo.Context) error {
		_, sessionErr = SessionUID(c)
		return sessionErr
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_v2", Value: "invalid-session"})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var httpErr *echo.HTTPError
	if !errors.As(sessionErr, &httpErr) {
		t.Fatalf("HTTP 오류가 필요합니다: %v", sessionErr)
	}
	if httpErr.Internal == nil {
		t.Fatal("세션 읽기 실패의 원본 오류가 보존되어야 합니다")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("응답 상태가 올바르지 않습니다: %d", rec.Code)
	}
	want := "{\"message\":\"로그인 세션이 만료되었거나 유효하지 않습니다. 다시 로그인해주세요.\"}\n"
	if rec.Body.String() != want {
		t.Fatalf("응답에는 사용자용 메시지만 포함되어야 합니다: %s", rec.Body.String())
	}
}
