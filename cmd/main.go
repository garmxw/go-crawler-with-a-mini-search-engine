package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/cli"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/search"
)

func main () {

	mode := flag.String("mode", "local", "Mode: web or local")
	query := flag.String("query", "", "Search query")
	//local
	path := flag.String("path", "", "Path to local file")
	// web (reuse crwaler's flags)
	var urls cli.MultiFlag
	flag.Var(&urls, "url", "URL to crawl")
	fileFlag := flag.String("file", "", "file with URLs")
	jsonFlag := flag.String("json", "", "JSON file with URLs")
	langFlag := flag.String("lang", "english", "Language of the text (english or french)")
	depthFlag := flag.Int("depth", 2, "Depth of crawl")
	maxPageFlag := flag.Int("maxPage", 0, "Max number of pages to crawl")
	limitFlag := flag.Int("limit", 1, "Max number of results to return")
	detailedFlag := flag.Bool("detailed", false, "Detailed output")
	flag.Parse()

	fmt.Println("\nmode:", *mode)
	fmt.Println("query:", *query)
	fmt.Println("path:", *path)
	fmt.Println("urls:", urls)
	fmt.Println("file:", *fileFlag)
	fmt.Println("json:", *jsonFlag)
	fmt.Println("lang:", *langFlag)
	fmt.Println("depth:", *depthFlag)
	fmt.Println("maxPage:", *maxPageFlag)
	fmt.Println("limit:", *limitFlag)
	fmt.Println("detailed:", *detailedFlag)

	if *mode == "" || *query == "" {
		log.Fatal("mode and query are required")
		return
	}

	if *mode != "web" && *mode != "local" {
		log.Fatal("invalid mode only web or local are supported")
		return
	}

	if *mode == "local" && *path == "" {
		log.Fatal("path is required for local mode")
		return
	}

	if *limitFlag <= 0 {
    slog.Warn("Warning: limit must be positive. Defaulting to 1.")
    *limitFlag = 1
}


	if *langFlag != "english" && *langFlag != "french" {
		slog.Warn(
			"invalid language, defaulting to english",
			"provided", *langFlag,
		)

		*langFlag = "english"
	}

	switch *mode {
	case "web":
		fmt.Println("\nRunning web mode...")
		results, err := search.RunWebMode(*query, *langFlag, "data/storage",*detailedFlag)
		if err != nil {
			log.Fatal(err)
		}
		if *limitFlag > 0 && len(results) > *limitFlag {
			results = results[:*limitFlag]
		}

		for i, r := range results {
			fmt.Printf(
				"%d. Score: %.4f | Path: %s\n",
				i+1,
				r.Score,
				r.Path,
			)
		}
	case "local":
		fmt.Println("\nRunning local mode...")
		results, err := search.RunLocalMode(*path, *query, *langFlag, *detailedFlag)
		if err != nil {
			log.Fatal(err)
		}
		if *limitFlag > 0 && len(results) > *limitFlag {
			results = results[:*limitFlag]
		}

		fmt.Println("\nResults:")
		for i, r := range results {
			fmt.Printf(
				"%d. Score: %.4f | Path: %s\n",
				i+1,
				r.Score,
				r.Path,
				)
		}
	default:
		log.Fatal("invalid mode use -mode=web or -mode=local")
		return
	}

}
