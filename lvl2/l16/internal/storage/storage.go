package storage

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"wget/internal/types"
)

// FileStorage реализует интерфейс Storage для сохранения файлов на диск
type FileStorage struct {
	baseDir string
}

// NewFileStorage создает новое файловое хранилище
func NewFileStorage(baseDir string) (*FileStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("ошибка создания директории: %w", err)
	}

	return &FileStorage{
		baseDir: baseDir,
	}, nil
}

// Save сохраняет ресурс на диск
func (s *FileStorage) Save(resource *types.Resource) error {
	localPath := s.GetLocalPath(resource.URL)
	resource.LocalPath = localPath

	// Создаем директории если нужно
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("ошибка создания директории %s: %w", dir, err)
	}

	// Сохраняем файл
	if err := os.WriteFile(localPath, resource.Content, 0644); err != nil {
		return fmt.Errorf("ошибка сохранения файла %s: %w", localPath, err)
	}

	return nil
}

// Exists проверяет существование файла
func (s *FileStorage) Exists(localPath string) bool {
	_, err := os.Stat(localPath)
	return err == nil
}

// GetLocalPath возвращает локальный путь для URL
func (s *FileStorage) GetLocalPath(resourceURL *url.URL) string {
	domain := resourceURL.Host
	urlPath := resourceURL.Path

	if urlPath == "" || strings.HasSuffix(urlPath, "/") {
		urlPath = urlPath + "index.html"
	}

	if !strings.Contains(filepath.Base(urlPath), ".") && !strings.HasSuffix(urlPath, "/") {
		urlPath = urlPath + ".html"
	}

	fullPath := filepath.Join(s.baseDir, domain, urlPath)

	fullPath = filepath.Clean(fullPath)

	return fullPath
}

// GetRelativePath возвращает относительный путь от базовой директории
func (s *FileStorage) GetRelativePath(absolutePath string) string {
	relPath, err := filepath.Rel(s.baseDir, absolutePath)
	if err != nil {
		return absolutePath
	}
	return relPath
}
