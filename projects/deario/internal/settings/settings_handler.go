package settings

import (
	"net/http"

	"simple-server/internal/validate"
	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views/pages"

	"github.com/labstack/echo/v4"
)

type updateSettingsDTO struct {
	IsPush      int64  `form:"is_push" json:"is_push" validate:"oneof=0 1" message:"알림 설정 값이 올바르지 않습니다."`
	PushTime    string `form:"push_time" json:"push_time" validate:"omitempty,datetime=15:04" message:"알림 시간이 올바르지 않습니다."`
	RandomRange *int64 `form:"random_range" json:"random_range" validate:"omitempty,min=0,max=3650" message:"랜덤일자 범위가 올바르지 않습니다."`
}

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

	var dto updateSettingsDTO
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "요청 본문이 올바르지 않습니다.")
	}
	if err := c.Validate(&dto); err != nil {
		return validate.HTTPError(err, &dto)
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	randomRange := int64(365)
	if dto.RandomRange != nil {
		randomRange = *dto.RandomRange
	}

	if err := queries.UpsertUserSetting(c.Request().Context(), db.UpsertUserSettingParams{
		Uid:         uid,
		IsPush:      dto.IsPush,
		PushTime:    dto.PushTime,
		RandomRange: randomRange,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "사용자 설정 저장 실패")
	}

	return c.NoContent(http.StatusNoContent)
}
