package crawler

import (
	"net/url"
	"strings"
)

func ExtractDomain(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
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
