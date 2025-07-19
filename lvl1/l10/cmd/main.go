package main

import "fmt"

func main() {
	temps := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5, 29.9, -29.9, 0, 2}
	groups := make(map[int][]float64)

	for _, temp := range temps {
		key := int(temp/10) * 10
		groups[key] = append(groups[key], temp)
	}

	for k, v := range groups {
		fmt.Printf("%d : %v\n\n", k, v)
	}
}
