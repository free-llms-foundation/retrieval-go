package retrieval

import (
	"io"
	"net/url"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/go-shiori/go-readability"
)

func (c *Client) extractFromURL(u string, body io.Reader) (*Document, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	page, err := readability.FromReader(body, parsedURL)
	if err != nil {
		return nil, err
	}

	markdown, err := htmltomarkdown.ConvertString(page.Content)
	if err != nil {
		return nil, err
	}

	doc := &Document{
		Title:    page.Title,
		Content:  markdown,
		Byline:   page.Byline,
		Length:   page.Length,
		Excerpt:  page.Excerpt,
		SiteName: page.SiteName,
		Image:    page.Image,
		Favicon:  page.Favicon,
		Language: page.Language,
	}

	return doc, nil
}
