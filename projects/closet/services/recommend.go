package services

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"simple-server/internal/config"
	"simple-server/projects/closet/db"

	"google.golang.org/genai"
)

type RecommendationResult struct {
	Kind  string
	Item  db.ListItemsByIDsRow
	Score float64
}

const (
	recommendTopK = 1
)

func RecommendOutfit(ctx context.Context, weather, style string) ([]RecommendationResult, error) {
	query := strings.TrimSpace(weather + " " + style)
	if query == "" {
		return nil, errors.New("조건을 입력해주세요")
	}

	queryVec, err := textEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	queries, err := db.GetQueries()
	if err != nil {
		return nil, err
	}

	embeddings, err := queries.ListEmbeddingItems(ctx)
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, errors.New("아직 임베딩 정보가 없습니다")
	}

	bestByKind := rankByKind(queryVec, embeddings)

	selectedIDs := collectIDs(bestByKind)
	if len(selectedIDs) == 0 {
		return nil, errors.New("조건에 맞는 추천을 찾지 못했습니다")
	}

	items, err := queries.ListItemsByIDs(ctx, selectedIDs)
	if err != nil {
		return nil, err
	}
	itemMap := make(map[int64]db.ListItemsByIDsRow, len(items))
	for _, row := range items {
		itemMap[row.ID] = row
	}

	results := buildResults(bestByKind, itemMap)
	if len(results) == 0 {
		return nil, errors.New("추천 결과를 구성하지 못했습니다")
	}
	return results, nil
}

func rankByKind(queryVec []float32, embeddings []db.ListEmbeddingItemsRow) map[string][]scoreItem {
	bestByKind := make(map[string][]scoreItem)
	for _, row := range embeddings {
		vec := bytesToFloat32Slice(row.VecF32)
		if len(vec) == 0 || len(vec) != len(queryVec) {
			continue
		}
		score := cosineSimilarity(queryVec, vec)
		list := bestByKind[row.Kind]
		list = append(list, scoreItem{ItemID: row.ID, Score: score})
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Score > list[j].Score
		})
		if len(list) > recommendTopK {
			list = list[:recommendTopK]
		}
		bestByKind[row.Kind] = list
	}
	return bestByKind
}

func collectIDs(best map[string][]scoreItem) []int64 {
	var ids []int64
	for _, kind := range []string{"top", "bottom", "shoes", "accessory"} {
		for _, item := range best[kind] {
			ids = append(ids, item.ItemID)
		}
	}
	return ids
}

func buildResults(best map[string][]scoreItem, items map[int64]db.ListItemsByIDsRow) []RecommendationResult {
	results := make([]RecommendationResult, 0, len(items))
	for _, kind := range []string{"top", "bottom", "shoes", "accessory"} {
		for _, candidate := range best[kind] {
			if item, ok := items[candidate.ItemID]; ok {
				results = append(results, RecommendationResult{
					Kind:  kind,
					Item:  item,
					Score: candidate.Score,
				})
			}
		}
	}
	return results
}

type scoreItem struct {
	ItemID int64
	Score  float64
}

func textEmbedding(ctx context.Context, text string) ([]float32, error) {
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

	resp, err := client.Models.EmbedContent(ctx, embeddingModelName, genai.Text(text), &genai.EmbedContentConfig{
		TaskType: "RETRIEVAL_QUERY",
	})
	if err != nil {
		return nil, fmt.Errorf("임베딩 요청 실패: %w", err)
	}
	if len(resp.Embeddings) == 0 || len(resp.Embeddings[0].Values) == 0 {
		return nil, fmt.Errorf("임베딩 응답이 비었습니다")
	}
	return resp.Embeddings[0].Values, nil
}

func bytesToFloat32Slice(data []byte) []float32 {
	if len(data)%4 != 0 {
		return nil
	}
	vec := make([]float32, len(data)/4)
	for i := range vec {
		bits := binary.LittleEndian.Uint32(data[i*4 : (i+1)*4])
		vec[i] = math.Float32frombits(bits)
	}
	return vec
}

func cosineSimilarity(a, b []float32) float64 {
	var dot float64
	var normA float64
	var normB float64
	for i := range a {
		valA := float64(a[i])
		valB := float64(b[i])
		dot += valA * valB
		normA += valA * valA
		normB += valB * valB
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
