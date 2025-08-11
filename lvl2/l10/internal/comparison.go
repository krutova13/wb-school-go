package internal

import (
	"regexp"
	"strconv"
	"strings"
)

// Константы для размеров файлов
const (
	Byte = 1
	KB   = 1024
	MB   = 1024 * KB
	GB   = 1024 * MB
	TB   = 1024 * GB
)

// Константы для регулярных выражений
const (
	// ExpectedGroups количество групп в регулярном выражении для парсинга размеров
	// Группа 0: полное совпадение
	// Группа 1: число
	// Группа 2: суффикс (KB, MB, GB, TB)
	ExpectedGroups = 3

	// NumberGroup Индексы групп в регулярном выражении
	NumberGroup = 1
	SuffixGroup = 2
)

// CompareNumeric сравнивает строки как числа, если это возможно.
// Если обе строки являются числами, сравнивает их числовые значения.
// Если только одна строка является числом, она считается меньше.
// В остальных случаях сравнивает как строки.
func CompareNumeric(a, b string) bool {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)

	numA, errA := strconv.ParseFloat(a, 64)
	numB, errB := strconv.ParseFloat(b, 64)

	if errA == nil && errB == nil {
		return numA < numB
	}

	if errA == nil {
		return true
	}
	if errB == nil {
		return false
	}

	return a < b
}

// CompareMonths сравнивает строки как названия месяцев.
// Поддерживает сокращенные названия месяцев (jan, feb, mar и т.д.).
// Если обе строки являются месяцами, сравнивает их порядок.
// Если только одна строка является месяцем, она считается меньше.
// В остальных случаях сравнивает как строки.
func CompareMonths(a, b string) bool {
	monthOrder := map[string]int{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4,
		"may": 5, "jun": 6, "jul": 7, "aug": 8,
		"sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}

	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))

	monthA, isMonthA := monthOrder[a]
	monthB, isMonthB := monthOrder[b]

	if isMonthA && isMonthB {
		return monthA < monthB
	}

	if isMonthA {
		return true // месяцы идут перед не-месяцами
	}
	if isMonthB {
		return false
	}

	return a < b
}

// CompareHumanReadable сравнивает строки как размеры в человекочитаемом формате.
// Поддерживает суффиксы K, M, G, T (KB, MB, GB, TB).
// Если обе строки являются размерами, сравнивает их числовые значения.
// Если только одна строка является размером, она считается меньше.
// В остальных случаях сравнивает как строки.
func CompareHumanReadable(a, b string) bool {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)

	sizeA, errA := parseHumanSize(a)
	sizeB, errB := parseHumanSize(b)

	if errA == nil && errB == nil {
		return sizeA < sizeB
	}

	if errA == nil {
		return true // размеры идут перед не-размерами
	}
	if errB == nil {
		return false
	}

	return a < b
}

// parseHumanSize парсит строку с размером в человекочитаемом формате
func parseHumanSize(s string) (float64, error) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([KMGT]?[B]?)$`)
	matches := re.FindStringSubmatch(strings.ToUpper(s))

	if len(matches) != ExpectedGroups {
		return 0, strconv.ErrSyntax
	}

	number, err := strconv.ParseFloat(matches[NumberGroup], 64)
	if err != nil {
		return 0, err
	}

	multipliers := map[string]float64{
		"":   Byte,
		"B":  Byte,
		"K":  KB,
		"KB": KB,
		"M":  MB,
		"MB": MB,
		"G":  GB,
		"GB": GB,
		"T":  TB,
		"TB": TB,
	}

	suffix := matches[SuffixGroup]
	multiplier, exists := multipliers[suffix]
	if !exists {
		return 0, strconv.ErrSyntax
	}

	return number * multiplier, nil
}
