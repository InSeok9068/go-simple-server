package aiclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"simple-server/internal/config"

	"google.golang.org/genai"
)

func Request(ctx context.Context, prompt string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.EnvMap["GEMINI_AI_KEY"],
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("AI 클라이언트 생성 실패: %w", err)
	}

	// gemini-2.5-pro
	// gemini-2.5-flash
	// gemini-2.5-flash-lite
	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("AI 요청 실패: %w", err)
	}

	resultText := result.Candidates[0].Content.Parts[0].Text

	return resultText, nil
}

func ImageRequest(ctx context.Context, prompt string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.EnvMap["GEMINI_AI_KEY"],
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("AI 클라이언트 생성 실패: %w", err)
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-preview-image-generation", genai.Text(prompt), &genai.GenerateContentConfig{
		ResponseModalities: []string{"Text", "Image"},
	})
	if err != nil {
		return "", fmt.Errorf("이미지 생성 요청 실패: %w", err)
	}

	for _, part := range result.Candidates[0].Content.Parts {
		if part.InlineData != nil {
			return base64.StdEncoding.EncodeToString(part.InlineData.Data), nil
		}
	}

	return "", fmt.Errorf("이미지 생성 실패")
}
