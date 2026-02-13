package connection

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/qustavo/sqlhooks/v2"
	"go.opentelemetry.io/otel/attribute"

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
	// isHooked := os.Getenv("ENV") == "dev"
	isHooked := true

	if len(hooked) > 0 {
		isHooked = hooked[0]
	}

	var (
		db  *sql.DB
		err error
	)

	if isHooked {
		once.Do(func() {
			sql.Register(driverName, sqlhooks.Wrap(&sqlite.Driver{}, &Hooks{}))
		})
		db, err = otelsql.Open(
			driverName,
			os.Getenv("APP_DATABASE_URL"),
			// DB 기본정보 포함
			otelsql.WithAttributes(
				attribute.String("db.system", "sqlite"),
				attribute.String("db.name", os.Getenv("SERVICE_NAME")),
			),
			// 실행 쿼리 파라미터 포함
			otelsql.WithAttributesGetter(func(ctx context.Context, method otelsql.Method, query string, args []driver.NamedValue) []attribute.KeyValue {
				if len(args) == 0 {
					return nil
				}

				params := make([]string, 0, len(args))
				for i, arg := range args {
					paramName := arg.Name
					if paramName == "" {
						paramName = fmt.Sprintf("$%d", i+1)
					}
					params = append(params, fmt.Sprintf("%s=%v", paramName, arg.Value))
				}

				return []attribute.KeyValue{
					attribute.String("db.parameters", strings.Join(params, ", ")),
				}
			}),
			// span 이름 Query ID 사용
			otelsql.WithSpanNameFormatter(sqlSpanNameFormatter),
		)
	} else {
		db, err = sql.Open("sqlite", os.Getenv("APP_DATABASE_URL"))
	}
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 실패: %w", err)
	}

	// 메인 DB 설정
	db.SetMaxOpenConns(5) // 최대 연결 수 (읽기 쓰기 동시)
	db.SetMaxIdleConns(5) // 최대 유휴 연결 수 (읽기 쓰기 동시)
	// TODO: 추후에 메모리 압박이 발생할 경우 추가
	// db.SetConnMaxIdleTime(1 * time.Hour) // 1h 유휴 연결 유지 시간

	// 연결 테스트

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func LogDBOpen() (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)
	db, err = sql.Open("sqlite", os.Getenv("LOG_DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("로그 데이터베이스 연결 실패: %w", err)
	}

	// 로그 DB 설정
	db.SetMaxOpenConns(1) // 최대 연결 수 (로그는 동시성 필요 없음)
	db.SetMaxIdleConns(1) // 최대 유휴 연결 수 (로그는 동시성 필요 없음)

	// 연결 테스트
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func sqlSpanNameFormatter(_ context.Context, method otelsql.Method, query string) string {
	if name := extractQueryName(query); name != "" {
		return "SQL " + name
	}

	return "SQL " + string(method)
}

func extractQueryName(query string) string {
	if query == "" {
		return ""
	}

	const prefix = "-- name:"

	for _, line := range strings.Split(query, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !strings.HasPrefix(strings.ToLower(trimmed), prefix) {
			continue
		}

		rest := strings.TrimSpace(trimmed[len(prefix):])
		if rest == "" {
			return ""
		}

		parts := strings.Fields(rest)
		if len(parts) == 0 {
			return ""
		}

		return parts[0]
	}

	return ""
}
