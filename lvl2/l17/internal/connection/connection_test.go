package connection

import (
	"context"
	"testing"
	"time"

	"telnet/internal/errorhandler"
	"telnet/internal/logger"
	"telnet/internal/types"

	"github.com/stretchr/testify/assert"
)

// MockDialer реализует мок-диалер для тестов
type MockDialer struct {
	shouldFail bool
}

func (m *MockDialer) Dial(network, address string) (types.Connection, error) {
	_ = network
	_ = address
	if m.shouldFail {
		return nil, assert.AnError
	}
	return &MockConnection{}, nil
}

func (m *MockDialer) DialContext(ctx context.Context, network, address string) (types.Connection, error) {
	_ = ctx
	return m.Dial(network, address)
}

// MockConnection реализует мок-соединение для тестов
type MockConnection struct {
	closed bool
}

func (m *MockConnection) Read(p []byte) (n int, err error) {
	_ = p
	return 0, nil
}

func (m *MockConnection) Write(p []byte) (n int, err error) {
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

func TestNewManager(t *testing.T) {
	config := &types.Config{
		Host:    "localhost",
		Port:    "8080",
		Timeout: 10 * time.Second,
	}
	dialer := &MockDialer{}
	log := logger.NewConsoleLogger()
	errorHandler := errorhandler.NewSilentErrorHandler()

	manager := NewManager(config, dialer, log, errorHandler)

	assert.NotNil(t, manager)
	assert.Equal(t, config, manager.config)
	assert.Equal(t, dialer, manager.dialer)
	assert.Equal(t, log, manager.logger)
	assert.Equal(t, errorHandler, manager.errorHandler)
}

func TestManager_Connect_Success(t *testing.T) {
	config := &types.Config{
		Host:    "localhost",
		Port:    "8080",
		Timeout: 10 * time.Second,
	}
	dialer := &MockDialer{shouldFail: false}
	log := logger.NewConsoleLogger()
	errorHandler := errorhandler.NewSilentErrorHandler()

	manager := NewManager(config, dialer, log, errorHandler)

	err := manager.Connect()
	assert.NoError(t, err)
	assert.NotNil(t, manager.conn)
}

func TestManager_Connect_Failure(t *testing.T) {
	config := &types.Config{
		Host:    "localhost",
		Port:    "8080",
		Timeout: 10 * time.Second,
	}
	dialer := &MockDialer{shouldFail: true}
	log := logger.NewConsoleLogger()
	errorHandler := errorhandler.NewSilentErrorHandler()

	manager := NewManager(config, dialer, log, errorHandler)

	err := manager.Connect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
}

func TestManager_Close(t *testing.T) {
	config := &types.Config{
		Host:    "localhost",
		Port:    "8080",
		Timeout: 10 * time.Second,
	}
	dialer := &MockDialer{shouldFail: false}
	log := logger.NewConsoleLogger()
	errorHandler := errorhandler.NewSilentErrorHandler()

	manager := NewManager(config, dialer, log, errorHandler)

	// Сначала устанавливаем соединение
	err := manager.Connect()
	assert.NoError(t, err)

	// Затем закрываем
	err = manager.Close()
	assert.NoError(t, err)
}

func TestManager_Close_WithoutConnection(t *testing.T) {
	config := &types.Config{
		Host:    "localhost",
		Port:    "8080",
		Timeout: 10 * time.Second,
	}
	dialer := &MockDialer{shouldFail: false}
	log := logger.NewConsoleLogger()
	errorHandler := errorhandler.NewSilentErrorHandler()

	manager := NewManager(config, dialer, log, errorHandler)

	err := manager.Close()
	assert.NoError(t, err)
}
