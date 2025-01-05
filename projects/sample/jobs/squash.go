package jobs

import (
	"context"
	"simple-server/internal"

	"github.com/chromedp/chromedp"
	"github.com/robfig/cron/v3"
)

func SquashJob(c *cron.Cron) {
	c.AddFunc("* * * * *", SquashExecute)
}

func SquashExecute() {
	on := internal.EnvMap["CHROMEDP_HEADLESS"]

	// Chromedp 컨텍스트 생성
	var ctx context.Context
	var cancel context.CancelFunc
	if on == "false" {
		// Chromedp 컨텍스트 옵션
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),          // 헤드리스 모드 OFF
			chromedp.Flag("disable-gpu", false),       // GPU 활성화
			chromedp.Flag("window-size", "1920,1080"), // 브라우저 창 크기 설정
		)
		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		ctx, cancel = chromedp.NewContext(allocCtx)
		defer cancel()
	} else {
		ctx, cancel = chromedp.NewContext(context.Background())
		defer cancel()
	}

	var result string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.auc.or.kr/sign/in/base/user"),
		chromedp.OuterHTML("html", &result),
	)

	if err != nil {
		panic(err)
	}

	// fmt.Println(result)
}
