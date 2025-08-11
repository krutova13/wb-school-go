package parser

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTMLParser_ParseHTML(t *testing.T) {
	parser := NewHTMLParser()
	baseURL, _ := url.Parse("https://example.com")

	tests := []struct {
		name     string
		html     string
		expected int // количество найденных ссылок
	}{
		{
			name: "Simple links",
			html: `<html>
				<body>
					<a href="/page1">Page 1</a>
					<a href="/page2">Page 2</a>
					<img src="/image.jpg" />
				</body>
			</html>`,
			expected: 3,
		},
		{
			name: "External links",
			html: `<html>
				<body>
					<a href="https://external.com/page">External</a>
					<a href="/internal">Internal</a>
				</body>
			</html>`,
			expected: 2,
		},
		{
			name: "CSS and JS links",
			html: `<html>
				<head>
					<link rel="stylesheet" href="/style.css" />
					<script src="/script.js"></script>
				</head>
				<body>
					<a href="/page">Page</a>
				</body>
			</html>`,
			expected: 3,
		},
		{
			name: "Filtered links",
			html: `<html>
				<body>
					<a href="javascript:void(0)">JS Link</a>
					<a href="mailto:test@example.com">Email</a>
					<a href="tel:+1234567890">Phone</a>
					<a href="/valid">Valid</a>
				</body>
			</html>`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls, err := parser.ParseHTML([]byte(tt.html), baseURL)
			require.NoError(t, err)
			assert.Len(t, urls, tt.expected)
		})
	}
}

func TestHTMLParser_ParseCSS(t *testing.T) {
	parser := NewHTMLParser()
	baseURL, _ := url.Parse("https://example.com")

	tests := []struct {
		name     string
		css      string
		expected int
	}{
		{
			name: "Simple CSS URLs",
			css: `
				body { background-image: url('/bg.jpg'); }
				.logo { background-image: url('/logo.png'); }
			`,
			expected: 2,
		},
		{
			name: "CSS with quotes",
			css: `
				body { background-image: url("/bg.jpg"); }
				.logo { background-image: url('/logo.png'); }
			`,
			expected: 2,
		},
		{
			name: "External CSS URLs",
			css: `
				body { background-image: url('https://cdn.com/image.jpg'); }
				.logo { background-image: url('/local.png'); }
			`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls, err := parser.ParseCSS([]byte(tt.css), baseURL)
			require.NoError(t, err)
			assert.Len(t, urls, tt.expected)
		})
	}
}

func TestHTMLParser_ResolveURL(t *testing.T) {
	parser := NewHTMLParser()
	baseURL, _ := url.Parse("https://example.com/path/")

	tests := []struct {
		name        string
		urlStr      string
		expected    string
		shouldBeNil bool
	}{
		{
			name:     "Absolute URL",
			urlStr:   "https://other.com/page",
			expected: "https://other.com/page",
		},
		{
			name:     "Relative URL",
			urlStr:   "/page",
			expected: "https://example.com/page",
		},
		{
			name:     "Relative path",
			urlStr:   "subpage",
			expected: "https://example.com/path/subpage",
		},
		{
			name:        "JavaScript link",
			urlStr:      "javascript:void(0)",
			shouldBeNil: true,
		},
		{
			name:        "Mailto link",
			urlStr:      "mailto:test@example.com",
			shouldBeNil: true,
		},
		{
			name:        "Empty URL",
			urlStr:      "",
			shouldBeNil: true,
		},
		{
			name:        "Anchor link",
			urlStr:      "#section",
			shouldBeNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.resolveURL(tt.urlStr, baseURL)

			if tt.shouldBeNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, result.String())
			}
		})
	}
}
