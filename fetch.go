package retrieval

import (
	"context"
	"errors"
	"io"
	"net/http"
)

func (c *Client) ParseContentFromLink(ctx context.Context, link string, robotsTxtAllowed bool) (*Document, error) {
	if link == "" {
		return nil, errors.New("link cannot be empty")
	}

	if robotsTxtAllowed {
		allowed, err := c.allowedByRobots(ctx, link)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, ErrRobotsDenied
		}
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

	limited := io.LimitReader(reader, c.maxBodyBytes)
	doc, err := c.extractFromURL(link, limited)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
