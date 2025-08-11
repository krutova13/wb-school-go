package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"wget/internal/types"
)

// HTTPDownloader реализует интерфейс Downloader для HTTP ресурсов
type HTTPDownloader struct {
	client    *http.Client
	userAgent string
}

// NewHTTPDownloader создает новый HTTP загрузчик
func NewHTTPDownloader(timeout time.Duration, userAgent string) *HTTPDownloader {
	return &HTTPDownloader{
		client: &http.Client{
			Timeout: timeout,
		},
		userAgent: userAgent,
	}
}

// Download загружает ресурс по указанному URL
func (d *HTTPDownloader) Download(ctx context.Context, resourceURL *url.URL) (*types.Resource, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", resourceURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("User-Agent", d.userAgent)
	req.Header.Set("Accept", "*/*")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP ошибка: %d %s", resp.StatusCode, resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	resourceType := d.determineResourceType(resourceURL, resp.Header.Get("Content-Type"))

	return &types.Resource{
		URL:      resourceURL,
		Type:     resourceType,
		Content:  content,
		MimeType: resp.Header.Get("Content-Type"),
	}, nil
}

func (d *HTTPDownloader) determineResourceType(u *url.URL, mimeType string) types.ResourceType {
	ext := strings.ToLower(path.Ext(u.Path))

	switch ext {
	case ".html", ".htm":
		return types.ResourceTypeHTML
	case ".css":
		return types.ResourceTypeCSS
	case ".js":
		return types.ResourceTypeJS
	case ".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp":
		return types.ResourceTypeImage
	}

	if strings.Contains(mimeType, "text/html") {
		return types.ResourceTypeHTML
	}
	if strings.Contains(mimeType, "text/css") {
		return types.ResourceTypeCSS
	}
	if strings.Contains(mimeType, "javascript") || strings.Contains(mimeType, "application/javascript") {
		return types.ResourceTypeJS
	}
	if strings.HasPrefix(mimeType, "image/") {
		return types.ResourceTypeImage
	}

	if u.Path == "" || strings.HasSuffix(u.Path, "/") {
		return types.ResourceTypeHTML
	}

	return types.ResourceTypeOther
}
