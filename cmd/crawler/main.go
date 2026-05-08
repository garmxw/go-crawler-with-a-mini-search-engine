package main

import (
	"flag"
	"fmt"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/configs"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/cli"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/crawler"
)







func main () {

	crawlDelay := configs.ReadDelay()

	//flags
	var urls cli.MultiFlag
	flag.Var(&urls, "url", "Add URL (can be used multiple times)")
	fileFlag := flag.String("file", "", "File containing URLs to crawl")
	jsonFlag := flag.String("json", "", "JSON file with URLs") 	// example json {"urls": ["https://example.com", "https://example.org"]}
	depthFlag := flag.Int("depth", 0, "Depth of the crawl")
	maxPageFlag := flag.Int("maxPages", 3, "Maximum number of pages to crawl")
	pagesPathFlag := flag.String("storage", "data/pages", "Path to save the crawl results")

	flag.Parse()

	fmt.Println("\nGo Crawl !!")
	fmt.Println("urls:", urls)
	fmt.Println("file:", *fileFlag)
	fmt.Println("json:", *jsonFlag)
	fmt.Println("depth:", *depthFlag)
	fmt.Println("maxPages:", *maxPageFlag)
	fmt.Println("storage:", *pagesPathFlag)

	// run the crawler logic
	crawler.RunCrawler(
		urls,
		*depthFlag,
		*maxPageFlag,
		crawlDelay,
		*fileFlag,
		*jsonFlag,
		*pagesPathFlag,
	)

}
