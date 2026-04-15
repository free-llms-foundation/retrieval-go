package retrieval

import (
	"os"
	"strings"
	"testing"
)

func extractorClient(t *testing.T) *Client {
	t.Helper()
	c, err := New()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	return c
}

var trashURLFragments = []string{
	"logo", "icon", "banner", "avatar", "tracking",
	"nav-", "social", "placeholder", "background", "badge",
}

func TestExtractFromURL_TrashImagesFiltered(t *testing.T) {
	t.Parallel()

	f, err := os.Open("testdata/extractor_article.html")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	c := extractorClient(t)
	doc, err := c.extractFromURL("https://example.com/article", f)
	if err != nil {
		t.Fatalf("extractFromURL: %v", err)
	}

	for _, imgURL := range doc.Images {
		lower := strings.ToLower(imgURL)

		// data URIs must never appear.
		if strings.HasPrefix(lower, "data:") {
			t.Errorf("data URI leaked into images: %s", imgURL)
		}

		// SVG must be filtered.
		if strings.HasSuffix(lower, ".svg") {
			t.Errorf("SVG image leaked: %s", imgURL)
		}

		// GIF must be filtered.
		if strings.HasSuffix(lower, ".gif") {
			t.Errorf("GIF image leaked: %s", imgURL)
		}

		// Trash keyword URLs must be filtered.
		for _, kw := range trashURLFragments {
			if strings.Contains(lower, kw) {
				t.Errorf("trash keyword %q leaked via URL: %s", kw, imgURL)
			}
		}
	}
}

func TestExtractFromURL_ValidImagesPresent(t *testing.T) {
	t.Parallel()

	f, err := os.Open("testdata/extractor_article.html")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	c := extractorClient(t)
	doc, err := c.extractFromURL("https://example.com/article", f)
	if err != nil {
		t.Fatalf("extractFromURL: %v", err)
	}

	validPrefixes := []string{
		"https://example.com/images/go-gopher.png",
		"https://example.com/diagrams/concurrency-channels.jpg",
		"https://example.com/screenshots/stdlib-coverage.png",
	}

	found := 0
	prefixIndex := 0
	for _, imgURL := range doc.Images {
		if prefixIndex < len(validPrefixes) && imgURL == validPrefixes[prefixIndex] {
			found++
			prefixIndex++
		}
	}

	if found == 0 {
		t.Errorf("no valid content images survived filtering; got images: %v", doc.Images)
	}
}

func TestExtractFromURL_BasicFields(t *testing.T) {
	t.Parallel()

	f, err := os.Open("testdata/extractor_article.html")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	c := extractorClient(t)
	doc, err := c.extractFromURL("https://example.com/article", f)
	if err != nil {
		t.Fatalf("extractFromURL: %v", err)
	}

	if doc.Title == "" {
		t.Error("document Title is empty")
	}
	if doc.Content == "" {
		t.Error("document Content (markdown) is empty")
	}

	if doc.Language == "" {
		t.Error("document Language is empty")
	}

	// Content should be markdown, not raw HTML.
	if strings.Contains(doc.Content, "<p>") {
		t.Errorf("Content appears to contain raw HTML <p> tags: %q", doc.Content[:min(200, len(doc.Content))])
	}
	// Table should be converted to markdown, not left as HTML.
	if strings.Contains(doc.Content, "<table>") {
		t.Errorf("Content contains raw <table> tag — table plugin may not be active")
	}
}

func TestExtractFromURL_TableInContent(t *testing.T) {
	t.Parallel()

	f, err := os.Open("testdata/extractor_article.html")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	c := extractorClient(t)
	doc, err := c.extractFromURL("https://example.com/article", f)
	if err != nil {
		t.Fatalf("extractFromURL: %v", err)
	}

	tableMarkers := []string{
		"Feature", "Goroutines", "Garbage collection", "Static typing",
	}
	for _, marker := range tableMarkers {
		if !strings.Contains(doc.Content, marker) {
			t.Errorf("table marker %q not found in markdown content", marker)
		}
	}

	if !strings.Contains(doc.Content, "|") {
		t.Errorf("no pipe characters found in content — markdown table may be missing:\n%s", doc.Content)
	}
}

// TestExtractFromURL_NoImagesURLsDuplicated checks that the image URL list
// contains no duplicates (the extractor itself deduplicates via URL parsing,
// but we want to guard against regressions).
func TestExtractFromURL_NoImagesURLsDuplicated(t *testing.T) {
	t.Parallel()

	f, err := os.Open("testdata/extractor_article.html")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	c := extractorClient(t)
	doc, err := c.extractFromURL("https://example.com/article", f)
	if err != nil {
		t.Fatalf("extractFromURL: %v", err)
	}

	seen := make(map[string]struct{}, len(doc.Images))
	for _, imgURL := range doc.Images {
		if _, dup := seen[imgURL]; dup {
			t.Errorf("duplicate image URL: %s", imgURL)
		}
		seen[imgURL] = struct{}{}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
