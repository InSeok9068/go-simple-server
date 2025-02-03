package jobs

import (
	"simple-server/internal"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/robfig/cron/v3"
)

func SquashJob(c *cron.Cron) {
	_, _ = c.AddFunc("* * * * *", SquashExecute)
}

func SquashExecute() {
	on := internal.EnvMap["CHROMEDP_HEADLESS"]

	u := launcher.New().
		Headless(on == "true").
		Leakless(false).
		MustLaunch()

	page := rod.New().ControlURL(u).MustConnect().MustPage("https://www.auc.or.kr/sign/in/base/user")
	page.MustClose()
}
