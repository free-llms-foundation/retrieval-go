package retrieval

import (
	"fmt"
	"io"
	"net/http"
)

func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
	for _, header := range c.headers {
		if req.Header.Get(header[0]) == "" {
			req.Header.Set(header[0], header[1])
		}
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		limited := io.LimitReader(resp.Body, c.maxErrBodyBytes)
		body, _ := io.ReadAll(limited)
		resp.Body.Close()
		return nil, fmt.Errorf("%w: %d url=%s body: %s", ErrUnexpectedStatusCode, resp.StatusCode, req.URL.String(), string(body))
	}

	return resp, nil
}
