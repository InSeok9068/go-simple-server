package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

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

		logArgs := []any{
			"error", err,
			"status", code,
			"method", c.Request().Method,
			"path", c.Path(),
		}
		if code >= http.StatusInternalServerError {
			logArgs = append(logArgs, "stack", string(debug.Stack()))
		}
		slog.ErrorContext(c.Request().Context(), "요청 처리 실패", logArgs...)

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
