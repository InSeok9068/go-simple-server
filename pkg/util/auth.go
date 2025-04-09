package util

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SesseionUid(c echo.Context) (string, error) {
	sess, _ := session.Get("session", c)
	uid := sess.Values["uid"]
	if uid == nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
	}

	strUid := uid.(string)

	if strUid == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
	}

	return strUid, nil
}
