package jobs

import (
	"log/slog"
	"simple-server/internal/config"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/robfig/cron/v3"
)

func SquashJob(c *cron.Cron) {
	if _, err := c.AddFunc("* * * * *", SquashExecute); err != nil {
		slog.Error("스케줄 등록 실패", "error", err)
	}
}

func SquashExecute() {
	on := config.EnvMap["CHROMEDP_HEADLESS"]

	u := launcher.New().
		Headless(on == "true").
		Leakless(false).
		MustLaunch()

	page := rod.New().ControlURL(u).MustConnect().MustPage("https://www.auc.or.kr/sign/in/base/user")
	page.MustClose()
}
