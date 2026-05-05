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

// custom flag type for multi-value flags
type multiFlag []string

func (m *multiFlag) String() string {
	return ""
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func ReadJSON(filename string) ([]string, error) {
	file, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		var data struct {
			URLs []string `json:"urls"`
		}

		err = json.Unmarshal(file, &data)
		return data.URLs, err
}


func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//check the extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".json"{
		return ReadJSON(filename)
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
	var urls multiFlag
	flag.Var(&urls, "url", "Add URL (can be used multiple times)")
	fileFlag := flag.String("file", "", "File containing URLs to crawl")
	// example json {"urls": ["https://example.com", "https://example.org"]}
	jsonFlag := flag.String("json", "", "JSON file with URLs")
	depthFlag := flag.Int("depth", 3, "Depth of the crawl")
	maxPageFlag := flag.Int("maxPages", 0, "Maximum number of pages to crawl")

	flag.Parse()

	if len(urls) == 0 && *fileFlag == "" {
		slog.Error("Error: No URLs provided")
		return
	}
	if *depthFlag < 1 {
		slog.Warn("Warning: Depth is less than 1, setting to default depth of 3")
		*depthFlag = 3
	}

	if *jsonFlag != "" {
		jsonUrls, err := ReadJSON(*jsonFlag)
		if err != nil {
			slog.Error("Error: Failed to read URLs from JSON file")
			os.Exit(1)
		}
		urls = append(urls, jsonUrls...)
	}

	if *maxPageFlag == 0 {
		*maxPageFlag = -1 //negative value indicates no limit
	}
	// create a schedular
	schedular := crawler.NewScheduler(100, *depthFlag, *maxPageFlag)

	if len(urls) == 0 {
		if *fileFlag != "" {
			fileUrls, err := readURLsFromFile(*fileFlag)
			if err != nil {
				slog.Error("Error: Failed to read URLs from file")
				os.Exit(1)
			}
			urls = append(urls, fileUrls...)
		}
	}
	//validation
	if len(urls) == 0 {
		slog.Error("Error: No URLs provided")
		return
	}
	// add a URL to the schedular
	for _,url := range urls {
		schedular.Add(url, 0)
	}
	// extract domains
	domainsMap := make(map[string]bool)
	for _,url := range urls{
		d := crawler.ExtractDomain(url)
		if d != "" {
			domainsMap[d] = true
		}
	}
	var domains []string
	for d := range domainsMap {
		domains = append(domains, d)
	}

	// create a fetcher
	f := crawler.NewFetcher(schedular,domains, int64(crawlDelay))
	// start crawling
	slog.Info("Crawling started ...")
	slog.Info("Crawling depth:", "depth" ,*depthFlag)
	f.Start()
}
