package notification

import (
	"net/http"

	"simple-server/internal/validate"
	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"

	"github.com/labstack/echo/v4"
)

// RegisterPushToken은 푸시 토큰을 등록한다.
func RegisterPushToken(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	type registerPushTokenDTO struct {
		Token string `json:"token" validate:"required"`
	}
	var dto registerPushTokenDTO
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

	if err := queries.UpsertPushKey(c.Request().Context(), db.UpsertPushKeyParams{
		Uid:       uid,
		PushToken: dto.Token,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "푸시 키 저장 실패")
	}

	return c.NoContent(http.StatusNoContent)
}
