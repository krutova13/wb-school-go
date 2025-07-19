package main

import (
	"fmt"
	"wbschoolgo/lvl1/l13"
)

func main() {
	a := 5
	b := 7
	a, b = l13.Exchange(a, b)
	fmt.Printf("Результат: %d, %d", a, b)
}
