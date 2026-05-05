package indexer

import (
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`[a-zA-Z]+`)//we gonna ignore the numbers

func Tokenizer(text string) []string{

	text = strings.ToLower(text)
	// extract words only (ignore numbers, punctuation)
	matches := reg.FindAllString(text, -1)

	return matches
}
