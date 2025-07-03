package config

import (
	"net/http"
	"simple-server/internal/connection"
	"time"

	"github.com/doganarif/govisual"
	"github.com/labstack/echo/v4"
)

func TransferEchoToGoVisualServerOnlyDev(e *echo.Echo, port string) *http.Server {
	db, _ := connection.LogDBOpen()

	server := govisual.Wrap(
		e,
		govisual.WithRequestBodyLogging(true),
		govisual.WithResponseBodyLogging(true),
		govisual.WithSQLiteStorageDB(db, "govisual"),
	)

	return &http.Server{
		Addr:         ":" + port,
		Handler:      server,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
