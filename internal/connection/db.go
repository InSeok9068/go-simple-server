package connection

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func AppDBOpen() (*sql.DB, error) {
	return sql.Open("sqlite3", os.Getenv("APP_DATABASE_URL"))
}

func LogDBOpen() (*sql.DB, error) {
	return sql.Open("sqlite3", os.Getenv("LOG_DATABASE_URL"))
}
