package l18

import "sync"

// Counter представляет потокобезопасный счетчик
type Counter struct {
	value int
	mu    sync.Mutex
}

// Inc увеличивает значение счетчика на 1
func (c *Counter) Inc() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

// Value возвращает текущее значение счетчика
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}
