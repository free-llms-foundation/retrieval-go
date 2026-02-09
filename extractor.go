package retrieval

import (
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	markdown, err := c.converter.ConvertString(page.Content)
	if err != nil {
		return nil, err
	}

	var images []string
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(page.Content)); err == nil {
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			if src, ok := s.Attr("src"); ok {
				if imgURL, err := url.Parse(src); err == nil {
					urlStr := strings.ToLower(imgURL.String())
					if strings.HasPrefix(urlStr, "data:") || strings.HasSuffix(urlStr, ".svg") || strings.HasSuffix(urlStr, ".gif") {
						return
					}

					if wStr, ok := s.Attr("width"); ok {
						if w, _ := strconv.Atoi(wStr); w > 0 && w < 120 {
							return
						}
					}
					if hStr, ok := s.Attr("height"); ok {
						if h, _ := strconv.Atoi(hStr); h > 0 && h < 120 {
							return
						}
					}

					trash := []string{
						"logo", "icon", "flag", "banner", "avatar", "spacer", "pixel", "ad-", "button",
						"social", "share", "nav", "menu", "footer", "header", "tracking",
						"bg-", "background", "placeholder", "sprite", "loader", "badge", "app-store",
						"google-play", "1x1", "transparent", "overlay", "theme", "plugin",
					}
					for _, kw := range trash {
						if strings.Contains(urlStr, kw) {
							return
						}
					}

					images = append(images, parsedURL.ResolveReference(imgURL).String())
				}
			}
		})
	}

	document := &Document{
		Title:    page.Title,
		Content:  markdown,
		Byline:   page.Byline,
		Length:   page.Length,
		Excerpt:  page.Excerpt,
		SiteName: page.SiteName,
		Image:    page.Image,
		Images:   images,
		Favicon:  page.Favicon,
		Language: page.Language,
	}

	return document, nil
}
