package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// url := "https://lenta.com"
	url := "https://samokat.ru"
	fmt.Println("trying to enter...")

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(10*time.Second),
		chromedp.OuterHTML("html", &html),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Website loaded successfully:", html)
}
