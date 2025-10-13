package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"simple-server/internal/validate"
	"simple-server/pkg/util/authutil"
	"simple-server/pkg/util/dateutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views"

	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// IndexPage는 메인 페이지를 렌더링한다.
func IndexPage(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("20060102")
	}
	date = strings.ReplaceAll(date, "-", "")

	uid, _ := authutil.SessionUID(c)

	if uid == "" {
		return views.Index(os.Getenv("APP_TITLE"), date, "0").Render(c.Response().Writer)
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	diary, errDiary := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: date,
	})

	if errDiary != nil && !errors.Is(errDiary, sql.ErrNoRows) {
		return errDiary
	}

	mood := diaryMood(diary, errDiary)

	return views.Index(os.Getenv("APP_TITLE"), date, mood).Render(c.Response().Writer)
}

// GetDiary는 특정 날짜의 일기를 조회한다.
func GetDiary(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := c.FormValue("date")
	if date == "" {
		date = time.Now().Format("20060102")
	} else {
		date = strings.ReplaceAll(date, "-", "")
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: date,
	})

	if err != nil {
		return views.DiaryContentForm(date, "").Render(c.Response().Writer)
	}
	return views.DiaryContentForm(diary.Date, diary.Content).Render(c.Response().Writer)
}

// ListDiaries는 일기 날짜 목록을 반환한다.
func ListDiaries(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	diarys, err := queries.ListDiarys(c.Request().Context(), db.ListDiarysParams{
		Uid:     uid,
		Column2: page,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "목록을 가져오지 못했습니다.")
	}

	var lis []Node
	for _, diary := range diarys {
		lis = append(lis,
			Li(
				A(Href(fmt.Sprintf("/?date=%s", diary.Date)),
					Text(dateutil.MustFormatDateKorWithWeekDay(diary.Date)),
				),
			),
		)
	}

	return Group(lis).Render(c.Response().Writer)
}

// MonthlyDiaryDates는 특정 월에 작성한 일기 날짜 목록을 반환한다.
func MonthlyDiaryDates(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	month := c.QueryParam("month")
	if month == "" {
		month = time.Now().Format("200601")
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	dates, err := queries.ListDiaryDatesByMonth(c.Request().Context(), db.ListDiaryDatesByMonthParams{
		Uid:  uid,
		Date: month,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "목록을 가져오지 못했습니다.")
	}

	return c.JSON(http.StatusOK, dates)
}

// RedirectToRandomDiary는 무작위 일기 날짜로 이동한다.
func RedirectToRandomDiary(c echo.Context) error {
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
		return err
	}

	dateLimit := time.Now().AddDate(0, 0, -int(userSetting.RandomRange)).Format("20060102")

	diary, err := queries.GetDiaryRandom(c.Request().Context(), db.GetDiaryRandomParams{
		Uid:  uid,
		Date: dateLimit,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "작성한 일기장이 없습니다.")
	}

	c.Response().Header().Set("HX-Redirect", "/?date="+diary.Date)
	return c.NoContent(http.StatusNoContent)
}

// SaveDiary는 일기를 저장하거나 수정한다.
func SaveDiary(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	type saveDiaryDTO struct {
		Date    string `form:"date" validate:"required" message:"날짜가 필요합니다."`
		Content string `form:"content"`
	}
	var dto saveDiaryDTO
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "요청 본문이 올바르지 않습니다.")
	}
	if err := c.Validate(&dto); err != nil {
		return validate.HTTPError(err, &dto)
	}

	date := dto.Date
	content := dto.Content

	queries, err := db.GetQueries()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "시스템 오류가 발생했습니다.")
	}

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: date,
	})

	if err != nil {
		if _, err := queries.CreateDiary(c.Request().Context(), db.CreateDiaryParams{
			Uid:     uid,
			Content: content,
			Date:    date,
		}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "일기 저장에 실패했습니다. 다시 시도해주세요.")
		}
	} else {
		if content == "" {
			if err := queries.DeleteDiary(c.Request().Context(), diary.ID); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "수정 실패")
			}
		} else {
			if _, err := queries.UpdateDiary(c.Request().Context(), db.UpdateDiaryParams{
				Content: content,
				ID:      diary.ID,
			}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "수정 실패")
			}
		}
	}

	return c.NoContent(http.StatusNoContent)
}
