//go:build integration

package retrieval

import (
	"context"
	"testing"
	"time"
)

// Run with: go test -tags=integration -timeout 60s ./...

func integrationClient(t *testing.T) *Client {
	t.Helper()
	c, err := New()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	return c
}

func integrationCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func TestIntegration_SearchWithQuery(t *testing.T) {
	c := integrationClient(t)
	ctx := integrationCtx(t)

	pages, err := c.SearchWithQuery(ctx, "golang", "")
	if err != nil {
		t.Fatalf("SearchWithQuery: %v", err)
	}
	if len(pages) == 0 {
		t.Fatal("SearchWithQuery: got 0 results, parser may be broken")
	}
	for _, p := range pages {
		if p.Link == "" {
			t.Error("result with empty Link")
		}
		if p.Title == "" {
			t.Error("result with empty Title")
		}
		if p.Snippet == "" {
			t.Error("result with empty Snippet")
		}
		if p.Source == "" {
			t.Error("result with empty Source")
		}

		if p.Link == "" {
			t.Error("result with empty Link")
		}
	}
	t.Logf("got %d search results", len(pages))
}

func TestIntegration_SearchWithQuery_DateFilter(t *testing.T) {
	c := integrationClient(t)
	ctx := integrationCtx(t)

	tests := []struct{ name, filter string }{
		{"day", "d"},
		{"week", "w"},
		{"month", "m"},
		{"year", "y"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages, err := c.SearchWithQuery(ctx, "golang", tt.filter)
			if err != nil {
				t.Fatalf("SearchWithQuery(%q): %v", tt.filter, err)
			}
			t.Logf("filter=%q results=%d", tt.filter, len(pages))
		})
	}
}

func TestIntegration_SearchImagesWithQuery(t *testing.T) {
	c := integrationClient(t)
	ctx := integrationCtx(t)

	images, err := c.SearchImagesWithQuery(ctx, "golang gopher", "")
	if err != nil {
		t.Fatalf("SearchImagesWithQuery: %v", err)
	}
	if len(images) == 0 {
		t.Fatal("SearchImagesWithQuery: got 0 results, Bing parser may be broken")
	}
	for _, img := range images {
		if img.URL == "" {
			t.Error("image with empty URL")
		}
		if img.Desc == "" {
			t.Error("image with empty Desc")
		}
		if img.Thumbnail == "" {
			t.Error("image with empty Thumbnail")
		}
		if img.PageURL == "" {
			t.Error("image with empty PageURL")
		}
	}
	t.Logf("got %d images", len(images))
}

func TestIntegration_SearchImagesWithQuery_DateFilter(t *testing.T) {
	c := integrationClient(t)
	ctx := integrationCtx(t)

	noFilter, err := c.SearchImagesWithQuery(ctx, "cats", "")
	if err != nil {
		t.Fatalf("no filter: %v", err)
	}

	year, err := c.SearchImagesWithQuery(ctx, "cats", "y")
	if err != nil {
		t.Fatalf("year filter: %v", err)
	}

	t.Logf("no filter: %d images, year filter: %d images", len(noFilter), len(year))

	if len(year) == 0 {
		t.Error("year-filtered images: got 0 results")
	}

	if len(year) > len(noFilter) {
		t.Errorf("year filter had no effect: filtered=%d > unfiltered=%d", len(year), len(noFilter))
	}
}

func TestIntegration_ParseContentFromLink(t *testing.T) {
	c := integrationClient(t)
	ctx := integrationCtx(t)

	doc, err := c.ParseContentFromLink(ctx, "https://go.dev/blog/", false)
	if err != nil {
		t.Fatalf("ParseContentFromLink: %v", err)
	}
	if doc.Content == "" {
		t.Error("parsed document has empty Content")
	}
	if doc.Title == "" {
		t.Error("parsed document has empty Title")
	}
	t.Logf("title: %q, content length: %d chars", doc.Title, len(doc.Content))
}
