package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterErrorHandler : 에코 전역 에러 핸들러 등록
func RegisterErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			code    int
			message string
		)

		if he, ok := err.(*echo.HTTPError); ok {
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

		slog.Error("요청 처리 실패",
			"error", err,
			"status", code,
			"method", c.Request().Method,
			"path", c.Path(),
		)

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				_ = c.NoContent(code)
			} else {
				_ = c.String(code, message)
			}
		}
	}
}
