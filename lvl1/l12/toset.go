package l12

import "fmt"

// ToSet преобразует слайс строк в множество (map)
func ToSet(words []string) map[string]struct{} {
	set := make(map[string]struct{})

	for _, word := range words {
		set[word] = struct{}{}
	}
	return set
}

// PrintResult выводит множество строк в консоль
func PrintResult(result map[string]struct{}) {
	fmt.Print("Полученное множество: {")
	first := true
	for word := range result {
		if !first {
			fmt.Print(", ")
		}
		fmt.Print(word)
		first = false
	}
	fmt.Println("}")
}
