package main

import (
	"fmt"
	"wbschoolgo/lvl1/l8"
)

func main() {
	var n int64 = 5
	n = l8.SetBit(n, 1, 0)
	fmt.Println(n)
}
