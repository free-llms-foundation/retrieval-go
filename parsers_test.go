package retrieval

import (
	"io"
	"os"
	"strings"
	"testing"
)

// --- DefaultDDGParser ---

func TestDefaultDDGParser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		file       string
		wantMinLen int
		wantErrNil bool
		wantNoAds  bool // no duckduckgo.com/y.js links
	}{
		{
			name:       "real DDG lite response",
			file:       "testdata/ddg_results.html",
			wantMinLen: 5,
			wantErrNil: true,
			wantNoAds:  true,
		},
		{
			name:       "empty page",
			file:       "",
			wantMinLen: 0,
			wantErrNil: true,
			wantNoAds:  true,
		},
	}

	parser := &DefaultDDGParser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var r io.ReadCloser
			if tt.file != "" {
				f, err := os.Open(tt.file)
				if err != nil {
					t.Fatalf("open fixture: %v", err)
				}
				r = f
			} else {
				r = io.NopCloser(strings.NewReader("<html><body></body></html>"))
			}

			pages, err := parser.Parse(r)

			if tt.wantErrNil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(pages) < tt.wantMinLen {
				t.Errorf("got %d results, want at least %d", len(pages), tt.wantMinLen)
			}

			if tt.wantNoAds {
				for _, p := range pages {
					if strings.Contains(p.Link, "duckduckgo.com/y.js") {
						t.Errorf("ad link leaked into results: %s", p.Link)
					}
					if strings.Contains(p.Link, "ad_provider") {
						t.Errorf("ad_provider link leaked into results: %s", p.Link)
					}
				}
			}

			for _, p := range pages {
				if p.Link == "" {
					t.Error("page with empty link")
				}
				if p.Title == "" {
					t.Error("page with empty title")
				}
				if p.Snippet == "" {
					t.Error("page with empty snippet")
				}
				if p.Source == "" {
					t.Error("page with empty source hostname")
				}
			}
		})
	}
}

func TestDefaultDDGParser_NoDuplicates(t *testing.T) {
	t.Parallel()

	f, err := os.Open("testdata/ddg_results.html")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	pages, err := (&DefaultDDGParser{}).Parse(f)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	seen := make(map[string]struct{}, len(pages))
	for _, p := range pages {
		if _, dup := seen[p.Link]; dup {
			t.Errorf("duplicate link: %s", p.Link)
		}
		if p.Link == "" {
			t.Error("page with empty link")
		}
		if p.Title == "" {
			t.Error("page with empty title")
		}
		if p.Snippet == "" {
			t.Error("page with empty snippet")
		}
		if p.Source == "" {
			t.Error("page with empty source hostname")
		}
		seen[p.Link] = struct{}{}
	}
}

// --- BingImagesParser ---

func TestBingImagesParser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		file       string
		wantMinLen int
		wantErrNil bool
	}{
		{
			name:       "real Bing response no filter",
			file:       "testdata/bing_images.html",
			wantMinLen: 10,
			wantErrNil: true,
		},
		{
			name:       "real Bing response year filter",
			file:       "testdata/bing_images_year.html",
			wantMinLen: 1,
			wantErrNil: true,
		},
		{
			name:       "empty page",
			file:       "",
			wantMinLen: 0,
			wantErrNil: true,
		},
	}

	parser := &BingImagesParser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var r io.ReadCloser
			if tt.file != "" {
				f, err := os.Open(tt.file)
				if err != nil {
					t.Fatalf("open fixture: %v", err)
				}
				r = f
			} else {
				r = io.NopCloser(strings.NewReader("<html><body></body></html>"))
			}

			images, err := parser.Parse(r)

			if tt.wantErrNil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErrNil && err == nil {
				t.Fatal("expected error, got nil")
			}

			if len(images) < tt.wantMinLen {
				t.Errorf("got %d results, want at least %d", len(images), tt.wantMinLen)
			}

			for _, img := range images {
				if img.URL == "" {
					t.Error("image with empty URL")
				}
				if img.Thumbnail == "" {
					t.Error("image with empty Thumbnail")
				}
				if img.PageURL == "" {
					t.Error("image with empty PageURL")
				}
			}
		})
	}
}

func TestBingImagesParser_NoDuplicates(t *testing.T) {
	t.Parallel()

	f, err := os.Open("testdata/bing_images.html")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	images, err := (&BingImagesParser{}).Parse(f)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	seenImages := make(map[string]struct{}, len(images))
	for _, img := range images {
		if _, dup := seenImages[img.URL]; dup {
			t.Errorf("duplicate image URL: %s", img.URL)
		}
		if img.URL == "" {
			t.Error("image with empty URL")
		}
		if img.Thumbnail == "" {
			t.Error("image with empty Thumbnail")
		}
		if img.PageURL == "" {
			t.Error("image with empty PageURL")
		}
		seenImages[img.URL] = struct{}{}
	}
}

func TestBingImagesParser_YearFilterFewerResults(t *testing.T) {
	t.Parallel()

	parseImages := func(path string) []Image {
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("open %s: %v", path, err)
		}
		defer f.Close()
		imgs, err := (&BingImagesParser{}).Parse(f)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		return imgs
	}

	all := parseImages("testdata/bing_images.html")
	year := parseImages("testdata/bing_images_year.html")

	if len(year) >= len(all) {
		t.Errorf("year-filtered results (%d) should be fewer than unfiltered (%d)", len(year), len(all))
	}

	seenYear := make(map[string]struct{}, len(year))
	for i, img := range year {
		if img.URL == "" {
			t.Errorf("year[%d]: empty URL", i)
		}
		if img.Thumbnail == "" {
			t.Errorf("year[%d]: empty Thumbnail (url=%s)", i, img.URL)
		}
		if img.PageURL == "" {
			t.Errorf("year[%d]: empty PageURL (url=%s)", i, img.URL)
		}

		if _, dup := seenYear[img.URL]; dup {
			t.Errorf("year[%d]: duplicate URL %s", i, img.URL)
		}
		seenYear[img.URL] = struct{}{}
	}
}
