package crawler

import (
	"sync"
)

type URLItem struct {
	URL string
	Depth int
}


type Scheduler struct {
	queue chan URLItem
	visited map[string]bool
	maxDepth int
	maxPages int
	count int
	mu  sync.Mutex
}

func NewScheduler (buffer int, maxDepth int, maxPages int ) *Scheduler {
	return &Scheduler{
		queue: make(chan URLItem, buffer),
		visited: make(map[string]bool),
		maxDepth: maxDepth,
		maxPages: maxPages,
	}
}

func (s *Scheduler) Add(url string, depth int) {
	url = NormalizeURL(url)
	s.mu.Lock()
	defer s.mu.Unlock()
	// Check if we've already hit the page limit
	if s.maxPages > 0 && s.count >= s.maxPages {
		return
	}
	if !s.visited[url] && depth <= s.maxDepth {
		s.visited[url] = true
		s.count++
		s.queue <- URLItem{URL: url, Depth: depth}

	}
}

func (s *Scheduler) Next() URLItem {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case url := <-s.queue:
		return url
	default:
		return URLItem{}
	}
}
func (s *Scheduler) IsVisited(url string) bool {
	url = NormalizeURL(url)
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.visited[url]
}

func (s *Scheduler) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.queue) == 0
}

func (s *Scheduler) Close() {
	close(s.queue)
}
