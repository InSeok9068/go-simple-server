package internal

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"sync"

	"simple-server/projects/homepage/db"

	_ "github.com/mattn/go-sqlite3"
)

var initOnce sync.Once

type DatabaseHandler struct {
	db *sql.DB
}

func (h *DatabaseHandler) Handle(ctx context.Context, r slog.Record) error {
	logMessage := r.Message
	logLevel := r.Level.Level()

	_, err := h.db.ExecContext(ctx, "INSERT INTO _logs (level, message, data) VALUES (?, ?, ?)", logLevel, logMessage, "{}")
	if err != nil {
		log.Println("Failed to insert log into database", "error", err)
	}

	log.Printf("%s - %s: %s", r.Time.Format("2006-01-02 15:04:05"), logLevel, logMessage)

	return nil
}

func (h *DatabaseHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelInfo
}

func (h *DatabaseHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *DatabaseHandler) WithGroup(name string) slog.Handler {
	return h
}

func DatabaseLogInit() {
	initOnce.Do(func() {
		dbCon, err := db.LogDbOpen()
		if err != nil {
			slog.Error("Failed to open database", "error", err)
			return
		}

		databaseHandler := &DatabaseHandler{db: dbCon}
		slog.SetDefault(slog.New(databaseHandler))
		log.SetOutput(os.Stdout)
	})
}
