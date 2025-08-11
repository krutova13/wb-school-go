package l3

import "fmt"

// Worker обрабатывает данные из канала
func Worker(id int, ch <-chan int) {
	for data := range ch {
		fmt.Printf("Воркер %d получил значение: %d\n", id, data)
	}
}
