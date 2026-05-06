package local

import (
	"os"
	"path/filepath"
	"strings"
)

type Document struct {
   ID int
   Path string
   Text string
}

type Loader struct {
   Path string
}

func NewLoader(path string) *Loader {
	return &Loader{
		Path: path,
	}
}

func (l *Loader) Load()([]Document ,error) {
	var docs []Document
	id := 1
	err := filepath.Walk(l.Path, func(path string, info os.FileInfo, err error) error{
		if err != nil {
			return err
		}
		//skip folders
		if info.IsDir() {
			return nil
		}
		// read file
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		//only text files
		if !isTextFile(path) {
			return nil
		}

		docs = append(docs, Document{
			ID:   id,
			Path: path,
			Text: string(content),
		})
		id++
		return nil
	})
	return docs, err
}


//helpers
func isTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	// allowed extensions
	case ".txt", ".md", ".html", ".log":
		return true
	default:
		return false
	}
}
