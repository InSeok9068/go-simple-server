package services

import (
	"context"
	"encoding/binary"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"simple-server/internal/config"
	"simple-server/projects/closet/db"

	"google.golang.org/genai"
)

const (
	// gemini-embedding-001? ?띿뒪???꾩슜 ?꾨쿋??紐⑤뜽?대떎.
	embeddingModelName = "gemini-embedding-001"
	embeddingTimeout   = time.Minute
)

// EnqueueEmbeddingJob? ?낅줈??吏곹썑 ?대?吏 ?꾨쿋???앹꽦??泥섎━?쒕떎.
// UploadItem?먯꽌 goroutine?쇰줈 ?몄텧?섎?濡? ?ш린?쒕뒗 ?숆린?곸쑝濡쒕쭔 ?ㅽ뻾?쒕떎.
func EnqueueEmbeddingJob(itemID int64, contextText string) {
	ctx, cancel := context.WithTimeout(context.Background(), embeddingTimeout)
	defer cancel()

	if err := ImageEmbedding(ctx, itemID, contextText); err != nil {
		slog.Error("?대?吏 ?꾨쿋???앹꽦 ?ㅽ뙣", "item_id", itemID, "error", err)
		return
	}
	slog.Info("?대?吏 ?꾨쿋???앹꽦 ?꾨즺", "item_id", itemID)
}

// ImageEmbedding? ??λ맂 ?대?吏瑜?遺덈윭? Gemini 紐⑤뜽???꾨쿋?⑹쓣 ?붿껌?섍퀬 寃곌낵瑜?DB????ν븳??
func ImageEmbedding(ctx context.Context, itemID int64, contextText string) error {
	apiKey := config.EnvMap["GEMINI_AI_KEY"]
	if apiKey == "" {
		return fmt.Errorf("gemini api 키가 비어 있습니다")
	}

	queries, err := db.GetQueries()
	if err != nil {
		return fmt.Errorf("쿼리 객체 생성 실패: %w", err)
	}

	item, err := queries.GetItemContent(ctx, itemID)
	if err != nil {
		return fmt.Errorf("아이템 이미지를 조회하지 못했어요: %w", err)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("gemini 클라이언트 생성 실패: %w", err)
	}

	payload := strings.TrimSpace(contextText)
	if payload == "" {
		payload = fmt.Sprintf("closet item #%d mime=%s size=%d bytes", itemID, item.MimeType, len(item.Bytes))
	}

	resp, err := client.Models.EmbedContent(ctx, embeddingModelName, genai.Text(payload), &genai.EmbedContentConfig{
		TaskType: "RETRIEVAL_DOCUMENT",
	})
	if err != nil {
		return fmt.Errorf("임베딩 요청 실패: %w", err)
	}
	if len(resp.Embeddings) == 0 || len(resp.Embeddings[0].Values) == 0 {
		return fmt.Errorf("임베딩 응답이 비어 있어요")
	}

	vector := resp.Embeddings[0].Values
	vecBytes := float32SliceToBytes(vector)

	if err := queries.PutEmbedding(ctx, db.PutEmbeddingParams{
		ItemID: itemID,
		Model:  embeddingModelName,
		Dim:    int64(len(vector)),
		VecF32: vecBytes,
	}); err != nil {
		return fmt.Errorf("임베딩 저장에 실패했어요: %w", err)
	}

	return nil
}

func float32SliceToBytes(values []float32) []byte {
	buf := make([]byte, len(values)*4)
	for i, v := range values {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(v))
	}
	return buf
}
