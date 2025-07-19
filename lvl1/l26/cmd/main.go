package main

import (
	"fmt"
	"wbschoolgo/lvl1/l26"
)

func main() {
	var s string
	fmt.Print("Введите строку: ")
	_, err := fmt.Scan(&s)
	if err != nil {
		return
	}
	if l26.HasUniqueChars(s) {
		fmt.Printf("Строка %v содержит только уникальные символы", s)
	} else {
		fmt.Printf("Строка %v содержит дубликаты", s)
	}
}
