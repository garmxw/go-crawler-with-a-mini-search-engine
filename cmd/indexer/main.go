package main

//needs to be removed and replaced by cmd/main.go this was just a test
import (
	"fmt" //we will import the indexer logic here
	"log"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/indexer"
)

func main () {

    idx := indexer.NewIndexer("english")//dynamique later

	err := indexer.LoadAndIndex("data/pages", idx)
	if err != nil {
			log.Fatal(err)
	}

	fmt.Println(idx.Index["go"])//just for now
}
