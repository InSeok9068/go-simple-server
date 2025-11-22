package diary

import (
	"net/http"

	"simple-server/pkg/util/authutil"
	"simple-server/pkg/util/maputil"
	"simple-server/projects/deario/db"

	"github.com/labstack/echo/v4"
)

// UpdateDiaryMood는 일기의 기분 정보를 저장한다.
func UpdateDiaryMood(c echo.Context) error {
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

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: maputil.GetString(data, "date", ""),
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "일기를 먼저 작성해주세요.")
	}

	if err := queries.UpdateDiaryOfMood(c.Request().Context(), db.UpdateDiaryOfMoodParams{
		ID:   diary.ID,
		Mood: maputil.GetString(data, "mood", "0"),
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "일기요정 저장 실패")
	}

	return c.NoContent(http.StatusNoContent)
}

// diaryMood는 일기에서 기분 값을 추출한다.
func diaryMood(d db.Diary, err error) string {
	if err == nil && d.ID != "" {
		return d.Mood
	}
	return "0"
}
