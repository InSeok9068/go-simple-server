package db

import (
	"database/sql"
	"fmt"
	"sync"

	"simple-server/internal/connection"
)

var (
	once           sync.Once
	dbConn         *sql.DB
	dbNotHookConn  *sql.DB
	queries        *Queries
	queriesNotHook *Queries
	errInit        error
)

func initDB() {
	dbConn, errInit = connection.AppDBOpen()
	dbNotHookConn, errInit = connection.AppDBOpen(false)
	if errInit == nil {
		queries = New(dbConn)
		queriesNotHook = New(dbNotHookConn)
	}
}

// GetDB 는 공용 DB 연결을 반환합니다
func GetDB(hooked ...bool) (*sql.DB, error) {
	once.Do(initDB)
	if errInit != nil {
		return nil, fmt.Errorf("데이터베이스 연결 실패: %w", errInit)
	}
	if len(hooked) > 0 && !hooked[0] {
		return dbNotHookConn, nil
	}
	return dbConn, nil
}

// GetQueries 는 공용 쿼리 객체를 반환합니다
func GetQueries(hooked ...bool) (*Queries, error) {
	if _, err := GetDB(hooked...); err != nil {
		return nil, err
	}
	if len(hooked) > 0 && !hooked[0] {
		return queriesNotHook, nil
	}
	return queries, nil
}

// Close 는 공용 DB 연결을 종료합니다
func Close() error {
	if dbConn != nil {
		return dbConn.Close()
	}
	if dbNotHookConn != nil {
		return dbNotHookConn.Close()
	}
	return nil
}
