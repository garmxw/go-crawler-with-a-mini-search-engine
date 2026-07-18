package main

import (
	"fmt"
	"log"
	"os"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/models"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/search"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/ui"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/ui/tui"
)

func main() {
	// 1. Banner prints once and stays on screen above the form
	ui.PrintSearchBanner()

	// 2. Interactive form renders inline below the banner
	cfg, err := tui.RunSearchForm(tui.SearchConfig{
		Mode:     "local",
		Lang:     "english",
		Limit:    5,
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
	//    Now run the actual search logic spinner + results print here below.
	fmt.Println()

	var results []models.SearchResult

	switch cfg.Mode {
	case "web":
		ui.Info("Loading stored pages...")
		spin := ui.NewSpinner("Indexing documents...")
		results, err = search.RunWebMode(cfg.Query, cfg.Lang, cfg.Storage, cfg.Detailed)
		if err != nil {
			spin.Fail("Web mode failed")
			log.Fatal(err)
		}
		spin.Done("Index built!")

	case "local":
		ui.Info("Loading local documents...")
		spin := ui.NewSpinner("Indexing documents...")
		results, err = search.RunLocalMode(cfg.Path, cfg.Query, cfg.Lang, cfg.Detailed)
		if err != nil {
			spin.Fail("Local mode failed")
			log.Fatal(err)
		}
		spin.Done("Index built!")

	case "live":
		ui.Info("Starting live crawl + index pipeline...")
		spin := ui.NewSpinner("Crawling the web...")
		results, err = search.RunWebLiveMode(
			cfg.Query,
			cfg.Lang,
			cfg.URLs,
			cfg.Depth,
			cfg.MaxPages,
			cfg.Delay,
			cfg.File,
			cfg.JSON,
			cfg.Detailed,
			cfg.Storage,
		)
		if err != nil {
			spin.Fail("Live mode failed")
			log.Fatal(err)
		}
		spin.Done("Crawl + index complete!")

	default:
		ui.ErrorBox("Unknown mode: " + cfg.Mode)
		os.Exit(1)
	}

	// Trim results
	if cfg.Limit > 0 && len(results) > cfg.Limit {
		results = results[:cfg.Limit]
	}

	// Display results prints below everything
	uiResults := make([]ui.SearchResult, len(results))
	for i, r := range results {
		uiResults[i] = ui.SearchResult{
			DocID: r.DocID,
			Path:  r.Path,
			Score: r.Score,
		}
	}
	ui.PrintSearchResults(uiResults, cfg.Query)
	ui.Divider()
}
