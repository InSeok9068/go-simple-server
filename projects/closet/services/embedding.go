package services

import (
	"context"
	"log/slog"
)

// EnqueueEmbeddingJob은 임베딩 생성을 비동기 작업 큐에 추가한다.
// 현재는 스텁으로 동작하며, 실제 임베딩 연동 시 이 위치에서 외부 서비스 호출을 추가하면 된다.
func EnqueueEmbeddingJob(itemID int64) {
	slog.Info("임베딩 작업을 준비합니다", "item_id", itemID)
	// TODO: 임베딩 API 연동 시 구현
}

// ImageEmbedding은 업로드 이미지의 임베딩을 생성한다.
// 실제 구현 전까지는 nil을 반환한다.
func ImageEmbedding(ctx context.Context, itemID int64) error {
	slog.Debug("임베딩 스텁 호출", "item_id", itemID)
	return nil
}
