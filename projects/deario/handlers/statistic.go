package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views/pages"

	"github.com/labstack/echo/v4"
)

// StatsPage는 통계 페이지를 렌더링한다.
func StatsPage(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	if _, err := queries.GetUserSetting(c.Request().Context(), uid); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	return pages.Statistic().Render(c.Request().Context(), c.Response().Writer)
}

// buildMoodMap은 월별 기분 데이터를 맵으로 변환한다.
func buildMoodMap(rows []db.MonthlyMoodCountRow) map[string]db.MonthlyMoodCountRow {
	m := make(map[string]db.MonthlyMoodCountRow)
	for _, r := range rows {
		m[r.Month] = r
	}
	return m
}

// nullFloat64ToInt64는 sql.NullFloat64를 int64로 변환한다.
func nullFloat64ToInt64(v sql.NullFloat64) int64 {
	if v.Valid {
		return int64(v.Float64)
	}
	return 0
}

// GetStatsData는 통계 데이터를 JSON으로 반환한다.
func GetStatsData(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	counts, err := queries.MonthlyDiaryCount(c.Request().Context(), uid)
	if err != nil {
		return err
	}

	moods, err := queries.MonthlyMoodCount(c.Request().Context(), uid)
	if err != nil {
		return err
	}

	moodMap := buildMoodMap(moods)

	var months []string
	var diaryCount []int64
	var mood1, mood2, mood3, mood4, mood5 []int64

	for _, cRow := range counts {
		month := cRow.Month
		months = append(months, month)
		diaryCount = append(diaryCount, cRow.Count)

		mm := moodMap[month]
		mood1 = append(mood1, nullFloat64ToInt64(mm.Mood1))
		mood2 = append(mood2, nullFloat64ToInt64(mm.Mood2))
		mood3 = append(mood3, nullFloat64ToInt64(mm.Mood3))
		mood4 = append(mood4, nullFloat64ToInt64(mm.Mood4))
		mood5 = append(mood5, nullFloat64ToInt64(mm.Mood5))
	}

	result := map[string]interface{}{
		"months":     months,
		"diaryCount": diaryCount,
		"mood1":      mood1,
		"mood2":      mood2,
		"mood3":      mood3,
		"mood4":      mood4,
		"mood5":      mood5,
	}

	return c.JSON(http.StatusOK, result)
}
