package main

import (
	"fmt"
	"wbschoolgo/lvl1/l16"
)

func main() {
	arr := []int{3, 5, 1, 8, 4, 0, 12, 2, 5, 8}
	result := l16.QuickSort(arr)
	fmt.Println("Отсортированный срез:", result)
}
