package retrieval

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) SearchWithQuery(ctx context.Context, query string) ([]Page, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	values := req.URL.Query()
	values.Set("q", query)
	req.URL.RawQuery = values.Encode()

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
	return c.parser.Parse(limited)
}
