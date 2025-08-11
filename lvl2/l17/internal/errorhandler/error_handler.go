package errorhandler

import (
	"fmt"
	"os"
)

// ConsoleErrorHandler обрабатывает ошибки через консоль
type ConsoleErrorHandler struct{}

// NewConsoleErrorHandler создает новый обработчик ошибок
func NewConsoleErrorHandler() *ConsoleErrorHandler {
	return &ConsoleErrorHandler{}
}

// HandleError обрабатывает ошибку, выводя её в stderr
func (h *ConsoleErrorHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// FatalErrorHandler обрабатывает ошибки с завершением программы
type FatalErrorHandler struct{}

// NewFatalErrorHandler создает новый фатальный обработчик ошибок
func NewFatalErrorHandler() *FatalErrorHandler {
	return &FatalErrorHandler{}
}

// HandleError обрабатывает ошибку и завершает программу
func (h *FatalErrorHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
	os.Exit(1)
}

// SilentErrorHandler тихий обработчик ошибок (для тестов)
type SilentErrorHandler struct{}

// NewSilentErrorHandler создает новый тихий обработчик ошибок
func NewSilentErrorHandler() *SilentErrorHandler {
	return &SilentErrorHandler{}
}

// HandleError ничего не делает
func (h *SilentErrorHandler) HandleError(err error) {
	// Тихий обработчик не обрабатывает ошибки
} 