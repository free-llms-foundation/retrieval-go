package retrieval

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type DefaultDDGParser struct{}

func (p *DefaultDDGParser) Parse(reader io.ReadCloser) ([]Page, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	pages := make([]Page, 0, 10)
	prefixes := map[string]struct{}{}

	doc.Find("a[href*='/l/?uddg=']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		if strings.HasPrefix(href, "//") {
			href = "https:" + href
		}
		if strings.HasPrefix(href, "/") {
			href = "https://duckduckgo.com" + href
		}

		finalLink := href
		u, err := url.Parse(href)
		if err != nil {
			return
		}

		target := u.Query().Get("uddg")
		if target != "" {
			finalLink = target
		}

		if finalLink == "" {
			return
		}

		if _, ok := prefixes[finalLink]; ok {
			return
		}

		prefixes[finalLink] = struct{}{}

		snippet := ""
		tr := s.ParentsFiltered("tr").First()
		if tr.Length() > 0 {
			snippetSelection := tr.Next().Find("td.result-snippet")
			if snippetSelection.Length() > 0 {
				snippet = snippetSelection.Text()
			}
		}

		pages = append(pages, Page{
			Link:    finalLink,
			Title:   strings.Join(strings.Fields(s.Text()), " "),
			Snippet: strings.Join(strings.Fields(snippet), " "),
		})
	})

	return pages, nil
}
