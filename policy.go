package retrieval

import (
	"context"
	"net/http"
	"net/url"

	"github.com/temoto/robotstxt"
)

func (c *Client) allowedByRobots(ctx context.Context, target string) (bool, error) {
	u, err := url.Parse(target)
	if err != nil || u.Host == "" || u.Scheme == "" {
		return false, nil
	}

	robotsTxt := u.Scheme + "://" + u.Host + "/robots.txt"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, robotsTxt, nil)
	if err != nil {
		return false, err
	}

	for _, header := range c.headers {
		req.Header.Set(header[0], header[1])
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, nil
	}

	robots, err := robotstxt.FromResponse(resp)
	if err != nil {
		return false, err
	}

	path := u.EscapedPath()
	if path == "" {
		path = "/"
	}

	return robots.TestAgent(path, "*"), nil
}
