package retrieval

import (
	"fmt"
	"io"
	"net/http"
)

func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		limited := io.LimitReader(resp.Body, c.maxErrBodyBytes)
		body, _ := io.ReadAll(limited)
		resp.Body.Close()
		return nil, fmt.Errorf("%w: %d url=%s body: %s", ErrUnexpectedStatusCode, resp.StatusCode, req.URL.String(), string(body))
	}

	return resp, nil
}
