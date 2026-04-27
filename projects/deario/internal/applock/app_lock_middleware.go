package applock

import (
	"net/http"

	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/views/pages"

	"github.com/labstack/echo/v4"
)

// GuardIndex는 로그인하지 않은 사용자는 통과시키고, 잠금 대상 사용자만 잠금 화면으로 보낸다.
func GuardIndex(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid, err := authutil.SessionUID(c)
		if err != nil {
			return next(c)
		}

		locked, err := isLocked(c, uid)
		if err != nil {
			return err
		}
		if locked {
			return pages.AppLock().Render(c.Request().Context(), c.Response().Writer)
		}

		return next(c)
	}
}

// RequireUnlocked는 로그인과 앱 잠금 해제를 모두 요구한다.
func RequireUnlocked(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid, err := authutil.SessionUID(c)
		if err != nil {
			return err
		}

		locked, err := isLocked(c, uid)
		if err != nil {
			return err
		}
		if locked {
			return lockedResponse(c)
		}

		return next(c)
	}
}

func isLocked(c echo.Context, uid string) (bool, error) {
	setting, err := userSetting(c, uid)
	if err != nil {
		return false, err
	}
	if setting.AppLockEnabled != 1 {
		return false, nil
	}

	return !isUnlockSessionValid(c, uid), nil
}

func lockedResponse(c echo.Context) error {
	if c.Request().Header.Get("Hx-Request") == "true" {
		c.Response().Header().Set("Hx-Redirect", "/app-lock")
		return c.NoContent(http.StatusNoContent)
	}
	if c.Request().Method == http.MethodGet {
		return pages.AppLock().Render(c.Request().Context(), c.Response().Writer)
	}

	return echo.NewHTTPError(http.StatusLocked, "PIN 잠금 해제가 필요합니다.")
}
