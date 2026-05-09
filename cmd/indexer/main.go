package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/configs"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/cli"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/models"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/search"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/ui"
)

func main() {
	crawlDelay := configs.ReadDelay()

	// Flags
	mode := flag.String("mode", "local", "Mode: web | local | live")
	query := flag.String("query", "", "Search query")
	path := flag.String("path", "", "Path to local files  (local mode)")

	var urls cli.MultiFlag
	flag.Var(&urls, "url", "URL to crawl  (live mode, repeatable)")
	langFlag := flag.String("lang", "english", "Language: english | french")
	limitFlag := flag.Int("limit", 1, "Max results to return")
	detailedFlag := flag.Bool("detailed", false, "Print detailed index stats")

	fileFlag := flag.String("file", "", "File containing URLs  (live mode)")
	jsonFlag := flag.String("json", "", "JSON file with URLs  (live mode)")
	depthFlag := flag.Int("depth", 0, "Crawl depth  (live mode)")
	maxPageFlag := flag.Int("maxPages", 3, "Max pages to crawl  (live mode)")
	pagesPathFlag := flag.String("storage", "data/pages", "Pages storage path")

	flag.Parse()

	//  Banner
	ui.PrintSearchBanner()

	// Validation
	if *mode == "" || *query == "" {
		ui.ErrorBox("-mode and -query flags are required.")
		os.Exit(1)
	}

	if *mode != "web" && *mode != "local" && *mode != "live" {
		ui.ErrorBox("Invalid mode: \"" + *mode + "\".  Use: web | local | live")
		os.Exit(1)
	}

	if *mode == "local" && *path == "" {
		ui.ErrorBox("-path is required when using local mode.")
		os.Exit(1)
	}

	if *limitFlag <= 0 {
		slog.Warn("limit must be positive, defaulting to 1")
		*limitFlag = 1
	}

	if *langFlag != "english" && *langFlag != "french" {
		ui.Warn("Invalid language \"" + *langFlag + "\", defaulting to english.")
		*langFlag = "english"
	}

	//Mode badge + config panel
	ui.PrintModeBadge(*mode)

	rows := []ui.ConfigRow{
		{Key: "Query", Value: "\"" + *query + "\""},
		{Key: "Language", Value: *langFlag},
		{Key: "Limit", Value: strconv.Itoa(*limitFlag)},
		{Key: "Detailed", Value: strconv.FormatBool(*detailedFlag)},
	}
	switch *mode {
	case "local":
		rows = append(rows, ui.ConfigRow{Key: "Path", Value: *path})
	case "web":
		rows = append(rows, ui.ConfigRow{Key: "Storage", Value: *pagesPathFlag})
	case "live":
		urlDisplay := "(none)"
		if len(urls) > 0 {
			urlDisplay = urls[0]
			if len(urls) > 1 {
				urlDisplay += " (+more)"
			}
		}
		rows = append(rows,
			ui.ConfigRow{Key: "URLs", Value: urlDisplay},
			ui.ConfigRow{Key: "Depth", Value: strconv.Itoa(*depthFlag)},
			ui.ConfigRow{Key: "Max Pages", Value: strconv.Itoa(*maxPageFlag)},
			ui.ConfigRow{Key: "Delay (s)", Value: strconv.Itoa(crawlDelay)},
			ui.ConfigRow{Key: "Storage", Value: *pagesPathFlag},
		)
	}
	ui.PrintConfigPanel("Search Configuration", rows)

	// Run
	var results []models.SearchResult
	var err error

	switch *mode {
	case "web":
		ui.Info("Loading stored pages...")
		spin := ui.NewSpinner("Indexing documents...")
		results, err = search.RunWebMode(*query, *langFlag, *pagesPathFlag, *detailedFlag)
		if err != nil {
			spin.Fail("Web mode failed")
			log.Fatal(err)
		}
		spin.Done("Index built!")

	case "local":
		ui.Info("Loading local documents...")
		spin := ui.NewSpinner("Indexing documents...")
		results, err = search.RunLocalMode(*path, *query, *langFlag, *detailedFlag)
		if err != nil {
			spin.Fail("Local mode failed")
			log.Fatal(err)
		}
		spin.Done("Index built!")

	case "live":
		ui.Info("Starting live crawl → index pipeline...")
		spin := ui.NewSpinner("Crawling the web...")
		results, err = search.RunWebLiveMode(
			*query,
			*langFlag,
			urls,
			*depthFlag,
			*maxPageFlag,
			crawlDelay,
			*fileFlag,
			*jsonFlag,
			*detailedFlag,
			*pagesPathFlag,
		)
		if err != nil {
			spin.Fail("Live mode failed")
			log.Fatal(err)
		}
		spin.Done("Crawl + index complete!")

	default:
		ui.ErrorBox("Invalid mode. Use -mode=web | local | live")
		os.Exit(1)
	}

	// Trim results
	if *limitFlag > 0 && len(results) > *limitFlag {
		results = results[:*limitFlag]
	}

	// Display results
	uiResults := make([]ui.SearchResult, len(results))
	for i, r := range results {
		uiResults[i] = ui.SearchResult{
			DocID: r.DocID,
			Path:  r.Path,
			Score: r.Score,
		}
	}

	ui.PrintSearchResults(uiResults, *query)
	ui.Divider()
}
