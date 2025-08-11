package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// UnpackString распаковывает строку, содержащую сжатые символы
// Поддерживает следующие форматы:
// - "a3b2c" -> "aaabbc" (повторение символов)
// - "a\\3b" -> "a3b" (escape-последовательности)
//
// # Возвращает распакованную строку и ошибку, если строка некорректна
//
// Ошибки возникают если:
// - строка начинается с цифры
// - некорректная escape-последовательность
// - некорректное число в строке
func UnpackString(s string) (string, error) {
	chars := []rune(s)
	var result []rune

	if len(s) == 0 {
		return s, nil
	}

	if err := validateFirstChar(chars); err != nil {
		return "", err
	}

	for i := 0; i < len(chars); i++ {
		escaping := '\\'

		if chars[i] == escaping {
			nextChar, newIndex, err := handleEscapeSequence(chars, i)

			if err != nil {
				return "", err
			}

			result = append(result, nextChar)
			i = newIndex
			continue
		}

		if !unicode.IsDigit(chars[i]) {
			result = append(result, chars[i])
			continue
		}

		number, newIndex, err := extractNumber(chars, i)

		if err != nil {
			return "", fmt.Errorf("некорректное число в строке: %s\n Детали: %s", s, err)
		}

		prevChar := chars[i-1]
		result = appendRepeatedChar(result, prevChar, number)
		i = newIndex
	}
	return string(result), nil
}

func validateFirstChar(chars []rune) error {
	if unicode.IsDigit(chars[0]) {
		return fmt.Errorf("строка начинается с цифры %c", chars[0])
	}
	return nil
}

func handleEscapeSequence(chars []rune, currentIndex int) (rune, int, error) {
	if currentIndex+1 >= len(chars) {
		return 0, 0, fmt.Errorf("некорректная escape-последовательность в конце строки")
	}
	return chars[currentIndex+1], currentIndex + 1, nil
}

func extractNumber(chars []rune, startIndex int) (int, int, error) {
	var numberString strings.Builder
	j := startIndex
	for j < len(chars) && unicode.IsDigit(chars[j]) {
		numberString.WriteRune(chars[j])
		j++
	}

	number, err := builderToInt(&numberString)
	if err != nil {
		return 0, 0, err
	}

	return number, j - 1, nil
}

func appendRepeatedChar(result []rune, char rune, count int) []rune {
	for k := 1; k < count; k++ {
		result = append(result, char)
	}
	return result
}

func builderToInt(builder *strings.Builder) (int, error) {
	return strconv.Atoi(builder.String())
}
