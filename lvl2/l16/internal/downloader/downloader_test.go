package downloader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"wget/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPDownloader_Download(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Test content</body></html>"))
	}))
	defer server.Close()

	downloader := NewHTTPDownloader(30*time.Second, "TestBot/1.0")

	testURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	ctx := context.Background()
	resource, err := downloader.Download(ctx, testURL)

	assert.NoError(t, err)
	assert.NotNil(t, resource)
	assert.Equal(t, testURL, resource.URL)
	assert.Equal(t, "text/html", resource.MimeType)
	assert.Contains(t, string(resource.Content), "Test content")
	assert.Equal(t, types.ResourceTypeHTML, resource.Type)
}

func TestHTTPDownloader_DetermineResourceType(t *testing.T) {
	downloader := NewHTTPDownloader(30*time.Second, "TestBot/1.0")

	tests := []struct {
		name     string
		url      string
		mimeType string
		expected types.ResourceType
	}{
		{
			name:     "HTML by extension",
			url:      "https://example.com/page.html",
			mimeType: "text/plain",
			expected: types.ResourceTypeHTML,
		},
		{
			name:     "CSS by extension",
			url:      "https://example.com/style.css",
			mimeType: "text/plain",
			expected: types.ResourceTypeCSS,
		},
		{
			name:     "JS by extension",
			url:      "https://example.com/script.js",
			mimeType: "text/plain",
			expected: types.ResourceTypeJS,
		},
		{
			name:     "Image by extension",
			url:      "https://example.com/image.png",
			mimeType: "text/plain",
			expected: types.ResourceTypeImage,
		},
		{
			name:     "HTML by MIME type",
			url:      "https://example.com/page",
			mimeType: "text/html",
			expected: types.ResourceTypeHTML,
		},
		{
			name:     "CSS by MIME type",
			url:      "https://example.com/style",
			mimeType: "text/css",
			expected: types.ResourceTypeCSS,
		},
		{
			name:     "JS by MIME type",
			url:      "https://example.com/script",
			mimeType: "application/javascript",
			expected: types.ResourceTypeJS,
		},
		{
			name:     "Image by MIME type",
			url:      "https://example.com/image",
			mimeType: "image/png",
			expected: types.ResourceTypeImage,
		},
		{
			name:     "Root path as HTML",
			url:      "https://example.com/",
			mimeType: "text/plain",
			expected: types.ResourceTypeHTML,
		},
		{
			name:     "Empty path as HTML",
			url:      "https://example.com",
			mimeType: "text/plain",
			expected: types.ResourceTypeHTML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testURL, err := url.Parse(tt.url)
			require.NoError(t, err)

			result := downloader.determineResourceType(testURL, tt.mimeType)
			assert.Equal(t, tt.expected, result)
		})
	}
}
