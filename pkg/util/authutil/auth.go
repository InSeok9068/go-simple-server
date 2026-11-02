package authutil

import (
	"errors"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SessionUID(c echo.Context) (string, error) {
	sess, err := session.Get("session_v2", c)
	if err != nil {
		var cookieErr securecookie.Error
		if errors.As(err, &cookieErr) && cookieErr.IsDecode() {
			return "", echo.NewHTTPError(http.StatusUnauthorized, "로그인 세션이 만료되었거나 유효하지 않습니다. 다시 로그인해주세요.").SetInternal(err)
		}
		return "", echo.NewHTTPError(http.StatusInternalServerError, "세션을 가져오는 중 오류가 발생했습니다.").SetInternal(err)
	}

	if sess == nil || sess.Values == nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 세션입니다.")
	}

	uid := sess.Values["uid"]
	if uid == nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
	}

	strUID, ok := uid.(string)
	if !ok || strUID == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
	}

	return strUID, nil
}
