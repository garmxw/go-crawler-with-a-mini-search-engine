package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/crawler"
)

func main () {
	delayStr := os.Getenv("CrawlDelay")
	crawlDelay,err := strconv.Atoi(delayStr)
	if err != nil {
		crawlDelay = 2
		fmt.Println("Error: CrawlerDelay is not a valid value")
	}
	f := crawler.NewFetcher(int64(crawlDelay))
	f.Fetch("https://example.com")
}
