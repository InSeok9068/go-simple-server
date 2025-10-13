package config

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/doganarif/govisual"
	"github.com/labstack/echo/v4"
)

func TransferEchoToGoVisualServerOnlyDev(e *echo.Echo, port string) *http.Server {
	url := "file:./shared/log/govisual.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&mode=rwc"
	db, err := sql.Open("sqlite", url)
	if err != nil {
		slog.Error("[Govisual] 로그 데이터베이스 연결 실패", "error", err)
		return nil
	}

	server := govisual.Wrap(
		e,
		govisual.WithRequestBodyLogging(true),
		govisual.WithResponseBodyLogging(true),
		govisual.WithSQLiteStorageDB(db, "govisual"),
	)

	return &http.Server{
		Addr:         ":" + port,
		Handler:      server,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
}
