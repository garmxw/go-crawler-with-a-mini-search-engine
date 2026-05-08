package crawler

import (
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/storage"
	"github.com/gocolly/colly"
)

type Fetcher struct {
	collector *colly.Collector
	Scheduler *Scheduler
	storage   *storage.Storage
	MaxDepth  int
	active    atomic.Int32 // Safely tracks in-flight requests
}

func NewFetcher(s *Scheduler, st *storage.Storage, domains []string, delay int64, maxDepth int) *Fetcher {
	c := colly.NewCollector(
		colly.AllowedDomains(domains...),
		colly.Async(true),
	)

	// Enable Colly's internal duplicate prevention
	c.AllowURLRevisit = false

	var randomDelay time.Duration
	if delay > 0 {
		randomDelay = time.Duration(delay/2) * time.Second
	}

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       time.Duration(delay) * time.Second,
		RandomDelay: randomDelay,
	})

	f := &Fetcher{
		collector: c,
		Scheduler: s,
		storage:   st,
		MaxDepth:  maxDepth,
	}

	f.registerCallbacks()

	slog.Info("Fetcher initialized",
		"domains", domains,
		"max_depth", maxDepth,
		"parallelism", 2,
	)

	return f
}

func (f *Fetcher) registerCallbacks() {
	// Request started
	f.collector.OnRequest(func(r *colly.Request) {
		f.active.Add(1) // Increment active workers
		r.Headers.Set("User-Agent", "Mozilla/5.0 (compatible; GoCrawler/1.0; +garmxBot)")
		slog.Debug("Visiting", "url", r.URL.String())
	})

	// Request finished successfully
	f.collector.OnScraped(func(r *colly.Response) {
		f.active.Add(-1) // Decrement active workers
	})

	// Request failed
	f.collector.OnError(func(r *colly.Response, err error) {
		f.active.Add(-1) // Decrement active workers
		slog.Error("Request failed", "url", r.Request.URL, "error", err)
	})

	// Parse the HTML
	f.collector.OnHTML("html", func(e *colly.HTMLElement) {
		depth := 0
		if d := e.Request.Ctx.GetAny("depth"); d != nil {
			depth = d.(int)
		}

		page := ParserPage(e)

		id := f.storage.SavePage(storage.Page{
			URL:   page.URL,
			Title: page.Title,
			Date:  time.Now(),
			Text:  page.Text,
		})

		slog.Info("Page processed", "id", id, "title", page.Title, "depth", depth)

		if depth >= f.MaxDepth {
			return
		}

		for _, link := range page.Links {
			f.Scheduler.Add(link, depth+1)
		}
	})
}

// Change the definition to remove (seedURL string)
func (f *Fetcher) Start() {
	slog.Info("Crawler engine started")

	for {
		// Check if we reached the max pages limit
		if f.Scheduler.GetMaxPages() > 0 && f.Scheduler.GetCount() >= f.Scheduler.GetMaxPages() {
					// THE FIX: We must ALSO check that the queue is empty (len == 0)
					if len(f.Scheduler.queue) == 0 && f.active.Load() == 0 {
						slog.Info("Target page count reached and processed. Shutting down.")
						break
					}
				}

		select {
		case item, ok := <-f.Scheduler.queue:
			if !ok {
				goto finish
			}

			ctx := colly.NewContext()
			ctx.Put("depth", item.Depth)

			err := f.collector.Request("GET", item.URL, nil, ctx, nil)
			if err != nil {
				slog.Warn("Request failed", "url", item.URL, "err", err)
			}

		case <-time.After(1 * time.Second):
			// Only exit if the queue is truly empty AND zero requests are running
			if len(f.Scheduler.queue) == 0 && f.active.Load() == 0 {
				slog.Info("Queue empty and no active tasks. Shutting down...")
				goto finish
			}
		}
	}

finish:
	f.collector.Wait()
	slog.Info("Crawler engine finished", "total_pages", f.Scheduler.GetCount())
}
