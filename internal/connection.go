package internal

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// func DBQueries() (*Queries, context.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	dbCon, err := AppDBOpen()
// 	if err != nil {
// 		slog.Error("Failed to open database", "error", err.Error())
// 	}
// 	queries := New(dbCon)
// 	return queries, ctx
// }

func AppDBOpen() (*sql.DB, error) {
	return sql.Open("sqlite3", os.Getenv("APP_DATABASE_URL"))
}

func LogDBOpen() (*sql.DB, error) {
	return sql.Open("sqlite3", os.Getenv("LOG_DATABASE_URL"))
}
