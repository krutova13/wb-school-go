package main

import (
	"fmt"
	"wbschoolgo/lvl1/l17"
)

func main() {
	arr := []int{1, 4, 6, 7, 9, 12, 16, 19}
	target := 5

	result := l17.BinarySearch(arr, target)

	if result == -1 {
		fmt.Printf("Элемент %d не найден", target)
	} else {
		fmt.Printf("Элемент %d найден", target)
	}
}
