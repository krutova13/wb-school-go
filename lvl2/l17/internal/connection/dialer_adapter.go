package connection

import (
	"context"
	"net"
	"time"

	"telnet/internal/types"
)

// DialerAdapter адаптирует net.Dialer к интерфейсу types.Dialer
type DialerAdapter struct {
	dialer *net.Dialer
}

// NewDialerAdapter создает новый адаптер для net.Dialer
func NewDialerAdapter(dialer *net.Dialer) *DialerAdapter {
	return &DialerAdapter{
		dialer: dialer,
	}
}

// Dial реализует метод Dial интерфейса types.Dialer
func (d *DialerAdapter) Dial(network, address string) (types.Connection, error) {
	conn, err := d.dialer.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &Adapter{conn: conn}, nil
}

// DialContext реализует метод DialContext интерфейса types.Dialer
func (d *DialerAdapter) DialContext(ctx context.Context, network, address string) (types.Connection, error) {
	conn, err := d.dialer.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	return &Adapter{conn: conn}, nil
}

// Adapter адаптирует net.Conn к интерфейсу types.Connection
type Adapter struct {
	conn net.Conn
}

// Read реализует io.Reader
func (c *Adapter) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

// Write реализует io.Writer
func (c *Adapter) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

// Close реализует io.Closer
func (c *Adapter) Close() error {
	return c.conn.Close()
}

// SetDeadline устанавливает дедлайн для соединения
func (c *Adapter) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}
