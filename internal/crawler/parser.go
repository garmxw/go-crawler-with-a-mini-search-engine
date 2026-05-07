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

func ParserPage(e *colly.HTMLElement) Page {
	page := Page{
		URL: e.Request.URL.String(),
	}

	// Title Extraction
	// Fallback to <h1> if <title> tag is missing or empty
	page.Title = e.ChildText("title")
	if page.Title == "" {
		page.Title = e.ChildText("h1")
	}

	// URL Normalization
	e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
		link := el.Attr("href")

		// Use AbsoluteURL to handle relative paths like "/about" automatically
		absolute := e.Request.AbsoluteURL(link)

		// Filter out empty strings, javascript: void(0), and anchor fragments
		if absolute != "" &&
		   !strings.Contains(absolute, "javascript:") &&
		   !strings.HasPrefix(link, "#") {
			page.Links = append(page.Links, absolute)
		}
	})

	// We use a clone of the selection to remove "junk" tags before grabbing text.
	// This prevents navbars, scripts, and footers from cluttering your data.
	domCopy := e.DOM.Clone()
	domCopy.Find("nav, footer, script, style, noscript, .sidebar, .menu, #footer, #header").Remove()

	// Get the text from the remaining content (usually <main> or <body> minus the junk)
	page.Text = strings.TrimSpace(domCopy.Text())

	return page
}
