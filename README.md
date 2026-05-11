# Go Crawler + Mini Search Engine

A modular web crawler and mini search engine written in Go.

This project combines:

- Recursive web crawling
- Document storage
- Tokenization
- Stop-word removal
- Stemming
- Inverted indexing
- TF / IDF / TF-IDF calculations
- Cosine similarity ranking
- Local file search
- Web page search
- Live crawl + instant search

---

# Branches

This project ships in three branches. Pick the one that matches your environment.

| Branch         | Interface                | Windows Terminal Preview | Editor Terminals | cmd.exe |
| -------------- | ------------------------ | ------------------------ | ---------------- | ------- |
| `main`         | Prompt-based TUI         | yes                      | yes              | yes     |
| `tui-advanced` | Full keyboard-driven TUI | no                       | yes              | yes     |
| `cli-flags`    | Classic CLI flags        | yes                      | yes              | yes     |

---

## `main` — Prompt TUI (recommended, works everywhere)

Interactive step-by-step prompts. No raw terminal mode.
Styled with colours and a confirm panel. Works on every terminal including
Windows Terminal Preview, cmd.exe, VSCode, Zed, and all Unix terminals.

```bash
git checkout main
go run ./cmd/search
go run ./cmd/crawler
```

No flags needed. The tool asks you everything interactively.

---

## `tui-advanced` — Full keyboard-driven TUI

A full keyboard-navigable form built with Bubbletea.
Tab between fields, arrow keys to change values, Enter to confirm.
Looks closest to tools like k9s and lazygit.

> **Not compatible with Windows Terminal Preview.**
> Use this branch on VSCode, Zed, Warp, iTerm2, or any Unix terminal.

```bash
git checkout tui-advanced
go run ./cmd/search
go run ./cmd/crawler
```

---

## `cli-flags` — Classic CLI flags

Original flag-based interface. No interactivity.
Fully scriptable and pipe-friendly. Works on every terminal and in CI.

```bash
git checkout cli-flags
go run ./cmd/crawler/main.go --url https://example.com --depth 2
go run ./cmd/search/main.go --mode=local --path=./docs --query="golang concurrency"
```

---

# Features

## Web Crawler

- Recursive crawling
- Depth control
- Maximum page limits
- Multi-URL crawling
- Async crawling using Colly
- URL deduplication
- Domain restriction
- JSON page storage

---

## Search Engine

- Tokenization
- Stop-word removal
- Stemming
- Inverted index
- TF matrix
- IDF vector
- TF-IDF matrix
- Cosine similarity ranking
- Query processing
- Ranked search results

---

# Architecture Overview

The project is separated into multiple independent systems.

---

## cmd/

Contains executable entry points.

### `cmd/crawler/main.go`

Runs the crawler only.

Responsible for:

- reading input (flags or TUI depending on branch)
- reading URLs
- starting crawler pipeline

It does NOT:

- index documents
- rank results
- search queries

---

### `cmd/search/main.go`

Runs the search engine.

Responsible for:

- local search
- web search
- live crawl + search

---

## internal/ui/

Contains all terminal UI and styling code.
Completely isolated from business logic.

### `banner.go`

Prints the ASCII art banner at startup for both tools.

### `messages.go`

Styled output helpers: `Success`, `Warn`, `Error`, `Info`, config panels, mode badges.

### `results.go`

Renders the ranked search results list with score bars.

### `spinner.go`

Animated spinner shown while crawling or indexing.

### `styles.go`

Colour palette and shared lipgloss style variables.

---

## internal/ui/tui/

Interactive input layer. Implementation varies by branch.

### `main` branch — `prompt.go`

Step-by-step prompt form using `bufio.Scanner`.
No raw terminal mode. Compatible with all terminals.

### `tui-advanced` branch — `form_search.go`, `form_crawler.go`, `components.go`, `theme.go`

Full Bubbletea TUI with keyboard navigation, tab focus, steppers, toggles.

---

## internal/crawler/

Contains the crawling system.

### `engine.go`

Reusable crawler runner.

Responsible for:

- creating scheduler
- creating storage
- creating fetcher
- starting crawler pipeline

---

### `fetcher.go`

Handles HTTP requests using Colly.

Responsible for:

- visiting pages
- handling callbacks
- extracting responses
- sending parsed pages to storage

---

### `parser.go`

Extracts data from HTML pages.

Responsible for:

- extracting titles
- extracting text
- extracting links

---

### `scheduler.go`

Controls crawl flow.

Responsible for:

- queueing URLs
- tracking visited URLs
- depth control
- page limits

---

### `utils.go`

Contains crawler helper functions:

- URL normalization
- domain extraction

---

## internal/indexer/

Contains the indexing pipeline.

### `indexer.go`

Core indexing logic.

Responsible for:

- inverted index
- TF calculation
- IDF calculation
- TF-IDF calculation

---

### `tokenizer.go`

Splits text into tokens.

---

### `stopwords.go`

Removes common words such as `the`, `and`, `is`.

---

### `stemmer.go`

Reduces words to their roots.

Examples:

- running → run
- playing → play

---

### `Loader.go`

Loads and indexes pages from disk (used by the indexer package).

---

## internal/local/

### `loader.go`

Loads local documents from the filesystem for local search mode.

---

## internal/search/

Contains search and ranking logic.

### `engine.go`

High-level search modes: local, web, live.

### `query.go`

Processes user queries: tokenization, stop-word removal, stemming, query vector creation.

### `cosine.go`

Calculates cosine similarity between query and document vectors.

### `ranking.go`

Sorts results from highest to lowest score.

---

## internal/storage/

### `storage.go`

Stores crawled pages as JSON files and loads them back.

---

## internal/models/

### `result.go`

Defines the `SearchResult` struct shared across packages.

---

## internal/utils/

Debug and print helpers:

- printing TF matrix
- printing TF-IDF matrix
- printing document stats

---

## internal/cli/

### `main.go`

`MultiFlag` type — allows repeated flags like `--url a --url b`.
Used only in the `cli-flags` branch.

---

# Storage Format

Crawled pages are stored inside `data/pages/`.

Each page is a JSON file:

```json
{
  "ID": 1,
  "URL": "https://example.com",
  "Title": "Example Domain",
  "Date": "2024-05-20T10:00:00Z",
  "Text": "This domain is for use in illustrative examples..."
}
```

---

# Workflows

## Crawler Workflow

```
URL(s)
 ↓
Scheduler
 ↓
Fetcher
 ↓
Parser
 ↓
Storage
```

The crawler visits pages, extracts text and links,
recursively crawls discovered links, and stores pages as JSON.

---

## Indexing Workflow

```
Documents
 ↓
Tokenizer
 ↓
Stop-word Removal
 ↓
Stemmer
 ↓
Inverted Index
 ↓
TF → IDF → TF-IDF
```

---

## Search Workflow

```
Query
 ↓
Query Processing
 ↓
Query Vector
 ↓
Cosine Similarity
 ↓
Ranking
 ↓
Results
```

---

# Usage

## `main` branch — Prompt TUI

Just run the binary. The tool guides you through all options interactively.

```bash
# Crawler
go run ./cmd/crawler

# Search engine
go run ./cmd/search
```

Example session for the search engine:

```
  Mode
  * 1. local
    2. web
    3. live
  choice > [1]  3

  Query > ralph waldo emerson life

  Language
  * 1. english
    2. french
  choice > [1]

  Limit > [5]  3
  Detailed output (y/n) > [n]

  URLs  (one per line, empty line to finish)
  url 1 > https://quotes.toscrape.com/tag/books/
    added: https://quotes.toscrape.com/tag/books/
  url 2 >

  Depth > [0]  2
  Max pages > [3]  10
  Delay (s) > [2]
  Storage path > [data/pages]

  ┌─────────────────────────────────────────┐
  │ [*] Ready to run                        │
  │     Mode  |  LIVE                       │
  │    Query  |  "ralph waldo emerson life" │
  │    Depth  |  2                          │
  └─────────────────────────────────────────┘

  Run with this config? (y/n) > [y]
```

---

## `tui-advanced` branch — Keyboard TUI

Run the binary. A full-screen form appears. Navigate with keyboard.

```bash
go run ./cmd/crawler
go run ./cmd/search
```

Key bindings:

| Key               | Action                               |
| ----------------- | ------------------------------------ |
| `tab` / `↓`       | Next field                           |
| `shift+tab` / `↑` | Previous field                       |
| `← →`             | Cycle options / change stepper value |
| `space`           | Toggle boolean field                 |
| `enter`           | Confirm field / add URL / submit     |
| `backspace`       | Delete character / remove last URL   |
| `ctrl+c`          | Quit                                 |

---

## `cli-flags` branch — Classic flags

### Crawler

```bash
# Single URL
go run ./cmd/crawler/main.go --url https://example.com

# Multiple URLs
go run ./cmd/crawler/main.go \
  --url https://example.com \
  --url https://go.dev

# From a TXT file
go run ./cmd/crawler/main.go --file ./urls.txt

# From a JSON file
go run ./cmd/crawler/main.go --json ./urls.json
```

TXT file format:

```
https://example.com
https://go.dev
```

JSON file format:

```json
{ "urls": ["https://example.com", "https://go.dev"] }
```

Crawler flags:

| Flag         | Required | Default      | Description                       |
| ------------ | -------- | ------------ | --------------------------------- |
| `--url`      | optional | none         | Add one or more URLs (repeatable) |
| `--file`     | optional | none         | TXT file of URLs                  |
| `--json`     | optional | none         | JSON file of URLs                 |
| `--depth`    | optional | `0`          | Maximum crawl depth               |
| `--maxPages` | optional | `3`          | Maximum pages to crawl            |
| `--storage`  | optional | `data/pages` | Where to save crawled pages       |

---

### Search engine

```bash
# Local mode
go run ./cmd/search/main.go \
  --mode=local \
  --path="./docs" \
  --query="golang concurrency"

# Web mode (search previously crawled pages)
go run ./cmd/search/main.go \
  --mode=web \
  --query="react hooks"

# Live mode (crawl + index + search in one step)
go run ./cmd/search/main.go \
  --mode=live \
  --url https://quotes.toscrape.com/tag/books/ \
  --query="life happiness" \
  --depth 2 \
  --maxPages 10
```

Search flags:

| Flag         | Required        | Default      | Description                      |
| ------------ | --------------- | ------------ | -------------------------------- |
| `--mode`     | yes             | none         | `local`, `web`, or `live`        |
| `--query`    | yes             | none         | Search query                     |
| `--path`     | local mode only | none         | Path to local documents          |
| `--storage`  | optional        | `data/pages` | Path to stored pages             |
| `--lang`     | optional        | `english`    | Language (`english` or `french`) |
| `--limit`    | optional        | `1`          | Max results to return            |
| `--detailed` | optional        | `false`      | Print TF/IDF debug info          |
| `--url`      | live mode       | none         | URL(s) to crawl (repeatable)     |
| `--file`     | live mode       | none         | TXT file of URLs                 |
| `--json`     | live mode       | none         | JSON file of URLs                |
| `--depth`    | live mode       | `0`          | Crawl depth                      |
| `--maxPages` | live mode       | `3`          | Max pages                        |

---

# Ranking System

Documents are ranked using TF-IDF and cosine similarity.

**TF**

```
TF = frequency(word in document) / max frequency in document
```

**IDF**

```
IDF = log2(total documents / document frequency)
```

**TF-IDF**

```
TF-IDF = TF × IDF
```

**Cosine Similarity** compares the query vector against each document vector.
Documents with higher similarity scores rank higher.

---

# Detailed Mode

Passing `--detailed` (flags branch) or enabling it in the TUI prints:

- Inverted index
- TF matrix
- IDF vector
- TF-IDF matrix
- Document statistics

---

# Dependencies

| Package                              | Purpose                                   |
| ------------------------------------ | ----------------------------------------- |
| `github.com/gocolly/colly`           | Web crawling                              |
| `github.com/PuerkitoBio/goquery`     | HTML parsing                              |
| `github.com/charmbracelet/lipgloss`  | Terminal styling                          |
| `github.com/pterm/pterm`             | Spinner and tables                        |
| `github.com/charmbracelet/bubbletea` | Keyboard TUI (`tui-advanced` branch only) |
| `github.com/kljensen/snowball`       | Word stemming                             |

---

# Setting Up Each Branch

```bash
# Clone the repo
git clone https://github.com/garmxw/go-crawler-with-a-mini-search-engine
cd go-crawler-with-a-mini-search-engine

# main branch — prompt TUI, works everywhere
git checkout main
go mod tidy
go run ./cmd/search

# tui-advanced — keyboard TUI, not for Windows Terminal Preview
git checkout tui-advanced
go mod tidy
go run ./cmd/search

# cli-flags — classic flags, fully scriptable
git checkout cli-flags
go mod tidy
go run ./cmd/search/main.go --mode=local --path=./docs --query="your query"
```

---

# Current Limitations

- No PageRank
- No fuzzy matching
- No distributed crawling
- No persistent index database
- No phrase search
- No JavaScript rendering
- No proxy rotation

---

# Future Improvements

- BM25 ranking
- PostgreSQL index storage
- Redis queues
- REST API
- Web interface
- Proxy rotation
- Headless browser support
- Distributed crawling
- Incremental indexing
- Snippet generation
- Query autocomplete

---

# Technologies Used

- Go
- Colly
- JSON storage
- TF-IDF
- Cosine Similarity
- Concurrent crawling
- Inverted indexes

---

# Educational Concepts Covered

- Information Retrieval (IR)
- Search Engine Architecture
- Vector Space Model
- Recursive and concurrent crawling
- Inverted indexing
- TF-IDF ranking
- Cosine similarity
- Query processing
- Tokenization and stemming
- Modular Go architecture
