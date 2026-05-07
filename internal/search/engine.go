package search

import (
	"fmt"
	"log"
	"strings"

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
	store := storage.NewStorage(storagePath)
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
