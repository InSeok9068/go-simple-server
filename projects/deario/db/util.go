package db

import (
	"context"
	"fmt"
	"simple-server/internal/connection"
)

// GetQueries 는 DB 연결을 열고 쿼리 객체를 반환합니다
func GetQueries(ctx context.Context) (*Queries, error) {
	dbCon, err := connection.AppDBOpen()
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 실패: %w", err)
	}

	return New(dbCon), nil
}
