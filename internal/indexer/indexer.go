package indexer

import "math"



type Indexer struct {
	Index map[string]map[int]int // token (word) -> docID -> count (freq)

	//cached matrices
	TF    map[string]map[int]float64
	IDF   map[string]float64
	TFIDF map[string]map[int]float64

	//metadata
	DocMaxFreq map[int]int // docID -> max frequency
	Documents map[int]string // docID -> path
	TotalDocs int

	lang  string
}

// NewIndex creates a new Index instance.
func NewIndexer(lang string) *Indexer {
	return &Indexer{
		Index: make(map[string]map[int]int),
		TF: make(map[string]map[int]float64),
		IDF: make(map[string]float64),
		TFIDF: make(map[string]map[int]float64),
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

func (i *Indexer) ComputeTF() {
	// Compute TF for each token in each document
	for word, docs := range i.Index {
		// compute TF for each doc that contains this token
		if _, ok := i.TF[word]; !ok {
			i.TF[word] = make(map[int]float64)
		}

		for docID, freq := range docs {
			// compute TF for this doc

			maxFreq := i.DocMaxFreq[docID]
			// store TF for this doc in the TF matrix
			i.TF[word][docID] =
				float64(freq) / float64(maxFreq)
		}
	}
}


func (i *Indexer) ComputeIDF() {
	// Compute IDF for each token
	for word, docs := range i.Index {
		// compute IDF for this token
		df := len(docs)
		// store IDF for this token in the IDF matrix
		i.IDF[word] =
			math.Log2(
				float64(i.TotalDocs) / float64(df),
			)
	}
}


func (i *Indexer) ComputeTFIDF() {
	// Compute TF-IDF for each token in each document
	for word, docs := range i.TF {
		// compute TF-IDF for each doc that contains this token
		if _, ok := i.TFIDF[word]; !ok {
			i.TFIDF[word] = make(map[int]float64)
		}
		// compute TF-IDF for each doc that contains this token
		for docID, tf := range docs {

			i.TFIDF[word][docID] =
				tf * i.IDF[word]
		}
	}
}

// Build the index by computing TF, IDF, and TF-IDF for each token in each document
func (i *Indexer) Build() {
	i.ComputeTF()
	i.ComputeIDF()
	i.ComputeTFIDF()
}
