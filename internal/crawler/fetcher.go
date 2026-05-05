package crawler

import (
	"log"
	"log/slog"
	"time"

	"github.com/gocolly/colly"
)

type Fetcher struct {
	collector *colly.Collector
	Scheduler *Scheduler
}

func NewFetcher(s *Scheduler, domains []string, delay int64) *Fetcher {
	c := colly.NewCollector(
		colly.AllowedDomains(domains...),
		colly.Async(true),
	)
	var randomDelay time.Duration
	if delay > 0 {
		randomDelay = time.Duration(delay / 2) * time.Second
	}

	//wait delay between requests (delay is provided by the user)
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Parallelism: 2,
		Delay: time.Duration(delay) * time.Second,
		RandomDelay: randomDelay,
	})

	// Identity: Set the User-Agent right here
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (compatible; GoCrawler/1.0)")
	})

	f := &Fetcher{
		collector: c,
		Scheduler: s,
	}
	f.registerCallbacks()
	// Log everything in one or two structured lines
    slog.Info("Fetcher initialized",
        "domains", domains,
        "delay", time.Duration(delay)*time.Second,
        "randomDelay", randomDelay,
    )

    // Logging the queue specifically
    slog.Info("Scheduler attached",
        "queue_capacity", cap(f.Scheduler.queue),
    )

	return f
}

func (f *Fetcher) Fetch(url string) {
	f.collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	f.collector.OnResponse(func(r *colly.Response){
		log.Println("Got response from:", r.Request.URL)
	})

	f.collector.OnHTML("a[href]", func(r *colly.HTMLElement) {
		link := r.Attr("href")
		log.Println("Found link:", link)
	})
	f.collector.Visit(url)
}

func (f *Fetcher) registerCallbacks() {
	//when visiting a page
	f.collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting: ", r.URL.String())
	})
	//when the html is recieved
	f.collector.OnHTML("html", func(e *colly.HTMLElement) {
		depth := 0
		if d:= e.Request.Ctx.GetAny("depth"); d !=nil {
			depth = d.(int)
		}
		page := ParserPage(e)
		log.Println("Parsed page: ", page)
		log.Println("Title: ", page.Title)
		//add new links to schedular
		for _, link := range page.Links {
			f.Scheduler.Add(link, depth+1)
		}
	})

	f.collector.OnError(func(r *colly.Response, err error) {
		log.Println("Error: ", err.Error())
	})

}

func (f *Fetcher) Start() {
	for {
		item := f.Scheduler.Next()
		f.collector.Request("GET", item.URL, nil, colly.NewContext(),nil)
		ctx := colly.NewContext()
		ctx.Put("depth", item.Depth)
	}
}
