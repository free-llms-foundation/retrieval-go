package retrieval

import "errors"

var ErrUnexpectedStatusCode = errors.New("unexpected status code")

type Document struct {
	Title   string
	Content string
	URL     string
}

type Page struct {
	Title   string
	Link    string
	Snippet string
}
