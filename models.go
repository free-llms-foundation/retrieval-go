package retrieval

import (
	"errors"
)

var ErrUnexpectedStatusCode = errors.New("unexpected status code")

type Document struct {
	Title       string
	Byline      string
	Content     string
	TextContent string
	Length      int
	Excerpt     string
	SiteName    string
	Image       string
	Images      []string
	Favicon     string
	Language    string
}

type Page struct {
	Link    string
	Title   string
	Source  string
	Snippet string
	Favicon string
}
