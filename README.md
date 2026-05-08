1. The predected project structure :
   my-search-engine/
   │── go.mod
   │ # Defines the Go module (project name + dependencies)
   │
   │── go.sum
   │ # Stores checksums of dependencies (auto-generated for security/integrity)
   │
   ├── cmd/
   │ # Entry points (each folder = separate executable program)
   │
   │ ├── crawler/
   │ │ └── main.go
   │ │ # Starts the crawler service
   │ │ # Responsible for fetching web pages and sending them to the queue
   │ │
   │ ├── indexer/
   │ │ └── main.go
   │ │ # Starts the indexing worker
   │ │ # Consumes crawled pages and builds the search index
   │ │
   │ ├── api/
   │ │ └── main.go
   │ │ # Starts the HTTP API server
   │ │ # Exposes endpoints like /search?q=...
   │
   ├── internal/
   │ # Core application logic (private to this project, cannot be imported externally)
   │
   │ ├── crawler/
   │ │ # Handles everything related to crawling the web
   │ │
   │ │ ├── fetcher.go
   │ │ │ # Sends HTTP requests to fetch web pages
   │ │ │ # Handles timeouts, headers, retries
   │ │
   │ │ ├── parser.go
   │ │ │ # Parses HTML content
   │ │ │ # Extracts links, titles, and text from pages
   │ │
   │ │ ├── scheduler.go
   │ │ # Decides which URLs to crawl next
   │ │ # Prevents duplicates and manages crawl queue
   │
   │ ├── indexer/
   │ │ # Responsible for transforming raw data into a searchable format
   │ │
   │ │ ├── indexer.go
   │ │ │ # Builds the search index (e.g., inverted index)
   │ │ │ # Maps words → documents
   │ │
   │ │ ├── tokenizer.go
   │ │ # Splits text into tokens (words)
   │ │ # Handles normalization (lowercase, removing punctuation, etc.)
   │
   │ ├── search/
   │ │ # Handles search queries and ranking results
   │ │
   │ │ ├── engine.go
   │ │ │ # Main search logic
   │ │ │ # Takes query input and returns matching documents
   │ │
   │ │ ├── ranking.go
   │ │ # Ranks results based on relevance
   │ │ # Example: keyword frequency, scoring algorithms
   │
   │ ├── storage/
   │ │ # Data persistence layer (database or file system)
   │ │
   │ │ ├── db.go
   │ │ │ # Handles database connection and queries
   │ │ │ # Abstracts storage implementation (Postgres, files, etc.)
   │ │
   │ │ ├── models.go
   │ │ # Defines data structures (Page, Document, Index, etc.)
   │
   │ ├── queue/
   │ │ # Communication layer between components (crawler → indexer)
   │ │
   │ │ └── queue.go
   │ │ # Implements job queue
   │ │ # Could use Go channels (simple) or external systems like Redis/Kafka
   │
   ├── pkg/
   │ # Reusable utilities (can be used outside this project if needed)
   │
   │ ├── httpclient/
   │ │ # Custom HTTP client wrapper
   │ │ # Adds retries, headers, logging, etc.
   │ │
   │ ├── logger/
   │ # Logging utility
   │ # Standardized logs (info, error, debug)
   │
   ├── configs/
   │ # Application configuration
   │
   │ └── config.go
   │ # Loads environment variables and config settings
   │ # Example: DB URL, API port, crawl limits

2. the db structure :
   -- 1. Store the actual pages
   CREATE TABLE documents (
   id SERIAL PRIMARY KEY,
   url TEXT UNIQUE,
   title TEXT,
   body TEXT,
   indexed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   -- 2. Store every unique word (Term) found across the web
   CREATE TABLE terms (
   id SERIAL PRIMARY KEY,
   word TEXT UNIQUE,
   idf_score FLOAT8 DEFAULT 0 -- We update this after every crawl session
   );

   -- 3. The Inverted Index: Links words to documents with their Frequency (TF)
   CREATE TABLE inverted_index (
   term_id INTEGER REFERENCES terms(id),
   doc_id INTEGER REFERENCES documents(id),
   tf_score FLOAT8,
   PRIMARY KEY (term_id, doc_id)
   );

   -- Indexing for speed
   CREATE INDEX idx_terms_word ON terms(word);

the commands :

1. gocrawl --url https://a.com --url https://b.com : to read 1 or more urls
2. gocrawl --file urls.txt : to read a file of urls (can be txt/json for now i will add to it later)
3. gocrawl --json urls.json : to read a json file of urls in this structure :
   example json {"urls": ["https://example.com", "https://example.org"]}
4. and they all take the --depth flag example : gocrawl --url https://example.com --depth 2

the storage :
storage will be on data/storage/("ID".json)
the json will be like this :
{
"ID": 1,
"URL": "https://example.com",
"Title": "Example Domain",
"Date": "(2024-05-20T10:00:00Z)",
"Text": "This domain is for use in illustrative examples..."
}

launches :

1. to launch the indexer :
   go run ./cmd/indexer/main.go --mode=web --query="Finish each day and be done with it. You have done what you could. by Ralph Waldo Emerson" --detailed --limit=3
   NOTE : i ommited banch of flags the modes are web,local , live
   Global flags : --lang to specify the language, --detailed to get detailed, --query to specify the query, --limit to limit the number of results
   Local flags : --path to specify the local file path
   web flags : --storage to specify the storage directory
   live flags : reuse the crawler's flags

2. to launch the crawler :
   go run ./cmd/crawler/main.go --file || --json ./inputs/urls.txt --depth 2 --maxPages 5
   -----or with a single url :
   go run ./cmd/crawler/main.go --url https://example.com --depth 2 --maxPages 5
   NOTE : they all can take the optional --storage flag to specify the storage directory and default to data/storage/
