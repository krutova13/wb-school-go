package l3

import "fmt"

func Worker(id int, ch <-chan int) {
	for data := range ch {
		fmt.Printf("Воркер %d получил значение: %d\n", id, data)
	}
}
