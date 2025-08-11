package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HTMLParser реализует интерфейс Parser для HTML документов
type HTMLParser struct{}

// NewHTMLParser создает новый HTML парсер
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// ParseHTML парсит HTML и извлекает все ссылки
func (p *HTMLParser) ParseHTML(content []byte, baseURL *url.URL) ([]*url.URL, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	var urls []*url.URL
	seen := make(map[string]bool)

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			if parsedURL := p.resolveURL(href, baseURL); parsedURL != nil {
				key := parsedURL.String()
				if !seen[key] {
					seen[key] = true
					urls = append(urls, parsedURL)
				}
			}
		}
	})

	doc.Find("[src]").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			if parsedURL := p.resolveURL(src, baseURL); parsedURL != nil {
				key := parsedURL.String()
				if !seen[key] {
					seen[key] = true
					urls = append(urls, parsedURL)
				}
			}
		}
	})

	doc.Find("link[href]").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			if parsedURL := p.resolveURL(href, baseURL); parsedURL != nil {
				key := parsedURL.String()
				if !seen[key] {
					seen[key] = true
					urls = append(urls, parsedURL)
				}
			}
		}
	})

	return urls, nil
}

// ParseCSS парсит CSS и извлекает все URL из правил url()
func (p *HTMLParser) ParseCSS(content []byte, baseURL *url.URL) ([]*url.URL, error) {
	var urls []*url.URL
	seen := make(map[string]bool)

	urlRegex := regexp.MustCompile(`url\(['"]?([^'")\s]+)['"]?\)`)
	matches := urlRegex.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) > 1 {
			urlStr := match[1]
			if parsedURL := p.resolveURL(urlStr, baseURL); parsedURL != nil {
				key := parsedURL.String()
				if !seen[key] {
					seen[key] = true
					urls = append(urls, parsedURL)
				}
			}
		}
	}

	return urls, nil
}

func (p *HTMLParser) resolveURL(urlStr string, baseURL *url.URL) *url.URL {
	if urlStr == "" || strings.HasPrefix(urlStr, "#") {
		return nil
	}

	if strings.HasPrefix(urlStr, "javascript:") || strings.HasPrefix(urlStr, "data:") {
		return nil
	}

	if strings.HasPrefix(urlStr, "mailto:") || strings.HasPrefix(urlStr, "tel:") {
		return nil
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}

	if parsedURL.IsAbs() {
		return parsedURL
	}

	resolvedURL := baseURL.ResolveReference(parsedURL)
	return resolvedURL
}
