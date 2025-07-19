package main

import (
	"fmt"
	"wbschoolgo/lvl1/l11"
)

func main() {
	sliceA := []int{1, 2, 3}
	sliceB := []int{2, 3, 4}

	fmt.Printf("Пересечение = %v", l11.Intersect(sliceA, sliceB))
}
