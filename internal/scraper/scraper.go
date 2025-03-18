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
	config   config.ConfigProvider
	logger   *slog.Logger
	proxyUrl string
}

func NewSamokatScraper(cfg config.ConfigProvider, log *slog.Logger) *samokatScraper {
	return &samokatScraper{
		config: cfg,
		logger: log,
	}
}

func (ss *samokatScraper) ChangeProxy(newProxyUrl string) {
	ss.proxyUrl = newProxyUrl
}

func (ss *samokatScraper) ScrapeCategory(url string) (string, error) {
	op := "scraper.SamokatScraper.ScraperCategory"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", ss.config.GetOptHeadless()),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"),
	)

	if ss.proxyUrl != "" {
		opts = append(opts, chromedp.Flag("proxy-server", ss.proxyUrl))
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ss.logger.Info("Trying to enter... ", "op", op)

	var html string
	delay, err := time.ParseDuration(ss.config.GetRequestDelay())
	if err != nil {
		ss.logger.Error("Failed convert request delay from string to int", logger.Err(err), "op", op)
		delay = 5 * time.Second
	}

	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(delay),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		ss.logger.Error("Error in run to website: ", logger.Err(err), "op", op)
		return "", err
	}

	return html, nil
}
