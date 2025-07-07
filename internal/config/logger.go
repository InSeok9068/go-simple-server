package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"simple-server/internal/connection"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
	"go.opentelemetry.io/otel/trace"
	_ "modernc.org/sqlite"
)

var initOnce sync.Once

// MultiHandler : 여러 slog.Handler를 지원하는 커스텀 핸들러
type MultiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler : MultiHandler 생성자
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

// Handle : 모든 핸들러에 로그를 전달
func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if spanCtx.IsValid() {
		r.AddAttrs(slog.String("trace_id", spanCtx.TraceID().String()))
	}
	for _, handler := range m.handlers {
		rec := r.Clone()
		if err := handler.Handle(ctx, rec); err != nil {
			return err
		}
	}
	return nil
}

// Enabled : 모든 핸들러가 동일하게 활성화될지 판단
func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range m.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// WithAttrs : 속성을 추가한 새로운 핸들러 생성
func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var newHandlers []slog.Handler
	for _, handler := range m.handlers {
		newHandlers = append(newHandlers, handler.WithAttrs(attrs))
	}
	return NewMultiHandler(newHandlers...)
}

// WithGroup : 그룹 이름을 추가한 새로운 핸들러 생성
func (m *MultiHandler) WithGroup(name string) slog.Handler {
	var newHandlers []slog.Handler
	for _, handler := range m.handlers {
		newHandlers = append(newHandlers, handler.WithGroup(name))
	}
	return NewMultiHandler(newHandlers...)
}

type logEntry struct {
	level   slog.Level
	message string
	data    string
}

type DatabaseHandler struct {
	db     *sql.DB
	level  slog.Leveler
	mu     sync.Mutex
	buffer []logEntry
}

func NewDatabaseHandler(db *sql.DB, level slog.Leveler) *DatabaseHandler {
	h := &DatabaseHandler{
		db:    db,
		level: level,
	}
	go h.start()
	return h
}

func (h *DatabaseHandler) start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		h.flush()
	}
}

func (h *DatabaseHandler) flush() {
	h.mu.Lock()
	entries := h.buffer
	h.buffer = nil
	h.mu.Unlock()
	if len(entries) == 0 {
		return
	}
	ctx := context.Background()
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("로그 배치 시작 실패", "error", err)
		return
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO _logs (level, message, data) VALUES (?, ?, ?)")
	if err != nil {
		slog.Error("로그 배치 쿼리 준비 실패", "error", err)
		_ = tx.Rollback()
		return
	}
	defer stmt.Close()
	for _, e := range entries {
		if _, err := stmt.ExecContext(ctx, e.level, e.message, e.data); err != nil {
			slog.Error("로그 배치 실패", "error", err)
			_ = tx.Rollback()
			return
		}
	}
	if err := tx.Commit(); err != nil {
		slog.Error("로그 배치 커밋 실패", "error", err)
	}
}

func (h *DatabaseHandler) Handle(ctx context.Context, r slog.Record) error {
	logMessage := r.Message
	logLevel := r.Level.Level()
	data := make(map[string]any)
	r.Attrs(func(attr slog.Attr) bool {
		data[attr.Key] = attr.Value.Any()
		return true
	})
	data["service"] = os.Getenv("SERVICE_NAME")
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	h.mu.Lock()
	h.buffer = append(h.buffer, logEntry{level: logLevel, message: logMessage, data: string(jsonBytes)})
	h.mu.Unlock()
	return nil
}

func (h *DatabaseHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *DatabaseHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *DatabaseHandler) WithGroup(name string) slog.Handler {
	return h
}

func InitLoggerWithDatabase() {
	initOnce.Do(func() {
		os.Setenv("LOG_DATABASE_URL", LogDatabaseURL())
		dbCon, err := connection.LogDBOpen()
		if err != nil {
			slog.Error("로그 데이터베이스 연결 실패", "error", err)
			return
		}

		var level slog.Leveler
		if IsDevEnv() {
			level = slog.LevelDebug
		} else {
			level = slog.LevelInfo
		}

		// Database Handler
		databaseHandler := NewDatabaseHandler(dbCon, level)

		// // File Handler
		// file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// if err != nil {
		// 	slog.Error("Failed to open log file", "error", err)
		// 	return
		// }
		// fileHandler := slog.NewTextHandler(file, &slog.HandlerOptions{})

		var consoleHandler slog.Handler
		if IsDevEnv() {
			consoleHandler = tint.NewHandler(os.Stderr, &tint.Options{
				Level: level,
			})
		} else {
			consoleHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				Level: level,
			})
		}

		// MultiHandler: Combine all handlers
		multiHandler := NewMultiHandler(consoleHandler, databaseHandler)
		slog.SetDefault(slog.New(multiHandler))
		log.SetOutput(os.Stderr)
	})
}

func CustomLogValuesFunc(c echo.Context, v middleware.RequestLoggerValues) error {
	// 요청 정보 로그 기록
	method := c.Request().Method
	requestURI := c.Request().RequestURI
	remoteIP := v.RemoteIP
	userIP := c.RealIP()

	slog.InfoContext(c.Request().Context(), fmt.Sprintf(`%s %s`, method, requestURI),
		"exec_time", v.Latency.Microseconds(),
		"id", v.RequestID,
		"type", "request",
		"status", c.Response().Status,
		"method", method,
		"url", requestURI,
		"referer", c.Request().Referer(),
		"remote_ip", remoteIP,
		"user_ip", userIP,
		"user_agent", c.Request().UserAgent(),
		"error", v.Error,
		"details", "")
	return nil
}
