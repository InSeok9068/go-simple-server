package auth

import (
	"net/http"

	shared "simple-server/shared/views"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// LoginPage는 로그인 화면을 렌더링한다.
func LoginPage(c echo.Context) error {
	return shared.Login().Render(c.Response().Writer)
}

// Logout은 사용자 세션을 종료한다.
func Logout(c echo.Context) error {
	sess, err := session.Get("session_v2", c)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{Path: "/", MaxAge: -1}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
