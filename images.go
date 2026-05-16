package retrieval

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// bingImageDateFilter maps DuckDuckGo-style date shorthand (d/w/m/y) to Bing
// Images' qft filter values. The colon in filterui:age-lt must NOT be
// percent-encoded, so these strings are appended to RawQuery directly.
var bingImageDateFilter = map[string]string{
	"d": "+filterui:age-lt1440",   // past 24 hours
	"w": "+filterui:age-lt10080",  // past week
	"m": "+filterui:age-lt43200",  // past month
	"y": "+filterui:age-lt525600", // past year
}

func (c *Client) SearchImagesWithQuery(ctx context.Context, query string, dateFilter string) ([]Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.imagesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Build RawQuery manually to ensure 'q' is FIRST.
	// Order of parameters is often a signal for bot detection.
	rawQuery := "q=" + url.QueryEscape(query)
	rawQuery += "&form=QBIDMH"
	rawQuery += "&first=1"
	rawQuery += "&adlt=strict"
	req.URL.RawQuery = rawQuery
	// Bing requires the colon in filterui:age-lt to be literal (not %3A),
	// so we append qft directly to RawQuery instead of using url.Values.
	if f, ok := bingImageDateFilter[dateFilter]; ok {
		// Append qft as a raw string — the colon and plus must stay literal;
		// url.QueryEscape would encode them and Bing would ignore the filter.
		req.URL.RawQuery += "&qft=" + f
	}

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
