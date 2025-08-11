package sort

import (
	"sortutil/internal"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func withOptions(options internal.SortOptions, testFunc func()) {
	originalOpts := *internal.GetOpts()

	defer func() {
		opts := internal.GetOpts()
		*opts = originalOpts
	}()

	opts := internal.GetOpts()
	*opts = options

	testFunc()
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		opts     internal.SortOptions
		expected []string
	}{
		{
			name:     "базовая сортировка",
			input:    []string{"zebra", "apple", "banana", "cherry"},
			opts:     internal.SortOptions{},
			expected: []string{"apple", "banana", "cherry", "zebra"},
		},
		{
			name:     "числовая сортировка",
			input:    []string{"10", "2", "1", "20"},
			opts:     internal.SortOptions{Numeric: true},
			expected: []string{"1", "2", "10", "20"},
		},
		{
			name:     "обратная сортировка",
			input:    []string{"apple", "banana", "cherry"},
			opts:     internal.SortOptions{Reverse: true},
			expected: []string{"cherry", "banana", "apple"},
		},
		{
			name:     "удаление дубликатов",
			input:    []string{"apple", "banana", "apple", "cherry", "banana"},
			opts:     internal.SortOptions{Unique: true},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "сортировка по столбцу",
			input:    []string{"zebra\t5", "apple\t3", "banana\t1", "cherry\t2"},
			opts:     internal.SortOptions{Column: 2, Numeric: true},
			expected: []string{"banana\t1", "cherry\t2", "apple\t3", "zebra\t5"},
		},
		{
			name:     "сортировка по месяцам",
			input:    []string{"Dec", "Jan", "Mar", "Feb"},
			opts:     internal.SortOptions{MonthSort: true},
			expected: []string{"Jan", "Feb", "Mar", "Dec"},
		},
		{
			name:     "сортировка по размерам",
			input:    []string{"1KB", "2MB", "500B", "1GB"},
			opts:     internal.SortOptions{HumanReadable: true},
			expected: []string{"500B", "1KB", "2MB", "1GB"},
		},
		{
			name:     "игнорирование пробелов",
			input:    []string{"apple   ", "banana", "cherry  "},
			opts:     internal.SortOptions{TrailingBlanks: true},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "числовая сортировка с обратным порядком",
			input:    []string{"1", "10", "2", "20"},
			opts:     internal.SortOptions{Numeric: true, Reverse: true},
			expected: []string{"20", "10", "2", "1"},
		},
		{
			name:     "сортировка с удалением дубликатов и обратным порядком",
			input:    []string{"apple", "banana", "apple", "cherry", "banana"},
			opts:     internal.SortOptions{Unique: true, Reverse: true},
			expected: []string{"cherry", "banana", "apple"},
		},
		{
			name:     "сортировка по столбцу с обратным порядком",
			input:    []string{"zebra\t5", "apple\t3", "banana\t1", "cherry\t2"},
			opts:     internal.SortOptions{Column: 2, Numeric: true, Reverse: true},
			expected: []string{"zebra\t5", "apple\t3", "cherry\t2", "banana\t1"},
		},
		{
			name:     "сортировка по месяцам с обратным порядком",
			input:    []string{"Dec", "Jan", "Mar", "Feb"},
			opts:     internal.SortOptions{MonthSort: true, Reverse: true},
			expected: []string{"Dec", "Mar", "Feb", "Jan"},
		},
		{
			name:     "сортировка по размерам с обратным порядком",
			input:    []string{"1KB", "2MB", "500B", "1GB"},
			opts:     internal.SortOptions{HumanReadable: true, Reverse: true},
			expected: []string{"1GB", "2MB", "1KB", "500B"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withOptions(tt.opts, func() {
				var input []string
				if tt.opts.TrailingBlanks {
					for _, line := range tt.input {
						input = append(input, strings.TrimRight(line, " \t"))
					}
				} else {
					input = tt.input
				}

				sortableLines := internal.PrepareLines(input)
				sortableLines.Sort()

				result := make([]string, len(sortableLines))
				for i, line := range sortableLines {
					result[i] = line.Original
				}

				assert.Equal(t, tt.expected, result, "Сортировка должна работать корректно")
			})
		})
	}
}

func TestCheckSorted(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		opts     internal.SortOptions
		expected bool
	}{
		{
			name:     "отсортированный массив",
			input:    []string{"apple", "banana", "cherry"},
			opts:     internal.SortOptions{},
			expected: true,
		},
		{
			name:     "неотсортированный массив",
			input:    []string{"cherry", "apple", "banana"},
			opts:     internal.SortOptions{},
			expected: false,
		},
		{
			name:     "отсортированный массив в обратном порядке",
			input:    []string{"cherry", "banana", "apple"},
			opts:     internal.SortOptions{Reverse: true},
			expected: true,
		},
		{
			name:     "неотсортированный массив в обратном порядке",
			input:    []string{"apple", "cherry", "banana"},
			opts:     internal.SortOptions{Reverse: true},
			expected: false,
		},
		{
			name:     "отсортированный числовой массив",
			input:    []string{"1", "2", "10", "20"},
			opts:     internal.SortOptions{Numeric: true},
			expected: true,
		},
		{
			name:     "неотсортированный числовой массив",
			input:    []string{"10", "1", "20", "2"},
			opts:     internal.SortOptions{Numeric: true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withOptions(tt.opts, func() {
				sortableLines := internal.PrepareLines(tt.input)
				result := sortableLines.IsSorted()

				assert.Equal(t, tt.expected, result, "Проверка отсортированности должна работать корректно")
			})
		})
	}
}
