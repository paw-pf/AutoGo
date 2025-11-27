package browser

import (
	"autogo/config"
	"context"

	"github.com/chromedp/chromedp"
)

// NewContext создает контекст с настроенным chromedp
func NewContext() (context.Context, context.CancelFunc, error) {
	ctx, _ := context.WithCancel(context.Background())

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx,
		chromedp.Flag("headless", config.Headless),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("start-maximized", ""),
	)

	chromeCtx, chromeCancel := chromedp.NewContext(allocCtx)
	return chromeCtx, func() {
		chromeCancel()
		allocCancel()
	}, nil

}
