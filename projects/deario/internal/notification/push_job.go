package notification

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"simple-server/internal/middleware"
	"simple-server/projects/deario/db"

	"firebase.google.com/go/v4/messaging"
	"github.com/robfig/cron/v3"
	"maragu.dev/goqite"
	"maragu.dev/goqite/jobs"
)

type Payload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Token string `json:"token"`
}

var pushQ *goqite.Queue

func PushSendCron(c *cron.Cron) {
	queries, err := db.GetQueries(false)
	if err != nil {
		slog.Error("쿼리 로드 실패", "error", err)
		return
	}

	if _, err := c.AddFunc("@every 1m", func() {
		ctx := context.Background()
		now := time.Now().Format("15:04")

		targets, err := queries.ListPushTargets(ctx)
		if err != nil {
			slog.Error("푸시 대상 조회 실패", "error", err)
			return
		}

		for _, target := range targets {
			if target.PushTime != now {
				continue
			}

			payload := Payload{
				Title: "매일 알림",
				Body:  "오늘 하루는 어땠나요?",
				Token: target.PushToken,
			}
			b, _ := json.Marshal(payload)

			if _, err := jobs.Create(ctx, pushQ, "send", goqite.Message{Body: b}); err != nil {
				slog.Error("푸시 발송 실패", "error", err)
			}
		}
	}); err != nil {
		slog.Error("스케줄 등록 실패", "error", err)
	}
}

func PushSendJob() {
	pushdb, err := db.GetDB(false)
	if err != nil {
		slog.Error("데이터베이스 연결 실패", "error", err)
		return
	}
	defer db.Close()
	pushQ = goqite.New(goqite.NewOpts{
		DB:   pushdb,
		Name: "push",
	})
	r := jobs.NewRunner(jobs.NewRunnerOpts{
		Limit:        1,
		Log:          slog.Default(),
		PollInterval: 1 * time.Second,
		Queue:        pushQ,
	})

	r.Register("send", func(ctx context.Context, m []byte) error {
		client, err := middleware.App.Messaging(ctx)
		if err != nil {
			slog.Error("메시징 클라이언트 생성 실패", "error", err)
			return err
		}

		var payload Payload
		if err := json.Unmarshal(m, &payload); err != nil {
			slog.Error("푸시 데이터 해독 실패", "error", err)
			return err
		}

		message := &messaging.Message{
			Data: map[string]string{
				"title": payload.Title,
				"body":  payload.Body,
			},
			Token: payload.Token,
		}

		if _, err := client.Send(ctx, message); err != nil {
			slog.Error("푸시 발송 실패", "error", err)
			return err
		}
		return nil
	})

	r.Start(context.Background())
}
