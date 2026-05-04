package crawler

import (
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Fetcher struct {
	collector *colly.Collector
}

func NewFetcher(delay int64) *Fetcher {
	c := colly.NewCollector(
		colly.AllowedDomains("example.com"),//dynamique later
		colly.Async(true),
	)
	//wait X secs between requests (X is provided by the user)
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Parallelism: 2,
		Delay: time.Duration(delay) * time.Second,
	})
	return &Fetcher{
		collector: c,
	}
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
