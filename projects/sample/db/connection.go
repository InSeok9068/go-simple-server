package db

import (
	"context"
	"database/sql"
	"log/slog"
)

func DbConnection() (*Queries, context.Context) {
	ctx := context.Background()
	dbCon, err := sql.Open("sqlite3", "file:./projects/sample/pb_data/data.db")
	if err != nil {
		slog.Error(err.Error())
	}
	queries := New(dbCon)
	return queries, ctx
}
