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
	if _, err := c.AddFunc("0 21 * * *", func() {
		ctx := context.Background()
		client, err := middleware.App.Messaging(ctx)
		if err != nil {
			slog.Error("메시징 클라이언트 생성 실패", "error", err)
			return
		}

		dbCon, err := connection.AppDBOpen()
		if err != nil {
			slog.Error("데이터베이스 연결 실패", "error", err)
			return
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
			slog.Error("푸시 발송 실패", "error", err)
		}

		slog.Debug("푸시 발송 응답", "response", response)
	}); err != nil {
		slog.Error("스케줄 등록 실패", "error", err)
	}
}
