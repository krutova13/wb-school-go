package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "пустая строка",
			input:    "",
			expected: "",
			hasError: false,
		},
		{
			name:     "строка без цифр",
			input:    "abc",
			expected: "abc",
			hasError: false,
		},
		{
			name:     "строка с повторениями",
			input:    "a3b2c",
			expected: "aaabbc",
			hasError: false,
		},
		{
			name:     "одиночные символы",
			input:    "a1b1c1",
			expected: "abc",
			hasError: false,
		},
		{
			name:     "escape-последовательности",
			input:    "a\\3b\\2c",
			expected: "a3b2c",
			hasError: false,
		},
		{
			name:     "escape с цифрами",
			input:    "a\\12b",
			expected: "a11b",
			hasError: false,
		},
		{
			name:     "смешанный тест",
			input:    "a10b3c2ac2d5e\\45",
			expected: "aaaaaaaaaabbbccaccddddde44444",
			hasError: false,
		},
		{
			name:     "начинается с цифры",
			input:    "3abc",
			expected: "",
			hasError: true,
		},
		{
			name:     "некорректная escape-последовательность",
			input:    "abc\\",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnpackString(tt.input)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateFirstChar(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		hasError bool
	}{
		{
			name:     "начинается с буквы",
			input:    []rune("abc"),
			hasError: false,
		},
		{
			name:     "начинается с цифры",
			input:    []rune("3abc"),
			hasError: true,
		},
		{
			name:     "начинается с символа",
			input:    []rune("!abc"),
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFirstChar(tt.input)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleEscapeSequence(t *testing.T) {
	tests := []struct {
		name          string
		chars         []rune
		currentIndex  int
		expectedChar  rune
		expectedIndex int
		hasError      bool
	}{
		{
			name:          "корректная escape-последовательность",
			chars:         []rune("a\\3b"),
			currentIndex:  1,
			expectedChar:  '3',
			expectedIndex: 2,
			hasError:      false,
		},
		{
			name:          "escape в конце строки",
			chars:         []rune("abc\\"),
			currentIndex:  3,
			expectedChar:  0,
			expectedIndex: 0,
			hasError:      true,
		},
		{
			name:          "escape с буквой",
			chars:         []rune("a\\bb"),
			currentIndex:  1,
			expectedChar:  'b',
			expectedIndex: 2,
			hasError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char, newIndex, err := handleEscapeSequence(tt.chars, tt.currentIndex)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedChar, char)
				assert.Equal(t, tt.expectedIndex, newIndex)
			}
		})
	}
}

func TestExtractNumber(t *testing.T) {
	tests := []struct {
		name          string
		chars         []rune
		startIndex    int
		expectedNum   int
		expectedIndex int
		hasError      bool
	}{
		{
			name:          "однозначное число",
			chars:         []rune("a5b"),
			startIndex:    1,
			expectedNum:   5,
			expectedIndex: 1,
			hasError:      false,
		},
		{
			name:          "многозначное число",
			chars:         []rune("a120b"),
			startIndex:    1,
			expectedNum:   120,
			expectedIndex: 3,
			hasError:      false,
		},
		{
			name:          "число в конце",
			chars:         []rune("a3"),
			startIndex:    1,
			expectedNum:   3,
			expectedIndex: 1,
			hasError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num, newIndex, err := extractNumber(tt.chars, tt.startIndex)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNum, num)
				assert.Equal(t, tt.expectedIndex, newIndex)
			}
		})
	}
}

func TestAppendRepeatedChar(t *testing.T) {
	tests := []struct {
		name     string
		result   []rune
		char     rune
		count    int
		expected []rune
	}{
		{
			name:     "повторение 3 раза",
			result:   []rune("abc"),
			char:     'x',
			count:    3,
			expected: []rune("abcxx"),
		},
		{
			name:     "повторение 1 раз",
			result:   []rune("abc"),
			char:     'x',
			count:    1,
			expected: []rune("abc"),
		},
		{
			name:     "пустой результат",
			result:   []rune{},
			char:     'a',
			count:    5,
			expected: []rune("aaaa"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appendRepeatedChar(tt.result, tt.char, tt.count)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuilderToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		hasError bool
	}{
		{
			name:     "корректное число",
			input:    "123",
			expected: 123,
			hasError: false,
		},
		{
			name:     "однозначное число",
			input:    "5",
			expected: 5,
			hasError: false,
		},
		{
			name:     "некорректное число",
			input:    "abc",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			builder.WriteString(tt.input)

			result, err := builderToInt(&builder)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
