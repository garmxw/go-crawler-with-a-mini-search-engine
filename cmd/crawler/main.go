package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/crawler"
)

func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//check the extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".json"{
		var urls []string
		if err := json.NewDecoder(file).Decode(&urls); err != nil {
			return nil, err
		}
		return urls, nil
	}
	if ext == ".txt" {
		var urls []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return urls, nil
	}
	return nil, nil
}



func main () {
	delayStr := os.Getenv("CrawlDelay")
	crawlDelay,err := strconv.Atoi(delayStr)
	if err != nil {
		crawlDelay = 2
		slog.Error("Error: CrawlerDelay is not a valid value")
	}
	//flags
	urlsFlag := flag.String("urls", "", "URLs to crawl")
	fileFlag := flag.String("file", "", "File containing URLs to crawl (txt || json)")
	depthFlag := flag.Int("depth", 3, "Depth of the crawl")

	flag.Parse()

	if *urlsFlag == "" && *fileFlag == "" {
		slog.Error("Error: No URLs provided")
		return
	}
	if *depthFlag < 1 {
		slog.Warn("Warning: Depth is less than 1, setting to default depth of 3")
		*depthFlag = 3
	}

	// create a schedular
	schedular := crawler.NewScheduler(100)

	urls := []string{}
	if *urlsFlag != "" {
		urls = append(urls, strings.Split(*urlsFlag, ",")...)
	}
	if *fileFlag != "" {
		fileUrls, err := readURLsFromFile(*fileFlag)
		if err != nil {
			slog.Error("Error: Failed to read URLs from file")
			os.Exit(1)
		}
		urls = append(urls, fileUrls...)
	}
	//validation
	if len(urls) == 0 {
		slog.Error("Error: No URLs provided")
		return
	}
	// add a URL to the schedular
	for _,url := range urls {
		schedular.Add(url)
	}
	// create a fetcher
	f := crawler.NewFetcher(schedular, int64(crawlDelay))
	// start crawling
	slog.Info("Crawling started ...")
	slog.Info("Crawling depth:", "depth" ,*depthFlag)
	f.Fetch("https://example.com")
}
