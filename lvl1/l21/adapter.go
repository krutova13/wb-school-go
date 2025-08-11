package l21

// LegacyReaderAdapter адаптирует LegacyReader к новому интерфейсу
type LegacyReaderAdapter struct {
	legacy *LegacyReader
}

// ReadData возвращает данные из legacy системы
func (a *LegacyReaderAdapter) ReadData() string {
	return a.legacy.GetLegacyData()
}

// NewLegacyReaderAdapter создает новый адаптер для LegacyReader
func NewLegacyReaderAdapter(legacy *LegacyReader) *LegacyReaderAdapter {
	return &LegacyReaderAdapter{legacy: legacy}
}
