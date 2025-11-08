package services

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
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

var kindOrder = []string{"top", "bottom", "shoes", "accessory"}

func RecommendOutfit(ctx context.Context, uid string, weather, style, skipRaw string, locks map[string]int64) ([]RecommendationResult, string, bool, error) {
	if strings.TrimSpace(uid) == "" {
		return nil, skipRaw, false, errors.New("로그인이 필요해요.")
	}

	query := strings.TrimSpace(weather + " " + style)
	if query == "" {
		return nil, skipRaw, false, errors.New("조건을 입력해주세요.")
	}

	preference := derivePreference(weather, style)
	skipSet := parseSkipIDs(skipRaw)
	locks = normalizeLocks(locks)
	for _, id := range locks {
		skipSet[id] = struct{}{}
	}

	queryVec, err := textEmbedding(ctx, query)
	if err != nil {
		return nil, skipRaw, false, err
	}

	queries, err := db.GetQueries()
	if err != nil {
		return nil, skipRaw, false, err
	}

	embeddings, err := queries.ListEmbeddingItems(ctx, uid)
	if err != nil {
		return nil, skipRaw, false, err
	}
	if len(embeddings) == 0 {
		return nil, skipRaw, false, errors.New("아직 추천에 활용할 데이터가 없어요.")
	}

	filtered := filterEmbeddings(embeddings, preference)
	bestByKind := rankByKind(queryVec, filtered, preference)
	if hasMissingKinds(bestByKind) {
		fallback := rankByKind(queryVec, embeddings, preference)
		mergeMissingKinds(bestByKind, fallback)
	}

	selected, hasMore := pickCandidates(bestByKind, skipSet, locks)
	selectedIDs := collectIDs(selected)
	if len(selectedIDs) == 0 {
		return nil, skipRaw, false, errors.New("조건에 맞는 추천을 찾지 못했어요.")
	}

	items, err := queries.ListItemsByIDs(ctx, db.ListItemsByIDsParams{
		UserUid: uid,
		Ids:     selectedIDs,
	})
	if err != nil {
		return nil, skipRaw, false, err
	}
	itemMap := make(map[int64]db.ListItemsByIDsRow, len(items))
	for _, row := range items {
		itemMap[row.ID] = row
	}

	results := buildResults(selected, itemMap)
	if len(results) == 0 {
		return nil, skipRaw, false, errors.New("추천 결과를 구성하지 못했어요.")
	}
	nextCache := formatSkipToken(skipSet)
	return results, nextCache, hasMore, nil
}

func rankByKind(queryVec []float32, embeddings []db.ListEmbeddingItemsRow, pref metadataPreference) map[string][]scoreItem {
	bestByKind := make(map[string][]scoreItem)
	for _, row := range embeddings {
		vec := bytesToFloat32Slice(row.VecF32)
		if len(vec) == 0 || len(vec) != len(queryVec) {
			continue
		}
		score := cosineSimilarity(queryVec, vec) + metadataBoost(row, pref)
		list := bestByKind[row.Kind]
		list = append(list, scoreItem{ItemID: row.ID, Score: score})
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Score > list[j].Score
		})
		bestByKind[row.Kind] = list
	}
	return bestByKind
}

func collectIDs(selected map[string]scoreItem) []int64 {
	var ids []int64
	for _, kind := range kindOrder {
		if item, ok := selected[kind]; ok {
			ids = append(ids, item.ItemID)
		}
	}
	return ids
}

func buildResults(selected map[string]scoreItem, items map[int64]db.ListItemsByIDsRow) []RecommendationResult {
	results := make([]RecommendationResult, 0, len(selected))
	for _, kind := range kindOrder {
		if candidate, ok := selected[kind]; ok {
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

func pickCandidates(best map[string][]scoreItem, skip map[int64]struct{}, locks map[string]int64) (map[string]scoreItem, bool) {
	selected := make(map[string]scoreItem)

	for kind, id := range locks {
		if id <= 0 {
			continue
		}
		if candidate, ok := findCandidateByID(best[kind], id); ok {
			selected[kind] = candidate
		} else {
			selected[kind] = scoreItem{ItemID: id, Score: math.MaxFloat32}
		}
		skip[id] = struct{}{}
	}

	for _, kind := range kindOrder {
		if _, locked := locks[kind]; locked {
			continue
		}
		list := best[kind]
		for _, candidate := range list {
			if _, exists := skip[candidate.ItemID]; exists {
				continue
			}
			selected[kind] = candidate
			skip[candidate.ItemID] = struct{}{}
			break
		}
	}

	hasMore := false
kindLoop:
	for _, kind := range kindOrder {
		if _, locked := locks[kind]; locked {
			continue
		}
		list := best[kind]
		for _, candidate := range list {
			if _, exists := skip[candidate.ItemID]; !exists {
				hasMore = true
				break kindLoop
			}
		}
	}
	return selected, hasMore
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
		return nil, fmt.Errorf("임베딩 응답이 비어 있어요")
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

// preference & metadata helpers ------------------------------------------------

type metadataPreference struct {
	seasons []string
	styles  []string
}

func derivePreference(weather, desiredStyle string) metadataPreference {
	pref := metadataPreference{}
	pref.seasons = appendUnique(pref.seasons, detectSeasons(weather)...)
	pref.seasons = appendUnique(pref.seasons, detectSeasons(desiredStyle)...)
	pref.styles = appendUnique(pref.styles, detectStyles(weather)...)
	pref.styles = appendUnique(pref.styles, detectStyles(desiredStyle)...)
	return pref
}

func filterEmbeddings(rows []db.ListEmbeddingItemsRow, pref metadataPreference) []db.ListEmbeddingItemsRow {
	if len(pref.seasons) == 0 && len(pref.styles) == 0 {
		return rows
	}
	filtered := make([]db.ListEmbeddingItemsRow, 0, len(rows))
	for _, row := range rows {
		if matchesPreference(row, pref) {
			filtered = append(filtered, row)
		}
	}
	if len(filtered) == 0 {
		return rows
	}
	return filtered
}

func matchesPreference(row db.ListEmbeddingItemsRow, pref metadataPreference) bool {
	if len(pref.seasons) > 0 && row.MetaSeason.Valid {
		if !containsAny(row.MetaSeason.String, pref.seasons) {
			return false
		}
	}
	if len(pref.styles) > 0 && row.MetaStyle.Valid {
		if !containsAny(row.MetaStyle.String, pref.styles) {
			return false
		}
	}
	return true
}

func metadataBoost(row db.ListEmbeddingItemsRow, pref metadataPreference) float64 {
	boost := 0.0
	if len(pref.seasons) > 0 && row.MetaSeason.Valid && containsAny(row.MetaSeason.String, pref.seasons) {
		boost += 0.15
	}
	if len(pref.styles) > 0 && row.MetaStyle.Valid && containsAny(row.MetaStyle.String, pref.styles) {
		boost += 0.1
	}
	return boost
}

func detectSeasons(text string) []string {
	lower := strings.ToLower(text)
	seasons := map[string][]string{
		"봄":  {"봄", "spring"},
		"여름": {"여름", "summer"},
		"가을": {"가을", "autumn", "fall"},
		"겨울": {"겨울", "winter"},
	}
	return findKeywords(lower, seasons)
}

func detectStyles(text string) []string {
	lower := strings.ToLower(text)
	styles := map[string][]string{
		"캐주얼":  {"캐주얼", "casual"},
		"포멀":   {"포멀", "정장", "formal"},
		"스트릿":  {"스트릿", "street"},
		"스포츠":  {"스포츠", "athleisure", "운동"},
		"미니멀":  {"미니멀", "minimal"},
		"로맨틱":  {"로맨틱", "페미닌", "feminine", "romantic"},
		"비즈니스": {"비즈니스", "오피스", "office"},
	}
	return findKeywords(lower, styles)
}

func findKeywords(input string, keywords map[string][]string) []string {
	var results []string
	for canonical, variants := range keywords {
		for _, keyword := range variants {
			if keyword == "" {
				continue
			}
			if strings.Contains(input, keyword) {
				results = appendUnique(results, canonical)
				break
			}
		}
	}
	return results
}

func appendUnique(list []string, values ...string) []string {
	seen := make(map[string]struct{}, len(list))
	for _, v := range list {
		seen[v] = struct{}{}
	}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		list = append(list, value)
	}
	return list
}

func containsAny(value string, candidates []string) bool {
	if len(candidates) == 0 {
		return true
	}
	valueTokens := splitTokens(value)
	for _, token := range valueTokens {
		for _, candidate := range candidates {
			if token == candidate {
				return true
			}
		}
	}
	return false
}

func splitTokens(value string) []string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return nil
	}
	separators := func(r rune) bool {
		return r == ',' || r == '/' || r == ' ' || r == '|' || r == '\n'
	}
	raw := strings.FieldsFunc(value, separators)
	result := make([]string, 0, len(raw))
	for _, token := range raw {
		token = strings.TrimSpace(token)
		if token != "" {
			result = append(result, token)
		}
	}
	return result
}

func parseSkipIDs(raw string) map[int64]struct{} {
	skip := make(map[int64]struct{})
	if strings.TrimSpace(raw) == "" {
		return skip
	}
	parts := strings.Split(raw, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil || id <= 0 {
			continue
		}
		skip[id] = struct{}{}
	}
	return skip
}

func formatSkipToken(skip map[int64]struct{}) string {
	if len(skip) == 0 {
		return ""
	}
	ids := make([]int64, 0, len(skip))
	for id := range skip {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.FormatInt(id, 10))
	}
	return strings.Join(parts, ",")
}

func findCandidateByID(list []scoreItem, id int64) (scoreItem, bool) {
	for _, candidate := range list {
		if candidate.ItemID == id {
			return candidate, true
		}
	}
	return scoreItem{}, false
}

func normalizeLocks(input map[string]int64) map[string]int64 {
	if len(input) == 0 {
		return nil
	}
	normalized := make(map[string]int64, len(input))
	for _, kind := range kindOrder {
		if id, ok := input[kind]; ok && id > 0 {
			normalized[kind] = id
		}
	}
	return normalized
}

func hasMissingKinds(best map[string][]scoreItem) bool {
	required := []string{"top", "bottom"}
	for _, kind := range required {
		if len(best[kind]) == 0 {
			return true
		}
	}
	return false
}

func mergeMissingKinds(base, fallback map[string][]scoreItem) {
	for _, kind := range kindOrder {
		if len(base[kind]) == 0 && len(fallback[kind]) > 0 {
			base[kind] = fallback[kind]
		}
	}
}
