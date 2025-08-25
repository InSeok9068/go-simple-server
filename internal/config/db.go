package config

import (
	"context"
	"database/sql"
	"expvar"
	"fmt"
	"strings"
	"time"
)

func AppDatabaseURL(serviceName string) string {
	appPragmas := []string{
		"_pragma=journal_mode(WAL)",            // 동시성 향상
		"_pragma=synchronous(NORMAL)",          // 정합성 향상
		"_pragma=busy_timeout(5000)",           // 쓰기 경합 대기 5초
		"_pragma=foreign_keys(ON)",             // 외래 키 허용
		"_pragma=temp_store(MEMORY)",           // 임시 저장소 메모리
		"_pragma=journal_size_limit(67108864)", // 64MB
	}
	pragmas := strings.Join(appPragmas, "&")
	var url string
	if IsDevEnv() {
		url = fmt.Sprintf(`file:./projects/%s/data/data.db?%s&mode=rwc`, serviceName, pragmas)
	} else {
		url = fmt.Sprintf(`file:/srv/%s/data/data.db?%s&mode=rwc`, serviceName, pragmas)
	}
	return url
}

func LogDatabaseURL() string {
	logPragmas := []string{
		"_pragma=journal_mode(WAL)",            // 동시성 향상
		"_pragma=synchronous(NORMAL)",          // 정합성 향상
		"_pragma=busy_timeout(5000)",           // 쓰기 경합 대기 5초
		"_pragma=temp_store(MEMORY)",           // 임시 저장소 메모리
		"_pragma=journal_size_limit(67108864)", // 64MB
	}
	pragmas := strings.Join(logPragmas, "&")
	var url string
	if IsDevEnv() {
		url = fmt.Sprintf(`file:./shared/log/auxiliary.db?%s&mode=rwc`, pragmas)
	} else {
		url = fmt.Sprintf(`file:/srv/log/auxiliary.db?%s&mode=rwc`, pragmas)
	}
	return url
}

func PublishDBVars(appDB *sql.DB) {
	expvar.Publish("db", expvar.Func(func() any {
		s := appDB.Stats()

		// 가볍게: 짧은 타임아웃으로 PRAGMA 조회
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		pr := map[string]any{
			"journal_mode":       pragmaString(ctx, appDB, "journal_mode"), // e.g. "wal"
			"synchronous":        pragmaEnum(ctx, appDB, "synchronous", map[int64]string{0: "OFF", 1: "NORMAL", 2: "FULL", 3: "EXTRA"}),
			"busy_timeout_ms":    pragmaInt(ctx, appDB, "busy_timeout"),
			"foreign_keys":       pragmaBool(ctx, appDB, "foreign_keys"),
			"temp_store":         pragmaEnum(ctx, appDB, "temp_store", map[int64]string{0: "DEFAULT", 1: "FILE", 2: "MEMORY"}),
			"journal_size_limit": pragmaInt(ctx, appDB, "journal_size_limit"), // -1이면 무제한
			"wal_autocheckpoint": pragmaInt(ctx, appDB, "wal_autocheckpoint"), // 페이지 단위
			// 필요시: page_size, cache_size 등도 추가
		}

		return map[string]any{
			"max_open":      s.MaxOpenConnections,
			"open":          s.OpenConnections,
			"in_use":        s.InUse,
			"idle":          s.Idle,
			"wait_count":    s.WaitCount,
			"wait_duration": s.WaitDuration.String(),
			"wait_ns":       s.WaitDuration.Nanoseconds(),
			"pragma":        pr,
		}
	}))
}

func pragmaString(ctx context.Context, db *sql.DB, name string) string {
	var v string
	_ = db.QueryRowContext(ctx, "PRAGMA "+name).Scan(&v)
	return v
}
func pragmaInt(ctx context.Context, db *sql.DB, name string) int64 {
	var v int64
	_ = db.QueryRowContext(ctx, "PRAGMA "+name).Scan(&v)
	return v
}
func pragmaBool(ctx context.Context, db *sql.DB, name string) bool {
	var v int64
	_ = db.QueryRowContext(ctx, "PRAGMA "+name).Scan(&v)
	return v == 1
}
func pragmaEnum(ctx context.Context, db *sql.DB, name string, m map[int64]string) any {
	var v int64
	if err := db.QueryRowContext(ctx, "PRAGMA "+name).Scan(&v); err == nil {
		if s, ok := m[v]; ok {
			return s
		}
		return v
	}
	return nil
}
