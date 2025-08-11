package cut

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FieldRange представляет диапазон полей
type FieldRange struct {
	Start int
	End   int
}

// Processor обрабатывает данные согласно параметрам cut
type Processor struct {
	fields    []FieldRange
	delimiter string
	separated bool
}

// NewProcessor создает новый процессор cut
func NewProcessor(fieldsStr, delimiter string, separated bool) (*Processor, error) {
	fields, err := parseFields(fieldsStr)
	if err != nil {
		return nil, err
	}

	return &Processor{
		fields:    fields,
		delimiter: delimiter,
		separated: separated,
	}, nil
}

// Process обрабатывает данные из stdin
func (cp *Processor) Process() error {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		result, shouldOutput := cp.ProcessLine(line)
		if shouldOutput {
			fmt.Println(result)
		}
	}

	return scanner.Err()
}

// ProcessLine обрабатывает одну строку
func (cp *Processor) ProcessLine(line string) (string, bool) {
	if cp.separated && !strings.Contains(line, cp.delimiter) {
		return "", false
	}

	fields := strings.Split(line, cp.delimiter)

	var result []string
	for _, fieldRange := range cp.fields {
		for i := fieldRange.Start; i <= fieldRange.End; i++ {
			fieldIndex := i - 1
			if fieldIndex >= 0 && fieldIndex < len(fields) {
				result = append(result, fields[fieldIndex])
			}
		}
	}

	return strings.Join(result, cp.delimiter), true
}

// parseFields парсит строку полей в список диапазонов
func parseFields(fieldsStr string) ([]FieldRange, error) {
	if fieldsStr == "" {
		return nil, fmt.Errorf("пустая строка полей")
	}

	var ranges []FieldRange
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("неверный формат диапазона: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("неверный номер поля: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("неверный номер поля: %s", rangeParts[1])
			}

			if start <= 0 || end <= 0 {
				return nil, fmt.Errorf("номера полей должны быть положительными")
			}

			if start > end {
				return nil, fmt.Errorf("начало диапазона должно быть меньше или равно концу: %s", part)
			}

			ranges = append(ranges, FieldRange{Start: start, End: end})
		} else {
			fieldNum, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("неверный номер поля: %s", part)
			}

			if fieldNum <= 0 {
				return nil, fmt.Errorf("номер поля должен быть положительным: %d", fieldNum)
			}

			ranges = append(ranges, FieldRange{Start: fieldNum, End: fieldNum})
		}
	}

	return ranges, nil
}
