package retrieval

import (
	"errors"

	"golang.org/x/net/html"
)

var ErrUnexpectedStatusCode = errors.New("unexpected status code")

type Document struct {
	Title       string
	Byline      string
	Node        *html.Node
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
}
