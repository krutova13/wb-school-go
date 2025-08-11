package types

import (
	"context"
	"io"
	"time"
)

// Config представляет конфигурацию telnet-клиента
type Config struct {
	Host    string
	Port    string
	Timeout time.Duration
}

// Connection представляет TCP-соединение
type Connection interface {
	io.ReadWriteCloser
	SetDeadline(t time.Time) error
}

// Dialer интерфейс для установки соединения
type Dialer interface {
	Dial(network, address string) (Connection, error)
	DialContext(ctx context.Context, network, address string) (Connection, error)
}

// Reader интерфейс для чтения данных
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer интерфейс для записи данных
type Writer interface {
	Write(p []byte) (n int, err error)
}

// ReadWriteCloser объединяет Reader, Writer и Closer
type ReadWriteCloser interface {
	Reader
	Writer
	io.Closer
}

// TelnetClient интерфейс для telnet-клиента
type TelnetClient interface {
	Connect() error
	Start() error
	Close() error
}

// ErrorHandler интерфейс для обработки ошибок
type ErrorHandler interface {
	HandleError(err error)
}

// Logger интерфейс для логирования
type Logger interface {
	Log(message string)
	LogError(err error)
}
