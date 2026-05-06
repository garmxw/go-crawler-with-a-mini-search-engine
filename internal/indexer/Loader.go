package indexer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)


type StoredPage struct {
	ID    int
	URL   string
	Title string
	Date  time.Time
	Text  string

}
func LoadAndIndex(path string, idx *Indexer) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		fullPath := filepath.Join(path, f.Name())

		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		var page StoredPage
		err = json.Unmarshal(data, &page)
		if err != nil {
			continue
		}

		idx.Add(page.ID, page.Text, fullPath)
	}

	return nil
}
