package logger

import (
	"log"
	"os"
)

// ConsoleLogger реализует логирование в консоль
type ConsoleLogger struct {
	logger *log.Logger
}

// NewConsoleLogger создает новый консольный логгер
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{
		logger: log.New(os.Stderr, "[TELNET] ", log.LstdFlags),
	}
}

// Log записывает сообщение в лог
func (l *ConsoleLogger) Log(message string) {
	l.logger.Println(message)
}

// LogError записывает ошибку в лог
func (l *ConsoleLogger) LogError(err error) {
	l.logger.Printf("ERROR: %v", err)
}
