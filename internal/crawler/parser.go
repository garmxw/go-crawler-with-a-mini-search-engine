package crawler

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Page struct {
	URL string
	Title string
	Links []string
	Text string
}

func ParserPage(e *colly.HTMLElement) Page {
	page := Page{
		URL: e.Request.URL.String(),
	}

	// Title Extraction
	page.Title = e.ChildText("title")
	if page.Title == "" {
		page.Title = e.ChildText("h1")
	}

	// URL Normalization
	e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
		link := el.Attr("href")
		absolute := e.Request.AbsoluteURL(link)
		if absolute != "" && !strings.Contains(absolute, "javascript:") && !strings.HasPrefix(link, "#") {
			page.Links = append(page.Links, absolute)
		}
	})

	var contentSelection *goquery.Selection
	if main := e.DOM.Find("main"); main.Length() > 0 {
		contentSelection = main
    } else if article := e.DOM.Find("article"); article.Length() > 0 {
        contentSelection = article
    } else {
        contentSelection = e.DOM
    }

    // Clone the selection so we don't mess up the original links parsing
    domCopy := contentSelection.Clone()

    // REMOVE BLOAT: Delete elements that usually contain "junk" text (buttons, forms, navs, and common class names for sidebars/widgets)
    domCopy.Find("nav, footer, script, style, noscript, header, aside").Remove()
    domCopy.Find("button, form, .sidebar, .menu, .ads, .social-share").Remove()

    // Clean the whitespace
    rawText := domCopy.Text()
    words := strings.Fields(rawText)

    // OPTIONAL (for later): Filter out "Micro-Words"
    // For a search engine, words like "a", "the", "to" (Stop Words) are often bloat.
    // We'll keep it simple for now and just join the words.
    page.Text = strings.Join(words, " ")

	return page
}
