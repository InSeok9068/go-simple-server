package handlers

import (
	"net/http"

	"simple-server/pkg/util/authutil"
	"simple-server/pkg/util/maputil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views/pages"

	"github.com/labstack/echo/v4"
)

// SettingsPage는 사용자 설정 페이지를 렌더링한다.
func SettingsPage(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	userSetting, err := queries.GetUserSetting(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "사용자 설정을 가져오지 못했습니다.")
	}

	return pages.Setting(userSetting).Render(c.Request().Context(), c.Response().Writer)
}

// UpdateSettings는 사용자 설정을 저장한다.
func UpdateSettings(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := c.Bind(&data); err != nil {
		return err
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	if err := queries.UpsertUserSetting(c.Request().Context(), db.UpsertUserSettingParams{
		Uid:         uid,
		IsPush:      maputil.GetInt64(data, "is_push", 0),
		PushTime:    maputil.GetString(data, "push_time", ""),
		RandomRange: maputil.GetInt64(data, "random_range", 365),
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "사용자 설정 저장 실패")
	}

	return c.NoContent(http.StatusNoContent)
}
