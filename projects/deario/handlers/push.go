package handlers

import (
	"net/http"

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

	var data map[string]interface{}
	if err := c.Bind(&data); err != nil {
		return err
	}

	token := data["token"].(string)

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	if err := queries.UpsertPushKey(c.Request().Context(), db.UpsertPushKeyParams{
		Uid:       uid,
		PushToken: token,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "푸시 키 저장 실패")
	}

	return nil
}
