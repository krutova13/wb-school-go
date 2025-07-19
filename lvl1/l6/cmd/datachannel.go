package main

import (
	"fmt"
	"time"
)

func main() {
	dataCh := make(chan int)
	go func() {
		for val := range dataCh {
			fmt.Println("Получено:", val)
		}
		fmt.Println("Остановлено по закрытию канала")
	}()
	for i := 0; i < 3; i++ {
		dataCh <- i
		time.Sleep(500 * time.Millisecond)
	}
	close(dataCh)
	time.Sleep(1 * time.Second)
}
