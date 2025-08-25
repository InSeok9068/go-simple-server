package connection

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/qustavo/sqlhooks/v2"

	"modernc.org/sqlite"
)

var (
	once       sync.Once
	driverName = "sqlite-hooked"
)

type Hooks struct{}

func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, "begin", time.Now()), nil //nolint:staticcheck
}

func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value("begin").(time.Time)
	slog.DebugContext(ctx, "SQL 실행", "query", query, "args", args, "duration", time.Since(begin))
	return ctx, nil
}

func AppDBOpen(hooked ...bool) (*sql.DB, error) {
	isHooked := os.Getenv("ENV") == "dev"

	if len(hooked) > 0 {
		isHooked = hooked[0]
	}

	var db *sql.DB
	var err error
	if isHooked {
		once.Do(func() {
			sql.Register(driverName, sqlhooks.Wrap(&sqlite.Driver{}, &Hooks{}))
		})
		db, err = sql.Open(driverName, os.Getenv("APP_DATABASE_URL"))
	} else {
		db, err = sql.Open("sqlite", os.Getenv("APP_DATABASE_URL"))
	}
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 실패: %w", err)
	}

	// 메인 DB 설정
	db.SetMaxOpenConns(5) // 최대 연결 수 (읽기 쓰기 동시)
	db.SetMaxIdleConns(5) // 최대 유휴 연결 수 (읽기 쓰기 동시)

	return db, nil
}

func LogDBOpen() (*sql.DB, error) {
	var db *sql.DB
	var err error
	db, err = sql.Open("sqlite", os.Getenv("LOG_DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("로그 데이터베이스 연결 실패: %w", err)
	}

	// 로그 DB 설정
	db.SetMaxOpenConns(1) // 최대 연결 수 (로그는 동시성 필요 없음)
	db.SetMaxIdleConns(1) // 최대 유휴 연결 수 (로그는 동시성 필요 없음)
	return db, nil
}
