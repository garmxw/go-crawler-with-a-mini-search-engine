package indexer

import (
	"github.com/kljensen/snowball"
)

// Stem returns the root of a word using the Snowball algorithm.
// It supports multiple languages (e.g., "english", "french").
func Stem(word string, lang string ) string {
	//handle the default language
	if lang == "" {
		lang = "english"
	}
	// Stem the word using the snowball stemmer
	stemmed, err := snowball.Stem(word, lang, true)
	if err != nil {
		return word
	}
	return stemmed
}

func StemWords(tokens []string, lang string) []string {
	stemmed := make([]string, 0, len(tokens))
	for _, token := range tokens {
		stemmed = append(stemmed, Stem(token, lang))
	}
	return stemmed
}
