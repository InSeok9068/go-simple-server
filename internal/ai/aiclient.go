package aiclient

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"simple-server/internal/config"

	"google.golang.org/genai"
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

func ImageRequest(ctx context.Context, prompt string) (string, error) {
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.EnvMap["GEMINI_AI_KEY"],
		Backend: genai.BackendGeminiAPI,
	})

	slog.Info(fmt.Sprintf(`프롬프트 요청 : %s`, prompt))

	result, _ := client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp-image-generation", genai.Text(prompt), &genai.GenerateContentConfig{
		ResponseModalities: []string{"Text", "Image"},
	})
	slog.Info("이미지 생성 응답", "결과", result)
	resultImage := result.Candidates[0].Content.Parts[0].InlineData.Data
	data := base64.StdEncoding.EncodeToString(resultImage)
	return data, nil

	// result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp-image-generation", genai.Text(prompt), nil)
	// if err != nil {
	// 	slog.Error("AI 요청 실패", "error", err.Error())
	// 	return nil, errors.New("AI 요청 실패")
	// }

	// resultImage := result.Candidates[0].Content.Parts[0].InlineData.Data

	// //slog.Info(fmt.Sprintf(`프롬프트 응답 : %s`, resultImage))
	// data, err := base64.StdEncoding.DecodeString(string(resultImage))

	// return data, nil
}
