package filter

import (
	"net/url"
	"strings"
	"sync"

	"wget/internal/config"
)

// URLFilter реализует интерфейс URLFilter для фильтрации URL
type URLFilter struct {
	config     *config.Config
	baseURL    *url.URL
	visited    map[string]bool
	visitedMux sync.RWMutex
}

// NewURLFilter создает новый фильтр URL
func NewURLFilter(cfg *config.Config) (*URLFilter, error) {
	baseURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	return &URLFilter{
		config:  cfg,
		baseURL: baseURL,
		visited: make(map[string]bool),
	}, nil
}

// ShouldDownload определяет нужно ли загружать URL
func (f *URLFilter) ShouldDownload(url *url.URL, depth int) bool {
	if depth > f.config.Depth {
		return false
	}

	f.visitedMux.RLock()
	if f.visited[url.String()] {
		f.visitedMux.RUnlock()
		return false
	}
	f.visitedMux.RUnlock()

	if !f.isSameDomain(url) {
		return false
	}

	if f.config.RespectRobots && !f.isAllowedByRobots(url) {
		return false
	}

	return true
}

// MarkVisited помечает URL как посещенный
func (f *URLFilter) MarkVisited(url *url.URL) {
	f.visitedMux.Lock()
	f.visited[url.String()] = true
	f.visitedMux.Unlock()
}

func (f *URLFilter) isSameDomain(url *url.URL) bool {
	return url.Host == f.baseURL.Host
}

func (f *URLFilter) isAllowedByRobots(url *url.URL) bool {
	if strings.Contains(url.Path, "robots.txt") {
		return false
	}

	adminPaths := []string{"/admin", "/wp-admin", "/administrator", "/manage"}
	for _, path := range adminPaths {
		if strings.HasPrefix(url.Path, path) {
			return false
		}
	}

	return true
}

// GetVisitedCount возвращает количество посещенных URL
func (f *URLFilter) GetVisitedCount() int {
	f.visitedMux.RLock()
	defer f.visitedMux.RUnlock()
	return len(f.visited)
}
