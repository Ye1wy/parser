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
	log := logger.NewLogger("local")

	scrap := scraper.NewSamokatScraper(cfg, log)
	parser := scraper.NewSamokatParser(cfg, log)
	categories := "https://samokat.ru/category/vsya-gotovaya-eda-6"

	htmlPage, err := scrap.ScrapeCategory(categories)
	if err != nil {
		fmt.Println("Something is wrong")
		return
	}

	products := parser.ParseHTML(htmlPage)
	fmt.Println(products)
}
