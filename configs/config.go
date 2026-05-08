package configs

import (
	"bufio"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port        string
	//optional DatabaseURL string
	CrawlDelay  string
}

func LoadConfig() *Config {
	// Load .env manually
	loadDotEnv(".env")

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		//DatabaseURL: getEnv("DATABASE_URL", ""),
		CrawlDelay:  getEnv("CRAWL_DELAY", "2"),
	}

	log.Println("Config loaded successfully")
	return cfg
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func loadDotEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(".env file not found, using system environment")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// split KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// set into environment
		os.Setenv(key, value)
	}
}

func ReadDelay() int {
	delayStr := os.Getenv("CrawlDelay")
	crawlDelay, err := strconv.Atoi(delayStr)
	if err != nil {
		crawlDelay = 2
		slog.Info("Info: CrawlerDelay is not a valid value")
	}
	return crawlDelay
}
