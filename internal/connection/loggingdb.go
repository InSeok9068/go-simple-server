package connection

import (
	"context"
	"database/sql"
	"log/slog"
)

type LoggingDB struct {
	*sql.DB
}

func (ldb *LoggingDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	slog.Debug("쿼리 Exec", "query", query, "args", args)
	res, err := ldb.DB.ExecContext(ctx, query, args...)
	if err != nil {
		slog.Debug("쿼리 Exec 실패", "error", err)
		return res, err
	}
	if rows, err := res.RowsAffected(); err == nil {
		slog.Debug("쿼리 Exec 성공", "rows", rows)
	}
	return res, nil
}

func (ldb *LoggingDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	slog.Debug("쿼리 Prepare", "query", query)
	stmt, err := ldb.DB.PrepareContext(ctx, query)
	if err != nil {
		slog.Debug("쿼리 Prepare 실패", "error", err)
	}
	return stmt, err
}

func (ldb *LoggingDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	slog.Debug("쿼리 Query", "query", query, "args", args)
	rows, err := ldb.DB.QueryContext(ctx, query, args...)
	if err != nil {
		slog.Debug("쿼리 Query 실패", "error", err)
	}
	return rows, err
}

func (ldb *LoggingDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	slog.Debug("쿼리 QueryRow", "query", query, "args", args)
	return ldb.DB.QueryRowContext(ctx, query, args...)
}
