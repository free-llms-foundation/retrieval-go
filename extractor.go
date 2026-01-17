package retrieval

import (
	"io"
	"net/url"

	"github.com/go-shiori/go-readability"
)

func (c *Client) extractFromURL(u string, body io.Reader) (*Document, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	article, err := readability.FromReader(body, parsedURL)
	if err != nil {
		return nil, err
	}

	doc := &Document{
		Title:    article.Title,
		Content:  article.TextContent,
		Byline:   article.Byline,
		Node:     article.Node,
		Length:   article.Length,
		Excerpt:  article.Excerpt,
		SiteName: article.SiteName,
		Image:    article.Image,
		Favicon:  article.Favicon,
		Language: article.Language,
	}

	return doc, nil
}
