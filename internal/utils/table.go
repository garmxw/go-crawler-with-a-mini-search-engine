package utils

import (
	"fmt"
	"strings"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/indexer"
)


func PrintIndexDetails(idx *indexer.Indexer) {

	fmt.Println("\n/| /| /| INVERTED INDEX MATRIX -> :")

	// header
	fmt.Printf(
		"%-15s %-10s %-10s %-10s %-10s %-10s\n",
		"WORD",
		"DOC",
		"FREQ",
		"TF",
		"IDF",
		"TF-IDF",
	)

	fmt.Println(strings.Repeat("-", 75))

	for word, docs := range idx.Index {

		for docID, freq := range docs {

			tf := idx.TF[word][docID]
			idf := idx.IDF[word]
			tfidf := idx.TFIDF[word][docID]

			fmt.Printf(
				"%-15s %-10d %-10d %-10.4f %-10.4f %-10.4f\n",
				word,
				docID,
				freq,
				tf,
				idf,
				tfidf,
			)
		}
	}
}

func PrintDocStats(idx *indexer.Indexer) {

	fmt.Println("\n/| /| /| DOCUMENT STATS -> :")

	fmt.Printf(
		"%-10s %-10s %-20s\n",
		"DOC ID",
		"MAX FREQ",
		"PATH",
	)

	fmt.Println(strings.Repeat("-", 50))

	for docID, maxFreq := range idx.DocMaxFreq {

		fmt.Printf(
			"%-10d %-10d %-20s\n",
			docID,
			maxFreq,
			idx.Documents[docID],
		)
	}
}

// not in use right now needs fixing
/*
func PrintQueryVector(vector map[string]float64) {

	fmt.Println("\n/| /| /| QUERY VECTOR -> :")

	fmt.Printf(
		"%-15s %-15s\n",
		"WORD",
		"TF-IDF",
	)

	fmt.Println(strings.Repeat("-", 35))

	for word, value := range vector {

		fmt.Printf(
			"%-15s %-15.4f\n",
			word,
			value,
		)
	}
}
*/
