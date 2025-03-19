package scraper

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"parser/config"
	"parser/internal/logger"
	"parser/internal/proxy"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/tebeka/selenium"
)

const (
	seleniumHost = "http://localhost:4444/wd/hub"
)

type Scraper interface {
	ScrapeCategory(url string) (string, error)
}

type samokatScraper struct {
	config config.ConfigProvider
	logger *slog.Logger
	pm     *proxy.ProxyManager
}

func NewSamokatScraper(cfg config.ConfigProvider, log *slog.Logger, pm *proxy.ProxyManager) *samokatScraper {
	return &samokatScraper{
		config: cfg,
		logger: log,
		pm:     pm,
	}
}

func (ss *samokatScraper) ScrapeCategory(url string) (string, error) {
	op := "scraper.samokatScraper.ScrapeCategory"
	ss.logger.Info("Start selenium driver", "op", op)
	chromeDriverPath := "/usr/bin/chromedriver"
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
	userAgentFlag := fmt.Sprintf("--user-agent=%s", userAgent)
	service, err := selenium.NewChromeDriverService(chromeDriverPath, 4444)
	if err != nil {
		ss.logger.Error("Error starting ChromeDriver", logger.Err(err), "op", op)
	}
	defer service.Stop()

	prox := ss.pm.GetRandomProxy()
	if prox == nil {
		ss.logger.Error("No proxies, session is terminated", "op", op)
		return "", err
	}

	ss.logger.Info("Uses proxy", "host", prox.Host, "port", prox.Port)
	time.Sleep(2 * time.Second)

	if err := ss.pm.SetProxyExt(prox.ProtocolType, prox.Host, prox.Port, prox.Username, prox.Password); err != nil {
		ss.logger.Error("Failed set proxt extentsion", logger.Err(err), "op", op)
		return "", err
	}

	extensionFile, err := os.ReadFile(proxy.ExtensionFilePath)
	if err != nil {
		ss.logger.Error("Failed in read file extension", logger.Err(err), "op", op)
		return "", err
	}

	encodedExtension := base64.StdEncoding.EncodeToString(extensionFile)

	caps := selenium.Capabilities{
		"browserName": "chrome",
		"goog:chromeOptions": map[string]interface{}{
			"extensions": []string{
				encodedExtension,
			},
			"args": []string{
				userAgentFlag,
				"--headless",
				"--disable-blink-features=AutomationControlled",
				"--no-sandbox",
				"--ignore-certificate-errors",
				"--start-maximized",
				"--disable-gpu",
			},
		},
	}

	wd, err := selenium.NewRemote(caps, seleniumHost)
	if err != nil {
		ss.logger.Error("Failed connecting to webdriver", logger.Err(err), "op", op)
		return "", err
	}
	defer func() {
		ss.logger.Info("Quitting webdriver")
		if err := wd.Quit(); err != nil {
			ss.logger.Error("Error quitting webdriver", logger.Err(err), "op", op)
		}
	}()

	if err = wd.Get(url); err != nil {
		ss.logger.Error("Failet visit in site", logger.Err(err), "op", op)
		return "", err
	}

	ss.logger.Info("Visited in", "site", url, "op", op)

	delay, err := time.ParseDuration(ss.config.GetRequestDelay())
	if err != nil {
		ss.logger.Error("Failed parse Request delay to time.Duration", logger.Err(err), "op", op)
		delay = 5 * time.Second
	}

	time.Sleep(delay)

	html, err := wd.PageSource()
	if err != nil {
		ss.logger.Error("Failed to get source page", logger.Err(err), "op", op)
	}

	return html, nil
}

/*
Don't work with my proxy.
Error: net:ERR_NO_SUPPORTED_PROXIES
*/
func (ss *samokatScraper) ScrapeCategoryWithChromedp(url string) (string, error) {
	op := "scraper.SamokatScraper.ScraperCategory"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", ss.config.GetOptHeadless()),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"),
	)

	// if ss.proxyUrl != "" {
	// 	opts = append(opts, chromedp.Flag("proxy-server", ss.proxyUrl))
	// }

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

// func setProxyAuth(username, password string) chromedp.Action {
// 	return chromedp.ActionFunc(func(ctx context.Context) error {
// 		executor, ok := chromedp.FromContext(ctx)
// 		if !ok {
// 			return errors.New("failed to get executor")
// 		}

// 		err := executor.Execute(ctx, "Network.enable", nil)
// 		if err != nil {
// 			return err
// 		}

// 		err = executor.Execute(ctx, "Network.setExtraHTTPHeaders", map[string]any{
// 			"headers": map[string]string{
// 				"Proxy-Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password)),
// 			},
// 		})
// 		return err
// 	})
// }
