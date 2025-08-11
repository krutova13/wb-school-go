package internal

// GrepOptions содержит опции для фильтрации
type GrepOptions struct {
	After  int  // -A N
	Before int  // -B N
	Circle int  // -C N
	Count  bool // -c
	Ignore bool // -i
	Invert bool // -v
	Fix    bool // -F
	Number bool // -n
}

// Глобальные опции фильтрации
var opts GrepOptions

// GetOpts возвращает указатель на глобальные опции фильтрации
func GetOpts() *GrepOptions {
	return &opts
}

// LineInfo содержит информацию о строке
type LineInfo struct {
	Number int    // Номер строки
	Text   string // Текст строки
	Match  bool   // Флаг совпадения с паттерном
}
