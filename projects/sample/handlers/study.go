package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"simple-server/internal/config"

	"github.com/labstack/echo/v4"
	"google.golang.org/genai"
)

func AIStudy(c echo.Context, random bool) error {
	ctx := c.Request().Context()
	input := c.Request().FormValue("input")
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.EnvMap["GEMINI_AI_KEY"],
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("AI 클라이언트 생성 실패", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "AI 초기화 실패")
	}

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

	result, err := client.Models.GenerateContent(ctx, "gemini-1.5-flash", genai.Text(prompt), nil)
	if err != nil {
		slog.Error("AI 요청 실패", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "AI 요청 실패")
	}

	resultText := result.Candidates[0].Content.Parts[0].Text
	re := regexp.MustCompile(`(?s)<ol>.*?</ol>`)
	resultText = re.FindString(resultText)

	return c.HTML(http.StatusOK, resultText)
}
