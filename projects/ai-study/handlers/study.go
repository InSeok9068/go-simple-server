package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	aiclient "simple-server/internal/ai"

	"github.com/labstack/echo/v4"
)

func AIStudy(c echo.Context, random bool) error {
	ctx := c.Request().Context()
	input := c.Request().FormValue("input")

	if random {
		input = "너가 정해줘"
	}

	prompt := fmt.Sprintf(`
	해당 주제로 공부할 주제를 짧게 10개 작성해줘

	주제 : %s

	output :
	<ol>
		<li>{주제}</li>
		<li>{주제}</li>
		<li>{주제}</li>
		....
	</ol>
	`, input)

	result, err := aiclient.Request(ctx, prompt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	re := regexp.MustCompile(`(?s)<ol>.*?</ol>`)
	result = re.FindString(result)

	return c.HTML(http.StatusOK, result)
}
