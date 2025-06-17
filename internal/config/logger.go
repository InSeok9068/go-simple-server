package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net"
	"os"
	"simple-server/internal/connection"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	for _, handler := range m.handlers {
		if err := handler.Handle(ctx, r); err != nil {
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

type DatabaseHandler struct {
	db *sql.DB
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
	jsonBytes, _ := json.Marshal(data)

	_, _ = h.db.ExecContext(ctx, "INSERT INTO _logs (level, message, data) VALUES (?, ?, ?)",
		logLevel,
		logMessage,
		string(jsonBytes))

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

func LoggerWithDatabaseInit() {
	initOnce.Do(func() {
		os.Setenv("LOG_DATABASE_URL", LogDatabaseURL())
		dbCon, err := connection.LogDBOpen()
		if err != nil {
			slog.Error("Failed to open database", "error", err)
			return
		}

		// Database Handler
		databaseHandler := &DatabaseHandler{db: dbCon}

		// // File Handler
		// file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// if err != nil {
		// 	slog.Error("Failed to open log file", "error", err)
		// 	return
		// }
		// fileHandler := slog.NewTextHandler(file, &slog.HandlerOptions{})

		// Console Handler
		consoleHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})

		// MultiHandler: Combine all handlers
		multiHandler := NewMultiHandler(consoleHandler, databaseHandler)
		slog.SetDefault(slog.New(multiHandler))
		log.SetOutput(os.Stdout)
	})
}

func CustomLogValuesFunc(c echo.Context, v middleware.RequestLoggerValues) error {
	// 요청 정보 로그 기록
	method := c.Request().Method
	requestURI := c.Request().RequestURI
	remoteIP, _, _ := net.SplitHostPort(v.RemoteIP)
	userIP, _, _ := net.SplitHostPort(c.RealIP())

	slog.Info("request",
		"execTime", v.Latency.Microseconds(),
		"id", v.RequestID,
		"type", "request",
		"status", c.Response().Status,
		"method", method,
		"url", requestURI,
		"referer", c.Request().Referer(),
		"remoteIP", remoteIP,
		"userIP", userIP,
		"userAgent", c.Request().UserAgent(),
		"error", v.Error,
		"details", "")
	return nil
}
