package scraper

import (
	"context"
	"log/slog"
	"parser/config"
	"parser/internal/logger"
	"time"

	"github.com/chromedp/chromedp"
)

type Scraper interface {
	ScrapeCategory(url string) (string, error)
}

type samokatScraper struct {
	Config config.ConfigProvider
	logger *slog.Logger
}

func NewSamokatScraper(cfg config.ConfigProvider, log *slog.Logger) *samokatScraper {
	return &samokatScraper{
		Config: cfg,
		logger: log,
	}
}

func (ss *samokatScraper) ScrapeCategory(url string) (string, error) {
	op := "scraper.SamokatScraper.ScraperCategory"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", ss.Config.GetOptHeadless()),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ss.logger.Info("Created chromedp new context", "op", op)

	ss.logger.Info("trying to enter... ", "op", op)

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &html),
	)

	if err != nil {
		ss.logger.Error("Error in run to website: ", logger.Err(err), "op", op)
		return "", err
	}

	return html, nil
}
