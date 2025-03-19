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

	proxyManager := proxy.NewProxyManager(log)
	log.Info("ProxyManager is created")

	if err := proxyManager.LoadProxy(); err != nil {
		os.Exit(1)
	}

	scrap := scraper.NewSamokatScraper(cfg, log, proxyManager)
	log.Info("Scraper is created")
	parser := scraper.NewSamokatParser(cfg, log)
	log.Info("Parser is created")
	saver := storage.NewStorageJson(cfg, log)
	log.Info("Storage is created")

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

	htmlPage, err := scrap.ScrapeCategory(categories[0])
	if err != nil {
		os.Exit(1)
	}

	products := parser.ParseHTML(htmlPage)
	saver.ClearFile(path)
	saver.Save(products, file)
}
