package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"simple-server/internal/connection"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var initOnce sync.Once
var logBuffer []LogEntry
var bufferMutex sync.Mutex
var logDB *sql.DB

// LogEntry 로그 항목을 저장하는 구조체
type LogEntry struct {
	Level   int
	Message string
	Data    string
	Time    time.Time
}

// BatchDatabaseHandler 배치 처리를 지원하는 핸들러
type BatchDatabaseHandler struct {
	flushInterval time.Duration
}

func (h *BatchDatabaseHandler) Handle(ctx context.Context, r slog.Record) error {
	logMessage := r.Message
	logLevel := r.Level.Level()
	data := make(map[string]any)

	r.Attrs(func(attr slog.Attr) bool {
		data[attr.Key] = attr.Value.Any()
		return true
	})

	data["service"] = os.Getenv("SERVICE_NAME")
	jsonBytes, _ := json.Marshal(data)

	// 버퍼에 로그 항목 추가
	bufferMutex.Lock()
	logBuffer = append(logBuffer, LogEntry{
		Level:   int(logLevel),
		Message: logMessage,
		Data:    string(jsonBytes),
		Time:    time.Now(),
	})
	bufferMutex.Unlock()

	return nil
}

func (h *BatchDatabaseHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelInfo
}

func (h *BatchDatabaseHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *BatchDatabaseHandler) WithGroup(name string) slog.Handler {
	return h
}

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

// FlushLogs 버퍼에 있는 모든 로그를 데이터베이스에 저장
func FlushLogs() {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	if len(logBuffer) == 0 {
		return
	}

	// 데이터베이스 트랜잭션 시작
	tx, err := logDB.Begin()
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		return
	}

	// 준비된 statement 생성
	stmt, err := tx.Prepare("INSERT INTO _logs (level, message, data, created) VALUES (?, ?, ?, ?)")
	if err != nil {
		slog.Error("Failed to prepare statement", "error", err)
		_ = tx.Rollback()
		return
	}
	defer stmt.Close()

	// 모든 로그 항목을 데이터베이스에 저장
	for _, entry := range logBuffer {
		_, err := stmt.Exec(entry.Level, entry.Message, entry.Data, entry.Time)
		if err != nil {
			slog.Error("Failed to insert log", "error", err)
			_ = tx.Rollback()
			return
		}
	}

	// 트랜잭션 커밋
	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		_ = tx.Rollback()
		return
	}

	// 버퍼 비우기
	logBuffer = []LogEntry{}
}

// startLogFlusher 주기적으로 로그를 데이터베이스에 저장하는 고루틴 시작
func startLogFlusher(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			<-ticker.C
			FlushLogs()
		}
	}()
}

// LoggerWithDatabaseInit 로거 초기화 함수 수정
func LoggerWithDatabaseInit() {
	initOnce.Do(func() {
		os.Setenv("LOG_DATABASE_URL", "file:./shared/log_data/auxiliary.db")
		var err error
		logDB, err = connection.LogDBOpen()
		if err != nil {
			slog.Error("Failed to open database", "error", err)
			return
		}

		//// 로그 테이블이 존재하는지 확인, 없으면 생성
		//_, err = logDB.Exec(`
		//	CREATE TABLE IF NOT EXISTS _logs (
		//		id INTEGER PRIMARY KEY AUTOINCREMENT,
		//		level INTEGER NOT NULL,
		//		message TEXT NOT NULL,
		//		data TEXT NOT NULL,
		//		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		//	)
		//`)
		//if err != nil {
		//	slog.Error("Failed to create log table", "error", err)
		//	return
		//}

		// 배치 핸들러 생성
		batchHandler := &BatchDatabaseHandler{
			flushInterval: 5 * time.Second,
		}

		// 콘솔 핸들러
		consoleHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})

		// 다중 핸들러: 모든 핸들러 결합
		multiHandler := NewMultiHandler(consoleHandler, batchHandler)
		slog.SetDefault(slog.New(multiHandler))
		log.SetOutput(os.Stdout)

		// 로그 플러셔 시작 - 5초마다 실행
		startLogFlusher(5 * time.Second)
	})
}

func CustomLogValuesFunc(c echo.Context, v middleware.RequestLoggerValues) error {
	// 요청 정보 로그 기록
	method := c.Request().Method
	requestURI := c.Request().RequestURI
	remoteIP, _, _ := net.SplitHostPort(v.RemoteIP)
	userIP, _, _ := net.SplitHostPort(c.RealIP())

	slog.Info(fmt.Sprintf(`%s %s`, method, requestURI),
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
