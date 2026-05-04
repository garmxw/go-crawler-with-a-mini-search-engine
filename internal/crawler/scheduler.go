package crawler

import (
	"sync"
)

type Scheduler struct {
	queue chan string
	visited map[string]bool
	mu  sync.Mutex
}

func NewScheduler (buffer int) *Scheduler {
	return &Scheduler{
		queue: make(chan string, buffer),
		visited: make(map[string]bool),
	}
}

func (s *Scheduler) Add(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.visited[url] {
		s.visited[url] = true
		s.queue <- url
	}
}

func (s *Scheduler) Next() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case url := <-s.queue:
		return url
	default:
		return ""
	}
}
func (s *Scheduler) IsVisited(url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.visited[url]
}
