package scraper

import (
	"log/slog"
	"parser/config"
	"parser/internal/logger"
	"parser/internal/models"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ParserHTML interface {
	ParseHTML(html string) []models.Product
}

type SamokatParser struct {
	config config.ConfigProvider
	logger *slog.Logger
}

func NewSamokatParser(cfg config.ConfigProvider, log *slog.Logger) *SamokatParser {
	return &SamokatParser{
		config: cfg,
		logger: log,
	}
}

func (sp *SamokatParser) ParseHTML(html string) []models.Product {
	op := "scraper.SamokatParser.ParseHTML"

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		sp.logger.Error("Failed to parse HTML: %v", logger.Err(err), "op", op)
		return nil
	}

	// maybe worst: if classes in html is changes then scraping needed update
	productNameCard := ".ProductCard_name__2VDcL"
	productPriceCard := ".ProductCardActions_text__3Uohy"
	productLinkClass := "a[href^='/product/']"

	var products []models.Product

	var productNames, productPrices, productLinks []string

	doc.Find(productNameCard).Each(func(i int, s *goquery.Selection) {
		productNames = append(productNames, strings.TrimSpace(s.Text()))
	})

	doc.Find(productPriceCard).Each(func(i int, s *goquery.Selection) {
		productPrices = append(productPrices, strings.TrimSpace(s.Text()))
	})

	doc.Find(productLinkClass).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			productLinks = append(productLinks, "https://samokat.ru"+href)
		}
	})

	for i := 0; i < len(productNames) && i < len(productPrices) && i < len(productLinks); i++ {
		product := models.Product{
			Name:  productNames[i],
			Price: productPrices[i],
			Link:  productLinks[i],
		}

		products = append(products, product)
	}

	sp.logger.Info("All products is parsed and retrived", "op", op)

	return products
}
