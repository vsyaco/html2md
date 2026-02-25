package internal

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-shiori/go-readability"
)

var noiseSelectors = []string{
	"nav", "footer", "header", "aside", "form", "button", "noscript", "script", "style",
}

var noiseClassKeywords = []string{
	"ad", "ads", "advert", "sidebar", "cookie", "menu", "popup", "banner", "promo",
	"newsletter", "subscribe", "related", "share", "social", "comment", "widget",
}

// StripNoise removes noisy elements from the document in-place.
func StripNoise(doc *goquery.Document) {
	for _, sel := range noiseSelectors {
		doc.Find(sel).Remove()
	}

	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		id, _ := s.Attr("id")
		combined := strings.ToLower(class + " " + id)
		for _, kw := range noiseClassKeywords {
			if strings.Contains(combined, kw) {
				s.Remove()
				return
			}
		}
	})
}

// ExtractArticle uses go-readability to extract the main article content as HTML.
func ExtractArticle(rawHTML string, rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || rawURL == "" {
		parsedURL = &url.URL{}
	}

	article, err := readability.FromReader(strings.NewReader(rawHTML), parsedURL)
	if err != nil {
		return "", err
	}

	return article.Content, nil
}
