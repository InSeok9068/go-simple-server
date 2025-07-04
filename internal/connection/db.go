package connection

import (
	"context"
	"database/sql"
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
	return ctx, nil
}

func AppDBOpen() (*sql.DB, error) {
	isHooked := os.Getenv("ENV") == "dev"
	if isHooked {
		once.Do(func() {
			sql.Register(driverName, sqlhooks.Wrap(&sqlite.Driver{}, &Hooks{}))
		})
		return sql.Open(driverName, os.Getenv("APP_DATABASE_URL"))
	} else {
		return sql.Open("sqlite", os.Getenv("APP_DATABASE_URL"))
	}
}

func LogDBOpen() (*sql.DB, error) {
	return sql.Open("sqlite", os.Getenv("LOG_DATABASE_URL"))
}
