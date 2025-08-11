package types

import (
	"context"
	"net/url"
)

// Resource представляет загружаемый ресурс
type Resource struct {
	URL       *url.URL
	Depth     int
	Type      ResourceType
	Content   []byte
	MimeType  string
	LocalPath string
}

// ResourceType определяет тип ресурса
type ResourceType int

const (
	// ResourceTypeHTML представляет HTML документ
	ResourceTypeHTML ResourceType = iota
	// ResourceTypeCSS представляет CSS файл
	ResourceTypeCSS
	// ResourceTypeJS представляет JavaScript файл
	ResourceTypeJS
	// ResourceTypeImage представляет изображение
	ResourceTypeImage
	// ResourceTypeOther представляет другие типы ресурсов
	ResourceTypeOther
)

// String возвращает строковое представление типа ресурса
func (rt ResourceType) String() string {
	switch rt {
	case ResourceTypeHTML:
		return "HTML"
	case ResourceTypeCSS:
		return "CSS"
	case ResourceTypeJS:
		return "JavaScript"
	case ResourceTypeImage:
		return "Image"
	case ResourceTypeOther:
		return "Other"
	default:
		return "Unknown"
	}
}

// Downloader интерфейс для загрузки ресурсов
type Downloader interface {
	Download(ctx context.Context, resourceURL *url.URL) (*Resource, error)
}

// Parser интерфейс для парсинга HTML и извлечения ссылок
type Parser interface {
	ParseHTML(content []byte, baseURL *url.URL) ([]*url.URL, error)
	ParseCSS(content []byte, baseURL *url.URL) ([]*url.URL, error)
}

// Storage интерфейс для сохранения ресурсов
type Storage interface {
	Save(resource *Resource) error
	Exists(localPath string) bool
	GetLocalPath(resourceURL *url.URL) string
}

// URLFilter интерфейс для фильтрации URL
type URLFilter interface {
	ShouldDownload(url *url.URL, depth int) bool
	MarkVisited(url *url.URL)
}

// ProgressReporter интерфейс для отчета о прогрессе
type ProgressReporter interface {
	ReportDownloaded(resource *Resource)
	ReportSkipped(url *url.URL, reason string)
	ReportError(url *url.URL, err error)
	PrintSummary()
}
