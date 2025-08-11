package main

import (
	"fmt"
)

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	result := FindAnagrams(words)
	fmt.Println(result)
}
