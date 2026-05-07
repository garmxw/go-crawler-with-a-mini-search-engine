package search

import (
	"fmt"
	"sort"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/indexer"
	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/models"
)



func ProcessQuery(query string, lang string) []string {
	tokens := indexer.Tokenizer(query)
	tokens = indexer.RemoveStopWords(tokens)
	tokens = indexer.StemWords(tokens, lang)
	return tokens
}


func BuildQueryVector(
	tokens []string,
	idf map[string]float64,
) map[string]float64 {

	vector := make(map[string]float64)

	freq := make(map[string]int)
	maxFreq := 0

	// count query frequencies
	for _, token := range tokens {

		freq[token]++

		if freq[token] > maxFreq {
			maxFreq = freq[token]
		}
	}

	// build TF-IDF query vector
	for token, count := range freq {

		tf :=
			float64(count) /
				float64(maxFreq)

		vector[token] =
			tf * idf[token]
	}

	return vector
}

func Search(query string, idx *indexer.Indexer, lang string) []models.SearchResult {
	tokens := ProcessQuery(query, lang)
	fmt.Println("\nProcessed query:", tokens)
	idf := idx.IDF
	queryVector := BuildQueryVector(tokens, idf)
	var results []models.SearchResult

	for docID, path := range idx.Documents {
		docVector := BuildDocumentVector(
			docID,
			idx.TFIDF,
		)
		score := cosineSimilarity(queryVector, docVector)
		if score > 0 {
			results = append(results, models.SearchResult{
				DocID: docID,
				Path:  path,
				Score: score,
			})
		}
	}
	sort.Slice(results, func(i, j int) bool{
		return results[i].Score > results[j].Score
	})
	return results
}
