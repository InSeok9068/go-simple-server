package wardrobe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"simple-server/internal/config"

	"google.golang.org/genai"
)

const imageAnalysisModel = "gemini-2.5-flash"

// ImageMetadata는 AI가 추출한 자동 메타데이터다.
type ImageMetadata struct {
	Summary string   `json:"summary"`
	Season  string   `json:"season"`
	Style   string   `json:"style"`
	Colors  []string `json:"colors"`
	Tags    []string `json:"tags"`
}

// AnalyzeClosetImage는 Gemini 2.5 Flash로 이미지를 분석해 요약/태그 정보를 반환한다.
func AnalyzeClosetImage(ctx context.Context, data []byte, mimeType string) (*ImageMetadata, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("이미지 데이터가 비어 있어요")
	}
	apiKey := config.EnvMap["GEMINI_AI_KEY"]
	if apiKey == "" {
		return nil, fmt.Errorf("gemini api 키가 비어 있습니다")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("gemini 클라이언트 생성 실패: %w", err)
	}

	prompt := `
당신은 옷장을 관리하는 패션 어시스턴트입니다.
입력된 이미지를 보고 아래 JSON 형식으로만 답변하세요.
{
  "summary": "30자 이하 한국어 설명",
  "season": "추천 계절 (예: 봄/여름/가을/겨울)",
  "style": "착장 스타일 (예: 캐주얼, 포멀 등)",
  "colors": ["주요 색상명", "..."],
  "tags": ["아이템 종류나 핵심 키워드", "..."]
}
`

	parts := []*genai.Part{
		{Text: prompt},
		{InlineData: &genai.Blob{Data: data, MIMEType: mimeType}},
	}
	contents := []*genai.Content{
		{Parts: parts},
	}

	resp, err := client.Models.GenerateContent(ctx, imageAnalysisModel, contents, &genai.GenerateContentConfig{
		ResponseModalities: []string{"Text"},
	})
	if err != nil {
		return nil, fmt.Errorf("이미지 분석 요청 실패: %w", err)
	}
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("이미지 분석 응답이 비어 있어요")
	}

	raw := strings.TrimSpace(resp.Candidates[0].Content.Parts[0].Text)
	clean := sanitizeJSONBlock(raw)

	meta := &ImageMetadata{}
	if err := json.Unmarshal([]byte(clean), meta); err != nil {
		return nil, fmt.Errorf("이미지 메타데이터 파싱 실패: %w", err)
	}
	meta.normalize()

	return meta, nil
}

func (m *ImageMetadata) normalize() {
	m.Summary = strings.TrimSpace(m.Summary)
	m.Season = normalizeWord(m.Season)
	m.Style = normalizeWord(m.Style)
	m.Colors = normalizeSlice(m.Colors)
	m.Tags = normalizeSlice(m.Tags)
}

func normalizeWord(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func normalizeSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, v := range values {
		normalized := normalizeWord(v)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}

func sanitizeJSONBlock(text string) string {
	trimmed := strings.TrimSpace(text)
	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```JSON")
	trimmed = strings.TrimSuffix(trimmed, "```")
	return strings.TrimSpace(trimmed)
}
