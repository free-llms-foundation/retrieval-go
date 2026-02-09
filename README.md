<p align="center">
  <img src="./logo.png" width="240" alt="retrieval-go logo" />
</p>

# retrieval-go

A fast, batteries-included web retrieval library for Go.

`retrieval-go` is a production-oriented DuckDuckGo HTML search client + web page content extractor designed for building LLM tools (RAG, agents, “browse the web” features) without relying on paid search APIs.

It provides:

- Search (`query -> []Page`) using DuckDuckGo’s Lite HTML endpoint.
- Content extraction (`url -> Document`) powered by [`go-readability`](https://github.com/go-shiori/go-readability).
- Robust defaults (headers, user-agent selection, decompression, timeouts).
- A configurable `Config` API plus ergonomic `Option` helpers.

---

## Why

Most LLM applications need three things from “internet access”:

- **Discover** relevant URLs for a query.
- **Fetch & extract** clean, readable text.
- **Keep it reliable** (timeouts, gzip/br/zstd, sane defaults).

`retrieval-go` is built specifically for that workflow.

---

## Features

- **DuckDuckGo Lite HTML search** (works without browser automation)
- **Clean text extraction**:
  - Returns **Markdown** (including tables!)
  - Filters clutter (ads, nav bars)
  - **Basic image filtering** (removes common icons, logos, flags, and tiny elements)
- **Automatic decompression** (`gzip`, `br`, `zstd`, `deflate`)
- **Production defaults**
  - request headers close to a real browser
  - error-body capture with size limit
- **Extensible parsing** via `Parser` interface
- **Two configuration styles**
  - `New(opts...)` (simple)
  - `NewWithConfig(cfg)` (explicit)

---

## Install

```bash
go get github.com/free-llms-foundation/retrieval-go
```

---

## Quick start

### 1) Search

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/free-llms-foundation/retrieval-go"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	c, err := retrieval.New()
	if err != nil {
		log.Fatal(err)
	}

	// Second argument is the search query
	// Third argument is the date filter (e.g., "d" for day, "w" for week, "m" for month, "y" for year, or "" for all time)
	pages, err := c.SearchWithQuery(ctx, "golang generics site:golang.org", "")
	if err != nil {
		log.Fatal(err)
	}

	for i, p := range pages {
		if i >= 5 {
			break
		}

		fmt.Println("Link:", p.Link)
		fmt.Println("Title:", p.Title)
		fmt.Println("Snippet:", p.Snippet)
		fmt.Println("Favicon:", p.Favicon)
		fmt.Print("-------------------------------------------------\n\n")
	}
}

```

### 2) Fetch + extract readable content

```go
func main() {
	c, _ := retrieval.New()
	ctx := context.Background()

	// Third argument is robotsTxtAllowed (bool)
	doc, err := c.ParseContentFromLink(ctx, "https://go.dev/blog/generics-next-step", true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Title: %s\n", doc.Title)
	fmt.Printf("Site: %s (%s)\n", doc.SiteName, doc.Language)
	fmt.Println("---")
	fmt.Println(doc.Content) // Clean Markdown
	
	if len(doc.Images) > 0 {
		fmt.Printf("\nExtracted %d images (filtered)\n", len(doc.Images))
	}
}
```

---

## Examples

See the `examples/` folder:

- `examples/search` — search only
- `examples/fetch` — fetch + readability extraction for a single URL
- `examples/research` — search → fetch multiple pages → print extracted text

---

## Data model

### `Page`

```go
type Page struct {
	Link    string
	Title   string
	Source  string  // Hostname or source name
	Snippet string  // Brief summary from search results
	Favicon string  // URL to the site's favicon
}
```

### `Document`

```go
type Document struct {
	Title       string
	Byline      string
	Content     string   // Main content in Markdown
	TextContent string   // Raw plain text
	Excerpt     string   // Short summary/teaser
	SiteName    string 
	Image       string   // Lead image URL
	Images      []string // List of article images (smart-filtered)
	Favicon     string   // Source favicon URL
	Language    string   // Detected language code
	Length      int      // Content length in characters
}
```

---

## Configuration

You can configure the client via `Option`s or via an explicit `Config`.

### Options

```go
c := retrieval.New(
	retrieval.WithTimeout(15*time.Second),
	retrieval.WithProxy("http://user:pass@host:port"), // Optional proxy
	retrieval.WithMaxErrBodyBytes(64*1024),
	retrieval.WithMaxBodyBytes(1024*1024),
)
```

### Config

```go
cfg := retrieval.DefaultConfig()
cfg.Timeout = 15 * time.Second
cfg.MaxErrBodyBytes = 64 * 1024
cfg.MaxBodyBytes = 1024 * 1024
cfg.Proxy = "http://localhost:8080"

c := retrieval.NewWithConfig(cfg)
```

### Proxy Support

You can easily route requests through a proxy (http, https, socks5) using `WithProxy`:

```go
c, err := retrieval.New(
    retrieval.WithProxy("http://user:pass@127.0.0.1:8080"),
)
```

### Custom HTTP client

If you need a custom transport (mTLS, custom DNS, etc.) that goes beyond simple proxying, pass your own `http.Client`:

```go
hc := &http.Client{ /* custom Transport, etc. */ }

c, err := retrieval.New(
	retrieval.WithClient(hc),
	// Optional: override timeout as well
	retrieval.WithTimeout(10*time.Second),
)
```

---

## Parser customization

Search parsing is abstracted behind:

```go
type Parser interface {
	Parse(reader io.ReadCloser) ([]Page, error)
}
```

You can supply your own parser if DuckDuckGo changes its markup or you want to target another HTML endpoint:

```go
c, err := retrieval.New(
	retrieval.WithParser(myParser{}),
)
```

---

## Error handling

When DuckDuckGo (or a target site) returns a non-2xx response, `retrieval-go` returns an error wrapping:

- `retrieval.ErrUnexpectedStatusCode`

When a target page is disallowed by `robots.txt` (and robots enforcement is enabled), `ParseContentFromLink` returns:

- `retrieval.ErrRobotsDenied`

You can handle it explicitly:

```go
doc, err := c.ParseContentFromLink(ctx, url, true)
if err != nil {
	if errors.Is(err, retrieval.ErrRobotsDenied) {
		// skip or handle in a policy-compliant way
		return
	}
	return
}
_ = doc
```

The error includes:

- status code
- request URL
- a truncated response body (up to `MaxErrBodyBytes`)

This is extremely helpful for debugging rate limits (429), blocks (403), and upstream failures (5xx).

---

## Notes for production & LLM usage

- **Rate limits / blocks**: DuckDuckGo may rate-limit or block aggressive traffic. Add caching, backoff, and reasonable concurrency.
- **Robots / Terms (IMPORTANT)**: Enforcement is now handled per-request via the third parameter of `ParseContentFromLink`. 
  - If `true`, the library will fetch and parse `robots.txt` before the main request.
  - If access is denied, it returns `retrieval.ErrRobotsDenied`.
  - It's **strongly recommended** to keep this enabled for third-party sites to comply with their crawling policies.
- **HTML extraction is heuristic**: readability works very well for articles, but not for all pages.
- **Security**: treat fetched content as untrusted input (sanitize before rendering; consider domain allow-lists).

---

## Contributing

PRs are welcome.

- Keep changes minimal and well-tested.
- Prefer backward-compatible API extensions.
- If you change parsing logic, include HTML fixtures.

Please read:

- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`

---

## License

MIT — see [`LICENSE`](./LICENSE).
