package aiclient

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/genai"
	"log/slog"
	"simple-server/internal/config"
)

func Request(ctx context.Context, prompt string) (string, error) {
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.EnvMap["GEMINI_AI_KEY"],
		Backend: genai.BackendGeminiAPI,
	})

	slog.Info(fmt.Sprintf(`프롬프트 요청 : %s`, prompt))

	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", genai.Text(prompt), nil)
	if err != nil {
		slog.Error("AI 요청 실패", "error", err.Error())
		return "", errors.New("AI 요청 실패")
	}

	resultText := result.Candidates[0].Content.Parts[0].Text

	slog.Info(fmt.Sprintf(`프롬프트 응답 : %s`, resultText))

	return resultText, nil
}
