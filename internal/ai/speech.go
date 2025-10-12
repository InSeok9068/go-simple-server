package aiclient

import (
	"context"
	"fmt"
	"net/http"
	"simple-server/internal/config"

	"google.golang.org/genai"
)

// TranscribeAudio는 음성 데이터를 텍스트로 변환한다.
func TranscribeAudio(ctx context.Context, data []byte, mimeType string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("오디오 데이터가 비어 있습니다")
	}
	if config.EnvMap["GEMINI_AI_KEY"] == "" {
		return "", fmt.Errorf("AI 키 설정을 찾을 수 없습니다")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.EnvMap["GEMINI_AI_KEY"],
		Backend: genai.BackendGeminiAPI,
		HTTPClient: &http.Client{
			Transport: http.DefaultTransport,
		},
	})
	if err != nil {
		return "", fmt.Errorf("AI 클라이언트 생성 실패: %w", err)
	}
	parts := []*genai.Part{
		{Text: "음성 데이터 내용을 분석하지 말고 그대로 한국어 텍스트로 변환해줘"},
		{InlineData: &genai.Blob{Data: data, MIMEType: mimeType}},
	}
	contents := []*genai.Content{{Parts: parts}}
	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, &genai.GenerateContentConfig{
		ResponseModalities: []string{"Text"},
	})
	if err != nil {
		return "", fmt.Errorf("오디오 변환 요청 실패: %w", err)
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("응답이 비어 있습니다")
	}
	return result.Candidates[0].Content.Parts[0].Text, nil
}
