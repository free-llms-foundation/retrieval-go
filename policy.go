package retrieval

import (
	"context"
	"net/http"
	"net/url"

	"github.com/temoto/robotstxt"
)

func (c *Client) allowedByRobots(ctx context.Context, target string) bool {
	u, err := url.Parse(target)
	if err != nil || u.Host == "" || u.Scheme == "" {
		return false
	}

	robotsTxt := u.Scheme + "://" + u.Host + "/robots.txt"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, robotsTxt, nil)
	if err != nil {
		return false
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false
	}

	robots, err := robotstxt.FromResponse(resp)
	if err != nil {
		return false
	}

	path := u.EscapedPath()
	if path == "" {
		path = "/"
	}

	return robots.TestAgent(path, c.userAgent)
}
