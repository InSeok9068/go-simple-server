package connection

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func AppDBOpen() (*sql.DB, error) {
	return sql.Open("sqlite", os.Getenv("APP_DATABASE_URL"))
}

func LogDBOpen() (*sql.DB, error) {
	return sql.Open("sqlite", os.Getenv("LOG_DATABASE_URL"))
}
