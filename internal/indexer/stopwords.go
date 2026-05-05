package indexer

import (
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/stopwords_list"
)

func RemoveStopWords(tokens []string) []string {
	var result []string

	for _, token := range tokens {
		if !stopwords_list.StopWords[token] {
			result = append(result, token)
		}
	}

	return result
}
