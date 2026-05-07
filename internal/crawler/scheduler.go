package crawler

import (
	"sync"
)

type URLItem struct {
	URL   string
	Depth int
}

type Scheduler struct {
	queue    chan URLItem
	visited  map[string]bool
	maxDepth int
	maxPages int
	count    int
	mu       sync.Mutex
}


// GetCount returns the current number of pages processed
func (s *Scheduler) GetCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.count
}

// GetMaxPages returns the limit
func (s *Scheduler) GetMaxPages() int {
	return s.maxPages
}

func NewScheduler(buffer int, maxDepth int, maxPages int) *Scheduler {
	return &Scheduler{
		// Use a large enough buffer to prevent early deadlocks
		queue:    make(chan URLItem, buffer),
		visited:  make(map[string]bool),
		maxDepth: maxDepth,
		maxPages: maxPages,
	}
}

func (s *Scheduler) Add(rawURL string, depth int) {
	normalized := NormalizeURL(rawURL)
	if normalized == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Check Limits
	if (s.maxPages > 0 && s.count >= s.maxPages) || depth > s.maxDepth {
		return
	}

	// 2. Duplication Check
	if s.visited[normalized] {
		return
	}

	// 3. Register and Enqueue
	s.visited[normalized] = true
	s.count++

	// Non-blocking send to prevent the Fetcher from hanging if the buffer is full
	select {
	case s.queue <- URLItem{URL: normalized, Depth: depth}:
	default:
		// If buffer is full, we log it or ignore it to keep the crawler moving
	}
}

func (s *Scheduler) Next() (URLItem, bool) {
	item, ok := <-s.queue
	return item, ok
}
