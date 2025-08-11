package db

import (
	"context"
	"testing"
)

func TestGetDBAndQueriesReuseConnection(t *testing.T) {
	t.Setenv("APP_DATABASE_URL", "file::memory:?cache=shared")

	ctx := context.Background()

	q1, err := GetQueries(ctx)
	if err != nil {
		t.Fatalf("첫 번째 쿼리 객체 생성 실패: %v", err)
	}

	db1, err := GetDB()
	if err != nil {
		t.Fatalf("첫 번째 DB 연결 실패: %v", err)
	}

	q2, err := GetQueries(ctx)
	if err != nil {
		t.Fatalf("두 번째 쿼리 객체 생성 실패: %v", err)
	}

	db2, err := GetDB()
	if err != nil {
		t.Fatalf("두 번째 DB 연결 실패: %v", err)
	}

	if q1 != q2 {
		t.Errorf("GetQueries가 다른 인스턴스를 반환했습니다")
	}

	if db1 != db2 {
		t.Errorf("GetDB가 다른 인스턴스를 반환했습니다")
	}

	if err := Close(); err != nil {
		t.Fatalf("DB 종료 실패: %v", err)
	}
}
