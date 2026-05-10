package main

import (
	"fmt"
	"os"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/crawler"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/ui"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/ui/tui"
)

func main() {
	ui.PrintCrawlerBanner()

	cfg, err := tui.RunCrawlerForm(tui.CrawlerConfig{
		MaxPages: 3,
		Delay:    2,
		Storage:  "data/pages",
	})
	if err != nil {
		ui.ErrorBox("Error: " + err.Error())
		os.Exit(1)
	}
	if !cfg.Submitted {
		fmt.Println()
		ui.Dim("Aborted.")
		os.Exit(0)
	}

	fmt.Println()
	spin := ui.NewSpinner("Crawling pages...")

	err = crawler.RunCrawler(
		cfg.URLs,
		cfg.Depth,
		cfg.MaxPages,
		cfg.Delay,
		cfg.File,
		cfg.JSON,
		cfg.Storage,
	)
	if err != nil {
		spin.Fail("Crawl failed")
		ui.ErrorBox("Error: " + err.Error())
		os.Exit(1)
	}

	spin.Done("Crawl complete!")
	ui.Divider()
	ui.SuccessBox("Pages saved to: " + cfg.Storage)
}
