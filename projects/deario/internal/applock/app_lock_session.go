package applock

import (
	"net/http"

	"simple-server/internal/config"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const unlockSessionName = "deario_app_unlock"

// SaveUnlockSession은 브라우저 세션 동안만 유지되는 앱 잠금 해제 상태를 저장한다.
func SaveUnlockSession(c echo.Context, uid string) error {
	sess, err := session.Get(unlockSessionName, c)
	if err != nil {
		return err
	}

	sess.Options = unlockSessionOptions(0)
	sess.Values["uid"] = uid
	sess.Values["unlocked"] = true

	return sess.Save(c.Request(), c.Response())
}

// ClearUnlockSession은 앱 잠금 해제 상태를 제거한다.
func ClearUnlockSession(c echo.Context) error {
	sess, err := session.Get(unlockSessionName, c)
	if err != nil {
		return err
	}

	sess.Options = unlockSessionOptions(-1)
	sess.Values = map[interface{}]interface{}{}

	return sess.Save(c.Request(), c.Response())
}

func isUnlockSessionValid(c echo.Context, uid string) bool {
	sess, err := session.Get(unlockSessionName, c)
	if err != nil || sess == nil || sess.Values == nil {
		return false
	}

	sessionUID, ok := sess.Values["uid"].(string)
	if !ok || sessionUID != uid {
		return false
	}

	unlocked, ok := sess.Values["unlocked"].(bool)
	return ok && unlocked
}

func unlockSessionOptions(maxAge int) *sessions.Options {
	return &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   config.IsProdEnv(),
		SameSite: http.SameSiteLaxMode,
	}
}
