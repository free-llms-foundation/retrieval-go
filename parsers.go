package retrieval

import (
	"encoding/json"
	"fmt"
	"html"
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

		if strings.Contains(finalLink, "duckduckgo.com/y.js") || strings.Contains(finalLink, "ad_provider") {
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

		u, err = url.Parse(finalLink)
		if err != nil {
			return
		}

		pages = append(pages, Page{
			Link:    finalLink,
			Source:  u.Hostname(),
			Title:   strings.Join(strings.Fields(s.Text()), " "),
			Snippet: strings.Join(strings.Fields(snippet), " "),
			Favicon: fmt.Sprintf(defaultFavicon, u.Host),
		})
	})

	return pages, nil
}

type BingImagesParser struct{}

type bingImageMeta struct {
	Murl string `json:"murl"`
	Turl string `json:"turl"`
	Purl string `json:"purl"`
	T    string `json:"t"`
	Desc string `json:"desc"`
}

func (p *BingImagesParser) Parse(reader io.ReadCloser) ([]Image, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("bing images: failed to parse HTML: %w", err)
	}

	var images []Image
	seen := make(map[string]struct{})

	doc.Find(".iusc").Each(func(_ int, s *goquery.Selection) {
		mRaw, exists := s.Attr("m")
		if !exists {
			return
		}

		var meta bingImageMeta
		if err := json.Unmarshal([]byte(html.UnescapeString(mRaw)), &meta); err != nil {
			return
		}

		if meta.Murl == "" || meta.Turl == "" || meta.Purl == "" {
			return
		}

		if _, dup := seen[meta.Murl]; dup {
			return
		}
		seen[meta.Murl] = struct{}{}

		images = append(images, Image{
			URL:       meta.Murl,
			Thumbnail: meta.Turl,
			PageURL:   meta.Purl,
			Desc:      strings.TrimSpace(meta.T),
		})
	})

	return images, nil
}
