package search

import (
	"math"
)



func BuildDocumentVector(
	docID int,
	tfidf map[string]map[int]float64,
) map[string]float64 {

	vector := make(map[string]float64)

	for word, docs := range tfidf {

		if score, ok := docs[docID]; ok {
			vector[word] = score
		}
	}

	return vector
}


func cosineSimilarity(a, b map[string]float64) float64 {
	var dotProduct float64
	var normA, normB float64

	//dot product + normalization a
	for key, valA := range a {
		valB := b[key]
		dotProduct += valA * valB
		normA += valA * valA
	}
	//normalization b
	for _, valB := range b {
		normB += valB * valB
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
