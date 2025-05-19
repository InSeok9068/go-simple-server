package tasks

import (
	"context"
	"log/slog"
	"simple-server/internal/connection"
	"simple-server/internal/middleware"
	"simple-server/projects/deario/db"

	"firebase.google.com/go/v4/messaging"
	"github.com/robfig/cron/v3"
)

func PushTask(c *cron.Cron) {
	_, _ = c.AddFunc("0 21 * * *", func() {
		ctx := context.Background()
		client, _ := middleware.App.Messaging(ctx)

		dbCon, err := connection.AppDBOpen()
		if err != nil {
			slog.Error("Failed to open database", "error", err.Error())
		}
		queries := db.New(dbCon)

		uid := "6KWofk1AVdZolC94UAuRuAB1wj13"
		pushKey, err := queries.GetPushKey(ctx, uid)

		if err != nil {
			return
		}

		message := &messaging.Message{
			Data: map[string]string{
				"title": "매일 알림",
				"body":  "오늘 하루는 어땠나요?",
			},
			Token: pushKey.Token,
		}

		response, err := client.Send(ctx, message)
		if err != nil {
			slog.Error("Failed to send push", "error", err.Error())
		}

		slog.Info("푸시 발송 응답", "response", response)
	})
}
