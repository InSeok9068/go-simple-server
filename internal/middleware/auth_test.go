package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func TestSaveFirebaseSessionReplacesExistingCookie(t *testing.T) {
	store := sessions.NewCookieStore([]byte("current-session-signing-key"))
	oldStore := sessions.NewCookieStore([]byte("previous-session-signing-key"))
	oldSession := sessions.NewSession(oldStore, "session_v2")
	oldSession.Options = &sessions.Options{Path: "/", MaxAge: 3600}
	oldSession.Values["uid"] = "previous-user"
	oldResponse := httptest.NewRecorder()
	if err := oldSession.Save(httptest.NewRequest(http.MethodGet, "/", nil), oldResponse); err != nil {
		t.Fatalf("기존 쿠키 준비 실패: %v", err)
	}

	for _, tc := range []struct {
		name   string
		cookie *http.Cookie
	}{
		{name: "쿠키 없음"},
		{name: "손상된 쿠키", cookie: &http.Cookie{Name: "session_v2", Value: "invalid-session"}},
		{name: "이전 서명 키의 쿠키", cookie: oldResponse.Result().Cookies()[0]},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/create-session", nil)
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			if err := saveFirebaseSession(c, store, "verified-user"); err != nil {
				t.Fatalf("새 세션 저장 실패: %v", err)
			}
			cookies := rec.Result().Cookies()
			if len(cookies) != 1 {
				t.Fatalf("새 로그인 쿠키 하나가 필요합니다: %d", len(cookies))
			}
			cookie := cookies[0]
			if cookie.Name != "session_v2" || cookie.Path != "/" || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode {
				t.Fatalf("로그인 쿠키의 보안 속성이 올바르지 않습니다")
			}
			nextReq := httptest.NewRequest(http.MethodGet, "/", nil)
			nextReq.AddCookie(cookie)
			sess, err := store.Get(nextReq, "session_v2")
			if err != nil {
				t.Fatalf("새 세션 읽기 실패: %v", err)
			}
			if sess.Values["uid"] != "verified-user" {
				t.Fatal("새 세션은 검증된 사용자 식별자를 사용해야 합니다")
			}
		})
	}
}
