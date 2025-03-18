package main

import (
	"fmt"
	"log"
	"parser/config"
	"parser/internal/logger"
	"parser/internal/scraper"
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
	categories := cfg.GetCategories()

	htmlPage, err := scrap.ScrapeCategory(categories[0])
	if err != nil {
		fmt.Println("Something is wrong")
		return
	}

	products := parser.ParseHTML(htmlPage)
	fmt.Println(products)
}
