package l21

// LegacyReader представляет устаревшую систему чтения данных
type LegacyReader struct{}

// GetLegacyData возвращает данные из legacy системы
func (l *LegacyReader) GetLegacyData() string {
	return "Данные из legacyReader"
}
