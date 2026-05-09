package crawler

import (
	"fmt"
	"net/url"
	"strings"
)

func ExtractDomain(rawURL string) string {
	//strip null bytes and trim spaces
	cleanURL := strings.Trim(rawURL, "\x00")
	cleanURL = strings.TrimSpace(cleanURL)

	parsed, err := url.Parse(cleanURL)
	if err != nil {
		fmt.Println("error in ExtractDomain: ", err)
		return ""
	}
	return parsed.Host
}

// NormalizeURL ensures we don't crawl "example.com/" and "example.com" as two pages
func NormalizeURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	u.Host = strings.ToLower(u.Host)
	u.Path = strings.TrimSuffix(u.Path, "/")
	u.Fragment = "" // Remove #anchors
	return u.String()
}
