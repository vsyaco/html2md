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
		if matchesNoiseKeyword(class) || matchesNoiseKeyword(id) {
			s.Remove()
			return
		}
	})
}

// matchesNoiseKeyword checks if any class/id token starts with a noise keyword.
// Splits by whitespace into individual tokens, then checks each token's word parts
// (split by - and _) to see if the first part matches a noise keyword.
func matchesNoiseKeyword(attr string) bool {
	if attr == "" {
		return false
	}
	for _, token := range strings.Fields(strings.ToLower(attr)) {
		parts := strings.FieldsFunc(token, func(r rune) bool {
			return r == '-' || r == '_'
		})
		if len(parts) == 0 {
			continue
		}
		first := parts[0]
		for _, kw := range noiseClassKeywords {
			if first == kw {
				return true
			}
		}
	}
	return false
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
