package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/configs"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/cli"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/crawler"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/ui"
)

func main() {
	crawlDelay := configs.ReadDelay()

	// Flags
	var urls cli.MultiFlag
	flag.Var(&urls, "url", "Add URL (can be used multiple times)")
	fileFlag := flag.String("file", "", "File containing URLs to crawl")
	jsonFlag := flag.String("json", "", "JSON file with URLs") // {"urls": ["https://example.com"]}
	depthFlag := flag.Int("depth", 0, "Depth of the crawl")
	maxPageFlag := flag.Int("maxPages", 3, "Maximum number of pages to crawl")
	pagesPathFlag := flag.String("storage", "data/pages", "Path to save the crawl results")

	flag.Parse()

	//  Banner
	ui.PrintCrawlerBanner()

	// Config panel
	urlDisplay := "(none)"
	if len(urls) > 0 {
		urlDisplay = urls[0]
		if len(urls) > 1 {
			urlDisplay += fmt.Sprintf("  (+%d more)", len(urls)-1)
		}
	}

	ui.PrintConfigPanel("Crawl Configuration", []ui.ConfigRow{
		{Key: "URLs", Value: strOr(urlDisplay, "(none)")},
		{Key: "File", Value: strOr(*fileFlag, "(none)")},
		{Key: "JSON", Value: strOr(*jsonFlag, "(none)")},
		{Key: "Depth", Value: strconv.Itoa(*depthFlag)},
		{Key: "Max Pages", Value: strconv.Itoa(*maxPageFlag)},
		{Key: "Delay (s)", Value: strconv.Itoa(crawlDelay)},
		{Key: "Storage", Value: *pagesPathFlag},
	})

	// Crawl
	spin := ui.NewSpinner("Crawling pages...")

	err := crawler.RunCrawler(
		urls,
		*depthFlag,
		*maxPageFlag,
		crawlDelay,
		*fileFlag,
		*jsonFlag,
		*pagesPathFlag,
	)

	if err != nil {
		spin.Fail("Crawl failed")
		ui.ErrorBox("Crawl encountered an error: " + err.Error())
		os.Exit(1)
	}

	spin.Done("Crawl complete!")
	ui.Divider()
	ui.SuccessBox("Pages saved to  →  " + *pagesPathFlag)
}

func strOr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
