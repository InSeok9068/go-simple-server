package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"simple-server/internal/config"

	"github.com/labstack/echo/v4"
	_ "github.com/openai/openai-go"
	"google.golang.org/genai"
)

func AIStudy(c echo.Context, random bool) error {
	ctx := c.Request().Context()
	input := c.Request().FormValue("input")
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.EnvMap["GEMINI_AI_KEY"],
		Backend: genai.BackendGeminiAPI,
	})

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

	slog.Info(fmt.Sprintf(`프롬프트 요청 : %s`, prompt))

	result, err := client.Models.GenerateContent(ctx, "gemini-1.5-flash", genai.Text(prompt), nil)
	if err != nil {
		slog.Error("AI 요청 실패", "error", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "AI 요청 실패")
	}

	resultText := result.Candidates[0].Content.Parts[0].Text
	re := regexp.MustCompile(`(?s)<ol>.*?</ol>`)
	resultText = re.FindString(resultText)

	slog.Info(fmt.Sprintf(`프롬프트 응답 : %s`, resultText))

	return c.HTML(http.StatusOK, resultText)
}
