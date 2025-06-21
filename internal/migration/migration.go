package migration

import (
	"context"
	"database/sql"
	"io/fs"
	"time"

	"github.com/pressly/goose/v3"
	"simple-server/internal/config"
)

// Up 함수는 주어진 데이터베이스와 마이그레이션 파일 시스템을 이용해
// goose 마이그레이션을 실행한다.
func Up(db *sql.DB, migrations fs.FS) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	provider, err := goose.NewProvider(
		goose.DialectSQLite3, db, migrations,
		goose.WithVerbose(config.IsDevEnv()),
	)
	if err != nil {
		return err
	}

	_, err = provider.Up(ctx)
	return err
}
