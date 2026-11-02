package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterErrorHandler : 에코 전역 에러 핸들러 등록
func RegisterErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var (
			he      *echo.HTTPError
			code    int
			message string
		)

		if errors.As(err, &he) {
			code = he.Code
			if m, ok := he.Message.(string); ok {
				message = m
			} else {
				message = fmt.Sprint(he.Message)
			}
		} else {
			code = http.StatusInternalServerError
			message = http.StatusText(code)
		}

		var internalErr error
		if he != nil {
			internalErr = he.Internal
		}
		sessionCookieCount := 0
		for _, cookie := range c.Request().Cookies() {
			if cookie.Name == "session_v2" {
				sessionCookieCount++
			}
		}

		slog.ErrorContext(c.Request().Context(), "요청 처리 실패",
			"error", err,
			"status", code,
			"method", c.Request().Method,
			"path", c.Path(),
			"내부_오류", internalErr,
			"접속_호스트", c.Request().Host,
			"브라우저", c.Request().UserAgent(),
			"로그인_쿠키_개수", sessionCookieCount,
		)

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				_ = c.NoContent(code)
			} else {
				_ = c.JSON(code, map[string]string{
					"message": message,
				})
			}
		}
	}
}
