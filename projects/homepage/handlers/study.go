package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"simple-server/internal"

	"github.com/labstack/echo/v4"
	"google.golang.org/genai"
)

func AIStudy(c echo.Context) error {
	ctx := c.Request().Context()
	input := c.Request().FormValue("input")
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  internal.EnvMap["GEMINI_AI_KEY"],
		Backend: genai.BackendGoogleAI,
	})

	prompt := fmt.Sprintf(`
	해당 주제로 공부할 주제를 짧게 작성해줘
	주제 : %s
	<ol>
		<li>{주제}</li>
		<li>{주제}</li>
		<li>{주제}</li>
		....
	</ol>
	`, input)

	slog.Info(fmt.Sprintf(`프롬프트 요청 : %s`, prompt))

	result, err := client.Models.GenerateContent(ctx, "gemini-1.5-flash", genai.Text(prompt), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "목록 조회 오류")
	}

	slog.Info(fmt.Sprintf(`프롬프트 응답 : %s`, result.Candidates[0].Content.Parts[0].Text))

	return c.HTML(http.StatusOK, result.Candidates[0].Content.Parts[0].Text)
}
