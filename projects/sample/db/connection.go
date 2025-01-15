package db

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func DbQueries() (*Queries, context.Context) {
	ctx := context.Background()
	dbCon, err := AppDbOpen()
	if err != nil {
		slog.Error(err.Error())
	}
	queries := New(dbCon)
	return queries, ctx
}

func AppDbOpen() (*sql.DB, error) {
	return sql.Open("sqlite3", os.Getenv("APP_DATABASE_URL"))
}

func LogDbOpen() (*sql.DB, error) {
	return sql.Open("sqlite3", os.Getenv("LOG_DATABASE_URL"))
}
