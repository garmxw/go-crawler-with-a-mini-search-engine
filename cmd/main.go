package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/cli"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/indexer"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/local"
)

func main () {

	mode := flag.String("mode", "", "Mode: web or local")
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

	flag.Parse()

	fmt.Println("mode:", *mode)
	fmt.Println("query:", *query)
	fmt.Println("path:", *path)
	fmt.Println("urls:", urls)
	fmt.Println("file:", *fileFlag)
	fmt.Println("json:", *jsonFlag)
	fmt.Println("lang:", *langFlag)
	fmt.Println("depth:", *depthFlag)
	fmt.Println("maxPage:", *maxPageFlag)

	if *mode == "" || *query == "" {
		log.Fatal("mode and query are required")
		return
	}

	if *mode == "local" && *path == "" {
		log.Fatal("path is required for local mode")
		return
	}

	if *langFlag != "english" && *langFlag != "french" {
		slog.Warn("invalid language use -lang=english or -lang=french, defaulting to english")
		return
	}

	switch *mode {
	case "web":
		fmt.Println("Running web mode...")
		//later:
		// run crawler with flags
		// then index
	case "local":
		fmt.Println("Running local mode...")
		loader := local.NewLoader(*path)
		docs, err := loader.Load()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Loaded docs:", len(docs))
		idx := indexer.NewIndexer(*langFlag)
		for _, doc := range docs {
			if strings.TrimSpace(doc.Text) == "" {
				continue
			}
			idx.Add(doc.ID, doc.Text, doc.Path)
		}
		fmt.Println("Indexed docs:", idx.TotalDocs)
		idx.TotalDocs = len(docs)
		fmt.Println("Total docs:", idx.TotalDocs)
		fmt.Print("Index built!")
	default:
		log.Fatal("invalid mode use -mode=web or -mode=local")
		return
	}

}
