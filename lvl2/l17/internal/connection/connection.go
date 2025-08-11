package connection

import (
	"context"
	"fmt"
	"time"

	"telnet/internal/types"
)

// Manager управляет TCP-соединением
type Manager struct {
	config       *types.Config
	dialer       types.Dialer
	conn         types.Connection
	logger       types.Logger
	errorHandler types.ErrorHandler
}

// NewManager создает новый менеджер соединений
func NewManager(config *types.Config, dialer types.Dialer, logger types.Logger, errorHandler types.ErrorHandler) *Manager {
	return &Manager{
		config:       config,
		dialer:       dialer,
		logger:       logger,
		errorHandler: errorHandler,
	}
}

// Connect устанавливает TCP-соединение
func (m *Manager) Connect() error {
	address := fmt.Sprintf("%s:%s", m.config.Host, m.config.Port)

	ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
	defer cancel()

	m.logger.Log(fmt.Sprintf("Connecting to %s...", address))

	conn, err := m.dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	m.conn = conn
	m.logger.Log(fmt.Sprintf("Connected to %s", address))

	return nil
}

// GetConnection возвращает текущее соединение
func (m *Manager) GetConnection() types.Connection {
	return m.conn
}

// Close закрывает соединение
func (m *Manager) Close() error {
	if m.conn != nil {
		m.logger.Log("Closing connection...")
		return m.conn.Close()
	}
	return nil
}

// SetDeadline устанавливает дедлайн для соединения
func (m *Manager) SetDeadline(deadline time.Time) error {
	if m.conn != nil {
		return m.conn.SetDeadline(deadline)
	}
	return fmt.Errorf("connection is not established")
}
