package internal

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/strikethrough"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

// Options holds conversion configuration.
type Options struct {
	Article   bool
	NoImages  bool
	Domain    string
}

// Convert reads the HTML file at path and returns Markdown.
func Convert(path string, opts Options) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	// Decode charset to UTF-8.
	decoded, err := decodeToUTF8(raw)
	if err != nil {
		decoded = string(raw)
	}

	var htmlContent string

	if opts.Article {
		baseURL := opts.Domain
		htmlContent, err = ExtractArticle(decoded, baseURL)
		if err != nil {
			return "", fmt.Errorf("extract article: %w", err)
		}
	} else {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(decoded))
		if err != nil {
			return "", fmt.Errorf("parse html: %w", err)
		}

		if opts.NoImages {
			doc.Find("img").Remove()
		}

		StripNoise(doc)

		htmlContent, err = doc.Html()
		if err != nil {
			return "", fmt.Errorf("serialize html: %w", err)
		}
	}

	if opts.NoImages && opts.Article {
		// Strip images from article HTML too.
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err == nil {
			doc.Find("img").Remove()
			if h, e := doc.Html(); e == nil {
				htmlContent = h
			}
		}
	}

	result, err := convertHTML(htmlContent)
	if err != nil {
		return "", fmt.Errorf("convert to markdown: %w", err)
	}

	return result, nil
}

func convertHTML(htmlContent string) (string, error) {
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			table.NewTablePlugin(),
			strikethrough.NewStrikethroughPlugin(),
		),
	)

	result, err := conv.ConvertString(htmlContent)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func decodeToUTF8(raw []byte) (string, error) {
	reader, err := charset.NewReader(bytes.NewReader(raw), "")
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
