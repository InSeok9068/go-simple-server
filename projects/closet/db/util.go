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

// initDB는 데이터베이스 연결을 초기화합니다.
func initDB() {
	var err error
	dbConn, err = connection.AppDBOpen()
	if err != nil {
		errInit = err
		return
	}
	dbNotHookConn, err = connection.AppDBOpen(false)
	if err != nil {
		errInit = err
		return
	}
	queries = New(dbConn)
	queriesNotHook = New(dbNotHookConn)
}

// GetDB 는 공용 DB 연결을 반환합니다.
// 기본 값 : hook된 연결을 반환합니다.
// hooked 파라미터가 true이면 hook된 연결을 반환하고, false이면 hook되지 않은 연결을 반환합니다.
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

// GetQueries 는 공용 쿼리 객체를 반환합니다.
// 기본 값 : hook된 쿼리 객체를 반환합니다.
// hooked 파라미터가 true이면 hook된 연결을 반환하고, false이면 hook되지 않은 연결을 반환합니다.
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
	var err1, err2 error
	if dbConn != nil {
		err1 = dbConn.Close()
	}
	if dbNotHookConn != nil {
		err2 = dbNotHookConn.Close()
	}
	if err1 != nil {
		return err1
	}
	return err2
}
