package io

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"telnet/internal/types"
)

// Handler управляет вводом-выводом между STDIN/STDOUT и сокетом
type Handler struct {
	conn         types.Connection
	logger       types.Logger
	errorHandler types.ErrorHandler
	mu           sync.Mutex
	closed       bool
}

// NewHandler создает новый обработчик ввода-вывода
func NewHandler(conn types.Connection, logger types.Logger, errorHandler types.ErrorHandler) *Handler {
	return &Handler{
		conn:         conn,
		logger:       logger,
		errorHandler: errorHandler,
	}
}

// Start запускает обработку ввода-вывода в отдельных горутинах
func (h *Handler) Start(ctx context.Context) error {
	if h.conn == nil {
		return fmt.Errorf("connection is not established")
	}

	// Канал для сигнализации о завершении
	done := make(chan struct{})
	defer close(done)

	// Горутина для чтения из STDIN и записи в сокет
	go h.handleInput(ctx, done)

	// Горутина для чтения из сокета и записи в STDOUT
	go h.handleOutput(ctx, done)

	// Ожидаем завершения контекста или сигнала
	<-ctx.Done()

	h.mu.Lock()
	h.closed = true
	h.mu.Unlock()

	return nil
}

// handleInput читает из STDIN и отправляет в сокет
func (h *Handler) handleInput(ctx context.Context, done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		h.mu.Lock()
		if h.closed {
			h.mu.Unlock()
			return
		}
		h.mu.Unlock()

		line := scanner.Text() + "\n"
		_, err := h.conn.Write([]byte(line))
		if err != nil {
			h.errorHandler.HandleError(fmt.Errorf("failed to write to connection: %w", err))
			return
		}
	}

	if err := scanner.Err(); err != nil {
		h.errorHandler.HandleError(fmt.Errorf("error reading from stdin: %w", err))
	}
}

// handleOutput читает из сокета и отправляет в STDOUT
func (h *Handler) handleOutput(ctx context.Context, done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	buffer := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		h.mu.Lock()
		if h.closed {
			h.mu.Unlock()
			return
		}
		h.mu.Unlock()

		n, err := h.conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				h.logger.Log("Server closed connection")
			} else {
				h.errorHandler.HandleError(fmt.Errorf("error reading from connection: %w", err))
			}
			return
		}

		if n > 0 {
			_, writeErr := os.Stdout.Write(buffer[:n])
			if writeErr != nil {
				h.errorHandler.HandleError(fmt.Errorf("error writing to stdout: %w", writeErr))
				return
			}
		}
	}
}

// Close закрывает обработчик
func (h *Handler) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.closed = true
	return nil
}
