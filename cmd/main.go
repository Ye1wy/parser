package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	originUrl := "https://samokat.ru"
	url := originUrl + "/category/vsya-gotovaya-eda-6"
	fmt.Println("trying to enter...")

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &html),
	)

	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	// maybe worst: if classes in html is changes then scraping needed update
	productNameCard := ".ProductCard_name__2VDcL"
	productPriceCard := ".ProductCardActions_text__3Uohy"
	productLinkClass := "a[href^='/product/']"

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
			productLinks = append(productLinks, originUrl+href)
		}
	})

	fmt.Println("Extracted data: ")
	fmt.Println(productNames)
	fmt.Println(productPrices)
	fmt.Println(productLinks)
}
