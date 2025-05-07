package util

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SesseionUID(c echo.Context) (string, error) {
	sess, _ := session.Get("session", c)
	uid := sess.Values["uid"]
	if uid == nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
	}

	strUID := uid.(string)

	if strUID == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
	}

	return strUID, nil
}
