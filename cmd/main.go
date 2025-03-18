package main

import (
	"log"
	"os"
	"parser/config"
	"parser/internal/logger"
	"parser/internal/proxy"
	"parser/internal/scraper"
	"parser/internal/storage"
)

func main() {
	log.Print("Starting scraper")

	cfg := config.MustLoad()
	log.Print("Config is loaded")
	log := logger.NewLogger("local")
	log.Info("Logger is created")

	scrap := scraper.NewSamokatScraper(cfg, log)
	log.Info("Scraper is created")
	parser := scraper.NewSamokatParser(cfg, log)
	log.Info("Parser is created")
	saver := storage.NewStorageJson(cfg, log)
	log.Info("Storage is created")
	proxyManager := proxy.NewProxyManager(cfg, log)
	log.Info("ProxyManager is created")

	path, err := saver.CreateFile("category")
	if err != nil {
		os.Exit(1)
	}

	file, err := saver.ReadFile(path)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	categories := cfg.GetCategories()

	proxyUrl := proxyManager.GetRandomProxy()
	scrap.ChangeProxy(proxyUrl)
	log.Info("Uses", "proxy", proxyUrl)
	htmlPage, err := scrap.ScrapeCategory(categories[0])
	if err != nil {
		os.Exit(1)
	}

	products := parser.ParseHTML(htmlPage)
	saver.Save(products, file)
}
