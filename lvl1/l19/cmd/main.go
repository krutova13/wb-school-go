package main

import (
	"fmt"
	"wbschoolgo/lvl1/l19"
)

func main() {
	word := "Проверка"
	reversedWord := l19.ReverseString(word)
	fmt.Printf("Перевернутая строка: %v", reversedWord)
}
