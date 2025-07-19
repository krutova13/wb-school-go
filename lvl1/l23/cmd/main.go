package main

import (
	"fmt"
	"wbschoolgo/lvl1/l23"
)

func main() {
	s := []int{1, 2, 3, 4, 5}
	index := 6
	fmt.Println("Начальное состояние:", s)
	result, err := l23.RemoveElementByIndex(s, index)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("После удаления элемента с индексом %d: %d", index, result)
	}
}
