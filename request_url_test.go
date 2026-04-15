package retrieval

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type mockSearchParser struct{}

func (m *mockSearchParser) Parse(r io.ReadCloser) ([]Page, error) { return nil, nil }

type mockImageParser struct{}

func (m *mockImageParser) Parse(r io.ReadCloser) ([]Image, error) { return nil, nil }

// captureURL returns an httptest.Server that records the last request URL and
// responds with 200 OK.
func captureURL(t *testing.T) (*httptest.Server, *url.URL) {
	t.Helper()
	var lastURL url.URL
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastURL = *r.URL
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "<html><body></body></html>")
	}))
	t.Cleanup(srv.Close)
	return srv, &lastURL
}

// --- SearchWithQuery URL tests ---

func TestSearchWithQuery_URLParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		query      string
		dateFilter string
		wantQ      string
		wantDF     string // expected "df" param; empty means absent
	}{
		{
			name:       "query only",
			query:      "golang testing",
			dateFilter: "",
			wantQ:      "golang testing",
			wantDF:     "",
		},
		{
			name:       "with date filter d",
			query:      "news",
			dateFilter: "d",
			wantQ:      "news",
			wantDF:     "d",
		},
		{
			name:       "with date filter y",
			query:      "cats",
			dateFilter: "y",
			wantQ:      "cats",
			wantDF:     "y",
		},
		{
			name:       "special chars in query",
			query:      "go & generics > 1.18",
			dateFilter: "",
			wantQ:      "go & generics > 1.18",
			wantDF:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv, lastURL := captureURL(t)

			c, err := New(
				WithBaseURL(srv.URL),
				WithSearchParser(&mockSearchParser{}),
			)
			if err != nil {
				t.Fatalf("new client: %v", err)
			}

			_, _ = c.SearchWithQuery(t.Context(), tt.query, tt.dateFilter)

			q := lastURL.Query()

			if got := q.Get("q"); got != tt.wantQ {
				t.Errorf("q param: got %q, want %q", got, tt.wantQ)
			}

			if tt.wantDF == "" {
				if q.Has("df") {
					t.Errorf("df param should be absent, got %q", q.Get("df"))
				}
			} else {
				if got := q.Get("df"); got != tt.wantDF {
					t.Errorf("df param: got %q, want %q", got, tt.wantDF)
				}
			}
		})
	}
}

// --- SearchImagesWithQuery URL tests ---

func TestSearchImagesWithQuery_URLParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		query      string
		dateFilter string
		wantQ      string
		wantQFT    string // expected raw "qft" value; empty means absent
	}{
		{
			name:       "no date filter",
			query:      "cats",
			dateFilter: "",
			wantQ:      "cats",
			wantQFT:    "",
		},
		{
			name:       "day filter",
			query:      "cats",
			dateFilter: "d",
			wantQ:      "cats",
			wantQFT:    "+filterui:age-lt1440",
		},
		{
			name:       "week filter",
			query:      "cats",
			dateFilter: "w",
			wantQ:      "cats",
			wantQFT:    "+filterui:age-lt10080",
		},
		{
			name:       "month filter",
			query:      "cats",
			dateFilter: "m",
			wantQ:      "cats",
			wantQFT:    "+filterui:age-lt43200",
		},
		{
			name:       "year filter",
			query:      "cats",
			dateFilter: "y",
			wantQ:      "cats",
			wantQFT:    "+filterui:age-lt525600",
		},
		{
			name:       "unknown filter ignored",
			query:      "cats",
			dateFilter: "century",
			wantQ:      "cats",
			wantQFT:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv, lastURL := captureURL(t)

			c, err := New(
				WithCustomImagesURL(srv.URL),
				WithImageParser(&mockImageParser{}),
			)
			if err != nil {
				t.Fatalf("new client: %v", err)
			}

			_, _ = c.SearchImagesWithQuery(t.Context(), tt.query, tt.dateFilter)

			raw := lastURL.RawQuery

			q := lastURL.Query()
			if got := q.Get("q"); got != tt.wantQ {
				t.Errorf("q param: got %q, want %q", got, tt.wantQ)
			}

			if got := q.Get("form"); got != "HDRSC3" {
				t.Errorf("form param: got %q, want %q", got, "HDRSC3")
			}
			if got := q.Get("first"); got != "1" {
				t.Errorf("first param: got %q, want %q", got, "1")
			}

			if tt.wantQFT == "" {
				if strings.Contains(raw, "qft=") {
					t.Errorf("qft should be absent in URL %q", raw)
				}
			} else {

				want := "qft=" + tt.wantQFT
				if !strings.Contains(raw, want) {
					t.Errorf("raw query %q does not contain %q", raw, want)
				}
			}
		})
	}
}
