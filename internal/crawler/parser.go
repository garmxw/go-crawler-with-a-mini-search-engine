package crawler

import (
	"strings"

	"github.com/gocolly/colly"
)

type Page struct {
	URL string
	Title string
	Links []string
	Text string
}

func ParserPage (e *colly.HTMLElement) Page {
	page := Page{
		URL: e.Request.URL.String(),
	}
	//title
	page.Title = e.ChildText("title")
	//links
	e.ForEach("a[href]", func(_ int, el *colly.HTMLElement){
		link := el.Attr("href")
		if strings.HasPrefix(link, "http") {
			page.Links = append(page.Links, link)
		}
	})
	//text (simple version)
	page.Text = strings.TrimSpace(e.Text)

	return page
}
