package tasks

import (
	"context"
	"log/slog"
	"simple-server/internal/connection"
	"time"

	"maragu.dev/goqite"
	"maragu.dev/goqite/jobs"
)

var AiReportQ *goqite.Queue

func GenerateAIReportJob() {
	apiReportDb, err := connection.AppDBOpen(false)
	if err != nil {
		slog.Error("데이터베이스 연결 실패", "error", err)
		return
	}
	AiReportQ = goqite.New(goqite.NewOpts{
		DB:   apiReportDb,
		Name: "ai-report",
	})
	r := jobs.NewRunner(jobs.NewRunnerOpts{
		Limit:        3,
		Log:          slog.Default(),
		PollInterval: 1 * time.Second,
		Queue:        AiReportQ,
	})

	r.Register("ai-report", func(ctx context.Context, m []byte) error {
		return nil
	})

	r.Start(context.Background())
}
