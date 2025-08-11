package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger представляет middleware для логирования
type Logger struct {
	logger *log.Logger
}

// NewLogger создает новый экземпляр логгера
func NewLogger() *Logger {
	return &Logger{
		logger: log.New(log.Writer(), "", log.LstdFlags),
	}
}

// LoggingMiddleware возвращает middleware для логирования HTTP-запросов
func (l *Logger) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l.logger.Printf("Запрос: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		l.logger.Printf("Запрос завершен: %s %s - %v", r.Method, r.URL.Path, duration)
	})
}
