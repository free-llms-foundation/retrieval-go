package retrieval

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
)

var ErrRobotsDenied = errors.New("robots.txt denied")

func (c *Client) ParseContentFromLink(ctx context.Context, link string) (*Document, error) {
	if link == "" {
		return nil, errors.New("link cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	reader, err := c.getDecodedReader(resp)
	if err != nil {
		return nil, err
	}

	if reader != resp.Body {
		defer reader.Close()
	}

	if c.respectRobots && !c.allowedByRobots(ctx, link) {
		return nil, ErrRobotsDenied
	}

	limited := io.LimitReader(reader, c.maxBodyBytes)
	doc, err := c.extractFromURL(link, limited)
	if err != nil {
		return nil, err
	}

	doc.Content = strings.Join(strings.Fields(strings.TrimSpace(doc.Content)), " ")

	return doc, nil
}
