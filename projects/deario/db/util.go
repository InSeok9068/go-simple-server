package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"simple-server/internal/connection"
)

var (
	once    sync.Once
	dbConn  *sql.DB
	queries *Queries
	errInit error
)

func initDB() {
	dbConn, errInit = connection.AppDBOpen()
	if errInit == nil {
		queries = New(dbConn)
	}
}

// GetDB 는 공용 DB 연결을 반환합니다
func GetDB() (*sql.DB, error) {
	once.Do(initDB)
	if errInit != nil {
		return nil, fmt.Errorf("데이터베이스 연결 실패: %w", errInit)
	}
	return dbConn, nil
}

// GetQueries 는 공용 쿼리 객체를 반환합니다
func GetQueries(ctx context.Context) (*Queries, error) {
	if _, err := GetDB(); err != nil {
		return nil, err
	}
	return queries, nil
}

// Close 는 공용 DB 연결을 종료합니다
func Close() error {
	if dbConn != nil {
		return dbConn.Close()
	}
	return nil
}
