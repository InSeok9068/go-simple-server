package cmd

import (
	"github.com/pocketbase/pocketbase"
	"log/slog"
)

func AdminServer() {
	// 어드민 서버
	app := pocketbase.New()

	if err := app.Start(); err != nil {
		slog.Error(err.Error())
	}
}
