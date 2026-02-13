package diary

import (
	"net/http"

	"simple-server/internal/validate"
	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"

	"github.com/labstack/echo/v4"
)

type updateDiaryMoodDTO struct {
	Date string `form:"date" json:"date" validate:"required,len=8,numeric" message:"날짜 형식이 올바르지 않습니다."`
	Mood string `form:"mood" json:"mood" validate:"required,oneof=0 1 2 3 4 5" message:"기분 값이 올바르지 않습니다."`
}

// UpdateDiaryMood는 일기의 기분 정보를 저장한다.
func UpdateDiaryMood(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	var dto updateDiaryMoodDTO
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

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: dto.Date,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "일기를 먼저 작성해주세요.")
	}

	if err := queries.UpdateDiaryOfMood(c.Request().Context(), db.UpdateDiaryOfMoodParams{
		ID:   diary.ID,
		Mood: dto.Mood,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "일기요정 저장 실패")
	}

	return c.NoContent(http.StatusNoContent)
}
