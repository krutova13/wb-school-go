package storage

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"wget/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_SaveAndExists(t *testing.T) {
	tempDir := t.TempDir()

	storage, err := NewFileStorage(tempDir)
	require.NoError(t, err)

	testURL, _ := url.Parse("https://example.com/page.html")
	resource := &types.Resource{
		URL:      testURL,
		Content:  []byte("<html><body>Test content</body></html>"),
		MimeType: "text/html",
		Type:     types.ResourceTypeHTML,
	}

	err = storage.Save(resource)
	require.NoError(t, err)

	localPath := storage.GetLocalPath(testURL)
	assert.True(t, storage.Exists(localPath))

	content, err := os.ReadFile(localPath)
	require.NoError(t, err)
	assert.Equal(t, resource.Content, content)
}

func TestFileStorage_GetLocalPath(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewFileStorage(tempDir)
	require.NoError(t, err)

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "HTML page",
			url:      "https://example.com/page.html",
			expected: filepath.Join(tempDir, "example.com", "page.html"),
		},
		{
			name:     "Root path",
			url:      "https://example.com/",
			expected: filepath.Join(tempDir, "example.com", "index.html"),
		},
		{
			name:     "Empty path",
			url:      "https://example.com",
			expected: filepath.Join(tempDir, "example.com", "index.html"),
		},
		{
			name:     "CSS file",
			url:      "https://example.com/css/style.css",
			expected: filepath.Join(tempDir, "example.com", "css", "style.css"),
		},
		{
			name:     "Image file",
			url:      "https://example.com/images/logo.png",
			expected: filepath.Join(tempDir, "example.com", "images", "logo.png"),
		},
		{
			name:     "Path without extension",
			url:      "https://example.com/about",
			expected: filepath.Join(tempDir, "example.com", "about.html"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testURL, _ := url.Parse(tt.url)
			localPath := storage.GetLocalPath(testURL)
			assert.Equal(t, tt.expected, localPath)
		})
	}
}

func TestFileStorage_GetRelativePath(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewFileStorage(tempDir)
	require.NoError(t, err)

	testFile := filepath.Join(tempDir, "example.com", "test.html")
	err = os.MkdirAll(filepath.Dir(testFile), 0755)
	require.NoError(t, err)

	err = os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	relativePath := storage.GetRelativePath(testFile)
	expectedPath := filepath.Join("example.com", "test.html")
	assert.Equal(t, expectedPath, relativePath)
}

func TestFileStorage_NewFileStorage(t *testing.T) {
	tempDir := filepath.Join(t.TempDir(), "new", "nested", "directory")

	storage, err := NewFileStorage(tempDir)
	require.NoError(t, err)
	assert.NotNil(t, storage)

	assert.DirExists(t, tempDir)
}

func TestFileStorage_Exists(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewFileStorage(tempDir)
	require.NoError(t, err)

	assert.False(t, storage.Exists("/nonexistent/file"))

	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	assert.True(t, storage.Exists(testFile))
}
