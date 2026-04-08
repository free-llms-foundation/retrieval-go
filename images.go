package retrieval

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetImagesPerQuery(ctx context.Context, query string) ([]Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.imagesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Set("q", query)
	q.Set("form", "HDRSC3")
	q.Set("first", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := c.sendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	reader, err := c.getDecodedReader(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get decoded reader: %w", err)
	}

	if reader != resp.Body {
		defer reader.Close()
	}

	limited := io.NopCloser(io.LimitReader(reader, c.maxBodyBytes))
	return c.imageParser.Parse(limited)
}
