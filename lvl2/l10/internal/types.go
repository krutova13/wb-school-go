package internal

import (
	"sort"
)

// SortOptions содержит опции для сортировки
type SortOptions struct {
	Column         int  // -k N
	Numeric        bool // -n
	Reverse        bool // -r
	Unique         bool // -u
	MonthSort      bool // -M
	TrailingBlanks bool // -b
	CheckSorted    bool // -c
	HumanReadable  bool // -h
}

// Line представляет строку для сортировки
type Line struct {
	Original string
	Fields   []string
	Key      string
}

// Глобальные опции сортировки
var opts SortOptions

// GetOpts возвращает указатель на глобальные опции сортировки
func GetOpts() *SortOptions {
	return &opts
}

// SortableLines представляет срез строк для сортировки
type SortableLines []Line

// Len возвращает длину среза
func (sl SortableLines) Len() int {
	return len(sl)
}

// Swap меняет местами элементы
func (sl SortableLines) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

// Less определяет порядок сортировки
func (sl SortableLines) Less(i, j int) bool {
	line1, line2 := sl[i], sl[j]

	key1 := line1.Key
	key2 := line2.Key

	opts := GetOpts()
	if opts.Numeric {
		return CompareNumeric(key1, key2)
	} else if opts.MonthSort {
		return CompareMonths(key1, key2)
	} else if opts.HumanReadable {
		return CompareHumanReadable(key1, key2)
	} else {
		return key1 < key2
	}
}

// Sort выполняет сортировку с учетом опций
func (sl SortableLines) Sort() {
	opts := GetOpts()
	if opts.Reverse {
		sort.Sort(sort.Reverse(sl))
	} else {
		sort.Sort(sl)
	}
}

// IsSorted проверяет, отсортированы ли данные
func (sl SortableLines) IsSorted() bool {
	for i := 1; i < len(sl); i++ {
		if !sl.Less(i-1, i) {
			return false
		}
	}
	return true
}
