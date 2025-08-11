package telnet

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"telnet/internal/connection"
	"telnet/internal/errorhandler"
	"telnet/internal/io"
	"telnet/internal/logger"
	"telnet/internal/types"
)

// Client реализует telnet-клиент
type Client struct {
	config       *types.Config
	connManager  *connection.Manager
	ioHandler    *io.Handler
	logger       types.Logger
	errorHandler types.ErrorHandler
}

// NewClient создает новый telnet-клиент
func NewClient(config *types.Config) *Client {
	log := logger.NewConsoleLogger()
	errorHandler := errorhandler.NewConsoleErrorHandler()

	// Создаем стандартный net.Dialer и адаптируем его
	netDialer := &net.Dialer{}
	dialer := connection.NewDialerAdapter(netDialer)

	connManager := connection.NewManager(config, dialer, log, errorHandler)

	return &Client{
		config:       config,
		connManager:  connManager,
		logger:       log,
		errorHandler: errorHandler,
	}
}

// Connect устанавливает соединение
func (c *Client) Connect() error {
	return c.connManager.Connect()
}

// Start запускает telnet-клиент
func (c *Client) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	conn := c.connManager.GetConnection()
	if conn == nil {
		return fmt.Errorf("connection is not established")
	}

	c.ioHandler = io.NewHandler(conn, c.logger, c.errorHandler)

	go func() {
		if err := c.ioHandler.Start(ctx); err != nil {
			c.errorHandler.HandleError(fmt.Errorf("IO handler error: %w", err))
		}
	}()

	select {
	case <-sigChan:
		c.logger.Log("Received termination signal, shutting down...")
	case <-ctx.Done():
		c.logger.Log("Context cancelled, shutting down...")
	}

	return nil
}

// Close закрывает клиент
func (c *Client) Close() error {
	if c.ioHandler != nil {
		if err := c.ioHandler.Close(); err != nil {
			c.errorHandler.HandleError(fmt.Errorf("failed to close IO handler: %w", err))
		}
	}

	if err := c.connManager.Close(); err != nil {
		c.errorHandler.HandleError(fmt.Errorf("failed to close connection: %w", err))
		return err
	}

	return nil
}
