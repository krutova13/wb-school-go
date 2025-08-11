package io

import (
	"context"
	"testing"
	"time"

	"telnet/internal/errorhandler"
	"telnet/internal/logger"

	"github.com/stretchr/testify/assert"
)

// MockConnection реализует мок-соединение для тестов
type MockConnection struct {
	readData  []byte
	writeData []byte
	closed    bool
	readErr   error
	writeErr  error
}

func (m *MockConnection) Read(p []byte) (n int, err error) {
	if m.readErr != nil {
		return 0, m.readErr
	}
	if len(m.readData) == 0 {
		return 0, nil
	}
	copy(p, m.readData)
	return len(m.readData), nil
}

func (m *MockConnection) Write(p []byte) (n int, err error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	m.writeData = append(m.writeData, p...)
	return len(p), nil
}

func (m *MockConnection) Close() error {
	m.closed = true
	return nil
}

func (m *MockConnection) SetDeadline(t time.Time) error {
	_ = t
	return nil
}

func TestNewHandler(t *testing.T) {
	mockConn := &MockConnection{}
	log := logger.NewConsoleLogger()
	errorHandler := errorhandler.NewSilentErrorHandler()

	handler := NewHandler(mockConn, log, errorHandler)

	assert.NotNil(t, handler)
	assert.Equal(t, mockConn, handler.conn)
	assert.Equal(t, log, handler.logger)
	assert.Equal(t, errorHandler, handler.errorHandler)
}

func TestHandler_Start_WithoutConnection(t *testing.T) {
	log := logger.NewConsoleLogger()
	errorHandler := errorhandler.NewSilentErrorHandler()

	handler := NewHandler(nil, log, errorHandler)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := handler.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection is not established")
}

func TestHandler_Close(t *testing.T) {
	mockConn := &MockConnection{}
	log := logger.NewConsoleLogger()
	errorHandler := errorhandler.NewSilentErrorHandler()

	handler := NewHandler(mockConn, log, errorHandler)

	err := handler.Close()
	assert.NoError(t, err)
	assert.True(t, handler.closed)
}
