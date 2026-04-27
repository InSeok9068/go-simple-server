package applock

import (
	"database/sql"
	"errors"
	"net/http"

	"simple-server/internal/validate"
	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views/pages"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type unlockDTO struct {
	PIN string `form:"pin" validate:"required,numeric,min=4,max=6" message:"PIN은 숫자 4~6자리로 입력해주세요."`
}

type updatePINDTO struct {
	CurrentPIN string `form:"current_pin" validate:"omitempty,numeric,min=4,max=6" message:"현재 PIN은 숫자 4~6자리로 입력해주세요."`
	PIN        string `form:"pin" validate:"required,numeric,min=4,max=6" message:"새 PIN은 숫자 4~6자리로 입력해주세요."`
	PINConfirm string `form:"pin_confirm" validate:"required,numeric,min=4,max=6" message:"PIN 확인은 숫자 4~6자리로 입력해주세요."`
}

type disableDTO struct {
	CurrentPIN string `form:"current_pin" validate:"required,numeric,min=4,max=6" message:"현재 PIN은 숫자 4~6자리로 입력해주세요."`
}

// LockPage는 앱 잠금 화면을 렌더링한다.
func LockPage(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	setting, err := userSetting(c, uid)
	if err != nil {
		return err
	}
	if setting.AppLockEnabled != 1 {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	if isUnlockSessionValid(c, uid) {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	return pages.AppLock().Render(c.Request().Context(), c.Response().Writer)
}

// Unlock은 PIN을 확인하고 현재 브라우저 세션에서 앱 잠금을 해제한다.
func Unlock(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	setting, err := userSetting(c, uid)
	if err != nil {
		return err
	}
	if setting.AppLockEnabled != 1 {
		return redirectAfterUnlock(c)
	}

	var dto unlockDTO
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "요청 본문이 올바르지 않습니다.")
	}
	if err := c.Validate(&dto); err != nil {
		return validate.HTTPError(err, &dto)
	}

	if !isPINValid(setting.AppLockPinHash, dto.PIN) {
		return echo.NewHTTPError(http.StatusUnauthorized, "PIN이 올바르지 않습니다.")
	}

	if err := SaveUnlockSession(c, uid); err != nil {
		return err
	}

	return redirectAfterUnlock(c)
}

// Lock은 현재 브라우저 세션의 앱 잠금 해제를 취소한다.
func Lock(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}
	setting, err := userSetting(c, uid)
	if err != nil {
		return err
	}
	if err := ClearUnlockSession(c); err != nil {
		return err
	}
	if setting.AppLockEnabled == 1 {
		c.Response().Header().Set("Hx-Redirect", "/app-lock")
	}
	return c.NoContent(http.StatusNoContent)
}

// UpdatePIN은 앱 잠금을 켜거나 PIN을 변경한다.
func UpdatePIN(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	var dto updatePINDTO
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "요청 본문이 올바르지 않습니다.")
	}
	if err := c.Validate(&dto); err != nil {
		return validate.HTTPError(err, &dto)
	}
	if dto.PIN != dto.PINConfirm {
		return echo.NewHTTPError(http.StatusBadRequest, "새 PIN과 PIN 확인이 일치하지 않습니다.")
	}

	setting, err := userSetting(c, uid)
	if err != nil {
		return err
	}
	if setting.AppLockEnabled == 1 && !isPINValid(setting.AppLockPinHash, dto.CurrentPIN) {
		return echo.NewHTTPError(http.StatusUnauthorized, "현재 PIN이 올바르지 않습니다.")
	}

	pinHash, err := hashPIN(dto.PIN)
	if err != nil {
		return err
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}
	if err := queries.UpdateAppLock(c.Request().Context(), db.UpdateAppLockParams{
		AppLockEnabled: 1,
		AppLockPinHash: pinHash,
		Uid:            uid,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "PIN 잠금 설정 저장에 실패했습니다.")
	}

	if err := SaveUnlockSession(c, uid); err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Redirect", "/setting")
	return c.NoContent(http.StatusNoContent)
}

// Disable은 PIN 확인 후 앱 잠금을 끈다.
func Disable(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	var dto disableDTO
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "요청 본문이 올바르지 않습니다.")
	}
	if err := c.Validate(&dto); err != nil {
		return validate.HTTPError(err, &dto)
	}

	setting, err := userSetting(c, uid)
	if err != nil {
		return err
	}
	if setting.AppLockEnabled == 1 && !isPINValid(setting.AppLockPinHash, dto.CurrentPIN) {
		return echo.NewHTTPError(http.StatusUnauthorized, "현재 PIN이 올바르지 않습니다.")
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}
	if err := queries.UpdateAppLock(c.Request().Context(), db.UpdateAppLockParams{
		AppLockEnabled: 0,
		AppLockPinHash: "",
		Uid:            uid,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "PIN 잠금 해제에 실패했습니다.")
	}

	if err := ClearUnlockSession(c); err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Redirect", "/setting")
	return c.NoContent(http.StatusNoContent)
}

func userSetting(c echo.Context, uid string) (db.UserSetting, error) {
	queries, err := db.GetQueries()
	if err != nil {
		return db.UserSetting{}, err
	}

	setting, err := queries.GetUserSetting(c.Request().Context(), uid)
	if err == nil {
		return setting, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return db.UserSetting{}, err
	}

	setting = db.UserSetting{
		Uid:         uid,
		IsPush:      0,
		PushTime:    "",
		RandomRange: 365,
	}
	if err := queries.UpsertUserSetting(c.Request().Context(), db.UpsertUserSettingParams{
		Uid:         setting.Uid,
		IsPush:      setting.IsPush,
		PushTime:    setting.PushTime,
		RandomRange: setting.RandomRange,
	}); err != nil {
		return db.UserSetting{}, err
	}

	return setting, nil
}

func isPINValid(pinHash string, pin string) bool {
	if pinHash == "" || pin == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(pinHash), []byte(pin)) == nil
}

func hashPIN(pin string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "PIN 저장 준비에 실패했습니다.")
	}
	return string(hashed), nil
}

func redirectAfterUnlock(c echo.Context) error {
	c.Response().Header().Set("Hx-Redirect", "/")
	return c.NoContent(http.StatusNoContent)
}
