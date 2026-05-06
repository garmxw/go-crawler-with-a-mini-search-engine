package indexer

import "math"



type Indexer struct {
	Index map[string]map[int]int // token (word) -> docID -> count (freq)
	DocMaxFreq map[int]int // docID -> max frequency
	Documents map[int]string // docID -> path
	TotalDocs int
	lang  string
}

// NewIndex creates a new Index instance.
func NewIndexer(lang string) *Indexer {
	return &Indexer{
		Index: make(map[string]map[int]int),
		DocMaxFreq: make(map[int]int),
		Documents: make(map[int]string),
		TotalDocs: 0,
		lang:  lang,
	}
}

func (i *Indexer) Add(docID int, text string, path string) {
	// store document path for later retrieval
	i.Documents[docID] = path
	// tokenize, remove stopwords, and stem the text
	tokens := Tokenizer(text)
	tokens = RemoveStopWords(tokens)
	tokens = StemWords(tokens, i.lang)
    // count token frequency in local doc
	localFreq := make(map[string]int)
	for _, token := range tokens{
		localFreq[token]++
	}
    // add token to index and update doc max freq
	for token, freq := range localFreq {
		if _, ok := i.Index[token]; !ok {
			i.Index[token] = make(map[int]int)
		}
		// add token to index for this doc
		i.Index[token][docID] = freq
		// update doc max freq  in the index struct
		if freq > i.DocMaxFreq[docID] {
			i.DocMaxFreq[docID] = freq
		}
	}
	i.TotalDocs++
}

func (i *Indexer) ComputeTF() map[string]map[int]float64 {
	// Compute TF for each token in each document
	tf := make(map[string]map[int]float64)

	for word, docs := range i.Index {
		// compute TF for each doc that contains this token
		tf[word] = make(map[int]float64)
		for docID, freq := range docs {
			maxFreq := i.DocMaxFreq[docID]
			tf[word][docID] = float64(freq) / float64(maxFreq)
		}
  	}
  	return tf
}


func (i *Indexer) ComputeIDF() map[string]float64 {
	idf := make(map[string]float64)

	for word, docs := range i.Index {
		df := len(docs)
		//compute IDF for this token
		idf[word] = math.Log2(float64(i.TotalDocs) /float64(df))
	}
	return idf
}
