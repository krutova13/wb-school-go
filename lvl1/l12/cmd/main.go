package main

import (
	"wbschoolgo/lvl1/l12"
)

func main() {
	strings := []string{"cat", "cat", "dog", "cat", "tree"}
	result := l12.ToSet(strings)
	l12.PrintResult(result)
}
