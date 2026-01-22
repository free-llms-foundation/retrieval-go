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
	Favicon     string
	Language    string
}

type Page struct {
	Title   string
	Link    string
	Snippet string
	Favicon string
}
