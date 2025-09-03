package aiclient

import (
	"context"
	"fmt"
	"simple-server/internal/config"

	"google.golang.org/genai"
)

// TranscribeAudio는 음성 데이터를 텍스트로 변환한다.
func TranscribeAudio(ctx context.Context, data []byte, mimeType string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("오디오 데이터가 비어있습니다")
	}
	if config.EnvMap["GEMINI_AI_KEY"] == "" {
		return "", fmt.Errorf("AI 키가 설정되지 않았습니다")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.EnvMap["GEMINI_AI_KEY"],
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("AI 클라이언트 생성 실패: %w", err)
	}
	parts := []*genai.Part{
		{Text: "다음 오디오 내용을 한국어 텍스트로 변환해줘."},
		{InlineData: &genai.Blob{Data: data, MIMEType: mimeType}},
	}
	contents := []*genai.Content{{Parts: parts}}
	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", contents, nil)
	if err != nil {
		return "", fmt.Errorf("오디오 인식 요청 실패: %w", err)
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("응답이 비어있습니다")
	}
	return result.Candidates[0].Content.Parts[0].Text, nil
}
