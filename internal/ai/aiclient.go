package aiclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
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

func Transcribe(ctx context.Context, r io.Reader) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.EnvMap["GEMINI_AI_KEY"],
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("AI 클라이언트 생성 실패: %w", err)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("오디오 읽기 실패: %w", err)
	}

	parts := []*genai.Part{
		{Text: "다음 음성을 한국어 텍스트로 정확히 전사해줘."},
		{InlineData: &genai.Blob{Data: data, MIMEType: "audio/webm"}},
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", []*genai.Content{{Parts: parts}}, nil)
	if err != nil {
		return "", fmt.Errorf("오디오 인식 실패: %w", err)
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("오디오 인식 실패: 빈 응답")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
