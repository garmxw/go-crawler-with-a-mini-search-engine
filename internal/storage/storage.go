package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)


type Page struct {
	ID int
	URL string
	Title string
	Date time.Time
	Text string

}

type Storage struct {
	mu *sync.Mutex
	nextID int
	path string
}
// NewStorage creates a new Storage instance with the given path.
func NewStorage(path string) *Storage {
	os.MkdirAll(path, os.ModePerm)

	return &Storage{
		mu:     &sync.Mutex{},
		nextID: 1,
		path:   path,
	}
}

func (s *Storage) SavePage(page Page) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	page.ID = s.nextID
	s.nextID++
	filename := filepath.Join(
		s.path,
		fmt.Sprintf("%d.json", page.ID),
)

	file, err := os.Create(filename)
	if err != nil {
		return -1
	}
	defer file.Close()

    encode := json.NewEncoder(file)
    if err := encode.Encode(page); err != nil {
    	return -1
    }
	return page.ID
}

func (s *Storage) LoadPages() ([]Page, error) {
	var pages []Page
	entries, err := os.ReadDir(s.path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
        if entry.IsDir(){
            continue
        }
        filePath := filepath.Join(s.path, entry.Name())
        file, err := os.Open(filePath)
        if err != nil {
            continue
        }
        var page Page
        decode := json.NewDecoder(file)
        if err := decode.Decode(&page); err != nil {
            file.Close()
            continue
        }
        file.Close()
        pages = append(pages, page)
	}
	return pages, nil
}
