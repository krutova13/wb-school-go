package config

import "time"

// Config содержит все настройки для утилиты wget
type Config struct {
	URL           string        // URL для загрузки
	Depth         int           // Глубина рекурсии
	OutputDir     string        // Директория для сохранения
	Concurrency   int           // Количество одновременных загрузок
	Timeout       time.Duration // Таймаут для HTTP запросов
	RespectRobots bool          // Соблюдать robots.txt
	UserAgent     string        // User-Agent для запросов
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.URL == "" {
		return ErrEmptyURL
	}

	if c.Depth < 0 {
		return ErrInvalidDepth
	}

	if c.Concurrency <= 0 {
		return ErrInvalidConcurrency
	}

	if c.Timeout <= 0 {
		return ErrInvalidTimeout
	}

	return nil
}
