package storage

import (
	"encoding/json"
	"fmt"
	"os"
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
	filename := fmt.Sprintf("%s%d.json", s.path, page.ID)

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
