package search

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/crawler"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/indexer"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/local"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/models"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/storage"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/utils"
)



func RunLocalMode(
	path string,
	query string,
	lang string,
	detailed bool,
) ([]models.SearchResult, error) {
	loader := local.NewLoader(path)
	docs, err := loader.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nLoaded docs:", len(docs))
	fmt.Println("Indexing...")

	idx := indexer.NewIndexer(lang)
	for _, doc := range docs {
		if strings.TrimSpace(doc.Text) == "" {
			continue
		}
		idx.Add(doc.ID, doc.Text, doc.Path)
	}
	fmt.Println("\nIndexed docs:", idx.TotalDocs)
	fmt.Println("Total docs:", idx.TotalDocs)
	idx.Build()
	fmt.Println("Index built!")
	if detailed {
		fmt.Println("\n/| /| /| Index Detailed List")
		utils.PrintIndexDetails(idx)
		utils.PrintDocStats(idx)

	}

	results := Search(query, idx, lang)

	return results, nil
}

func RunWebMode(
	query string,
	lang string,
	storagePath string,
	detailed bool,
) ([]models.SearchResult, error) {
	// Create a new storage instance which loads pages from the given path
	store := storage.NewStorage(storagePath)
	// Load pages from storage
	pages, err := store.LoadPages()
	if err != nil {
		return nil, err
	}
	fmt.Println("Loaded pages:", len(pages))
	idx := indexer.NewIndexer(lang)
	for _,page := range pages {
		idx.Add(page.ID, page.Text, page.URL)
		fmt.Println(page.URL)
	}
	fmt.Println("\nIndexed docs:", idx.TotalDocs)
	fmt.Println("Total docs:", idx.TotalDocs)
	idx.Build()
	fmt.Println("Index built!")
	if detailed {
		fmt.Println("\n/| /| /| Index Detailed List")
		utils.PrintIndexDetails(idx)
		utils.PrintDocStats(idx)

	}
	results := Search(query, idx, lang)
	return results, nil
}

func RunWebLiveMode(
	query string,
	lang string,
	urls []string,
	depth int,
	maxPages int,
	delay int,
	filePath string,
	jsonFilePath string,
	detailed bool,
	storagePath string,
) ([]models.SearchResult, error) {

	// clear old pages first
	os.RemoveAll(storagePath)

	err := os.MkdirAll(
		storagePath,
		os.ModePerm,
	)

	if err != nil {
		return nil, err
	}

	// crawl first
	err = crawler.RunCrawler(
		urls,
		depth,
		maxPages,
		delay,
		filePath,
		jsonFilePath,
		storagePath,
	)

	if err != nil {
		return nil, err
	}

	// THEN search stored pages
	return RunWebMode(
		query,
		lang,
		storagePath,
		detailed,
	)
}
