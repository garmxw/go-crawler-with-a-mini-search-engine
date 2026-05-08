package crawler

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/garmxw/go-crawler-with-a-mini-search-engine/internal/storage"
)

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

func RunCrawler(
	urls []string,
	depth int,
	maxPages int,
	delay int,
	filePath string,
	jsonFilePath string,
	storagePath string,
) error {


	if len(urls) == 0 && filePath == "" && jsonFilePath == "" {
		slog.Error("Error: No URLs provided and no file specified")
		os.Exit(1)
	}
	if depth < 0 {
		slog.Warn("Warning: Depth is less than 0, setting to default depth of 0")
		depth = 0
	}

	if jsonFilePath != "" {
		jsonUrls, err := ReadJSON(jsonFilePath)
		if err != nil {
			slog.Error("Error: Failed to read URLs from JSON file")
			os.Exit(1)
		}
		urls = append(urls, jsonUrls...)
	}

	if maxPages == 0 {
		maxPages = -1 //negative value indicates no limit
	}

	// create a schedular
	schedular := NewScheduler(100, depth, maxPages)

	if len(urls) == 0 {
		if filePath != "" {
			fileUrls, err := readURLsFromFile(filePath)
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
		os.Exit(1)
	}
	// add a URL to the schedular
	for _,url := range urls {
		schedular.Add(url, 0)
	}
	// extract domains
	domainsMap := make(map[string]bool)
	for _,url := range urls{
		d := ExtractDomain(url)
		if d != "" {
			domainsMap[d] = true
		}
	}
	var domains []string
	for d := range domainsMap {
		domains = append(domains, d)
	}

	// create a storage
	st := storage.NewStorage(storagePath) //defaults to "data/pages"

	// create a fetcher
	f := NewFetcher(schedular, st, domains, int64(delay), depth)
	// start crawling
	slog.Info("Crawling started ...")
	slog.Info("Crawling depth:", "depth" ,depth)
	f.Start()

	return nil
}
