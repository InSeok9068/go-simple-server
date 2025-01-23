package db

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func DbQueries() (*Queries, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbCon, err := AppDbOpen()
	if err != nil {
		slog.Error("Failed to open database", "error", err.Error())
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
