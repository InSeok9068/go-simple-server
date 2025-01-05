package jobs

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/robfig/cron/v3"
)

func SquashJob(c *cron.Cron) {
	c.AddFunc("* * * * *", SquashExecute)
}

func SquashExecute() {
	// Chromedp 컨텍스트 생성
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var result string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.auc.or.kr/sign/in/base/user"),
		chromedp.OuterHTML("html", &result),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
