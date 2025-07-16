package tasks

import (
	"context"
	"log/slog"
	"time"

	"simple-server/internal/middleware"
	"simple-server/projects/deario/db"

	"firebase.google.com/go/v4/messaging"
	"github.com/robfig/cron/v3"
)

func PushTask(c *cron.Cron) {
	if _, err := c.AddFunc("@every 1m", func() {
		ctx := context.Background()
		now := time.Now().Format("15:04")

		client, err := middleware.App.Messaging(ctx)
		if err != nil {
			slog.Error("메시징 클라이언트 생성 실패", "error", err)
			return
		}

		queries, err := db.GetQueries(ctx)
		if err != nil {
			slog.Error("쿼리 로드 실패", "error", err)
			return
		}

		targets, err := queries.ListPushTargets(ctx)
		if err != nil {
			slog.Error("푸시 대상 조회 실패", "error", err)
			return
		}

		for _, target := range targets {
			if target.PushTime != now {
				continue
			}

			message := &messaging.Message{
				Data: map[string]string{
					"title": "매일 알림",
					"body":  "오늘 하루는 어땠나요?",
				},
				Token: target.PushToken,
			}

			if _, err := client.Send(ctx, message); err != nil {
				slog.Error("푸시 발송 실패", "error", err, "uid", target.Uid)
				continue
			}
		}
	}); err != nil {
		slog.Error("스케줄 등록 실패", "error", err)
	}
}
