package crawler

import (
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Fetcher struct {
	collector *colly.Collector
	Scheduler *Scheduler
}

func NewFetcher(s *Scheduler, delay int64) *Fetcher {
	c := colly.NewCollector(
		colly.AllowedDomains("example.com"),//dynamique later
		colly.Async(true),
	)
	//wait X secs between requests (X is provided by the user)
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Parallelism: 2,
		Delay: time.Duration(delay) * time.Second,
		RandomDelay: time.Duration(delay - 1) * time.Second,
	})
	 f := &Fetcher{
		collector: c,
		Scheduler: s,
	}
	f.registerCallbacks()

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
		page := ParserPage(e)
		log.Println("Parsed page: ", page)
		log.Println("Title: ", page.Title)
		//add new links to schedular
		for _, link := range page.Links {
			f.Scheduler.Add(link)
		}
	})

	f.collector.OnError(func(r *colly.Response, err error) {
		log.Println("Error: ", err.Error())
	})

}

func (f *Fetcher) Start() {
	for {
		url := f.Scheduler.Next()
		f.collector.Visit(url)
	}
}
