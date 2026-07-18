package main

import (
	"fmt"
	"os"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/crawler"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/ui"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/ui/tui"
)

func main() {
	// 1. Banner prints once and stays on screen above the form
	ui.PrintCrawlerBanner()

	// 2. Interactive form renders inline below the banner
	cfg, err := tui.RunCrawlerForm(tui.CrawlerConfig{
		MaxPages: 3,
		Delay:    2,
		Storage:  "data/pages",
	})
	if err != nil {
		ui.ErrorBox("TUI error: " + err.Error())
		os.Exit(1)
	}
	if !cfg.Submitted {
		fmt.Println()
		ui.Dim("Aborted.")
		os.Exit(0)
	}

	// 3. Confirm panel already printed inside the TUI (below the locked form).
	// Now run the actual crawler spinner + result print here below everything.
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
		ui.ErrorBox("Crawl encountered an error: " + err.Error())
		os.Exit(1)
	}

	spin.Done("Crawl complete!")
	ui.Divider()
	ui.SuccessBox("Pages saved to  →  " + cfg.Storage)
}
