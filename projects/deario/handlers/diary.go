package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"simple-server/pkg/util/authutil"
	"simple-server/pkg/util/dateutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views"

	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Index는 메인 페이지를 렌더링한다.
func Index(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("20060102")
	}
	date = strings.ReplaceAll(date, "-", "")

	uid, _ := authutil.SessionUID(c)

	if uid == "" {
		return views.Index(os.Getenv("APP_TITLE"), date, "0").Render(c.Response().Writer)
	}

	queries, err := db.GetQueries(c.Request().Context())
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

	mood := moodValue(diary, errDiary)

	return views.Index(os.Getenv("APP_TITLE"), date, mood).Render(c.Response().Writer)
}

// Diary는 특정 날짜의 일기를 조회한다.
func Diary(c echo.Context) error {
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

	queries, err := db.GetQueries(c.Request().Context())
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

// DiaryList는 일기 목록을 반환한다.
func DiaryList(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}

	queries, err := db.GetQueries(c.Request().Context())
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

// DiaryRandom은 무작위 일기 날짜로 이동한다.
func DiaryRandom(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	queries, err := db.GetQueries(c.Request().Context())
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

	return c.HTML(http.StatusOK, fmt.Sprintf(`<script>location.href = "/?date=%s";</script>`, diary.Date))
}

// Save는 일기를 저장하거나 수정한다.
func Save(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := c.FormValue("date")
	content := c.FormValue("content")
	if date == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "날짜는 필수 입력값입니다.")
	}

	queries, err := db.GetQueries(c.Request().Context())
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

	return nil
}
