package wardrobe

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"simple-server/pkg/util/authutil"
	"simple-server/projects/closet/views"

	"github.com/labstack/echo/v4"
)

// RecommendOutfitHandler는 날씨와 스타일 조건에 맞는 추천 다이얼로그를 렌더링한다.
func RecommendOutfitHandler(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	if err := c.Request().ParseForm(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "입력값을 확인해주세요.")
	}
	weather := strings.TrimSpace(c.FormValue("weather"))
	style := strings.TrimSpace(c.FormValue("style"))
	skipIDs := strings.TrimSpace(c.FormValue("skip_ids"))
	locks := parseLockSelections(c)

	results, cacheToken, hasMore, err := RecommendOutfit(c.Request().Context(), uid, weather, style, skipIDs, locks)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	viewResults := make([]views.RecommendationItem, 0, len(results))
	for _, result := range results {
		viewResults = append(viewResults, views.RecommendationItem{
			Kind: result.Kind,
			Item: views.NewClosetItem(convertIDsRow(result.Item)),
		})
	}

	rows := views.BuildRecommendationRows(viewResults)
	hasResults := len(viewResults) > 0

	var builder strings.Builder
	if err := views.RecommendationDialog(rows, weather, style, cacheToken, hasMore, locks, hasResults).Render(c.Request().Context(), &builder); err != nil {
		return err
	}
	return c.HTML(http.StatusOK, builder.String())
}

func parseLockSelections(c echo.Context) map[string]int64 {
	locks := make(map[string]int64)
	for _, kind := range kindOrder {
		field := fmt.Sprintf("lock_%s", kind)
		value := strings.TrimSpace(c.FormValue(field))
		if value == "" {
			continue
		}
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			continue
		}
		locks[kind] = id
	}
	if len(locks) == 0 {
		return nil
	}
	return locks
}
