package main

import (
	"fmt"
	"time"
)

func main() {
	stopCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopCh:
				fmt.Println("Остановлено через канал")
				return
			default:
				fmt.Println("Горутина работает")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	time.Sleep(2 * time.Second)
	close(stopCh)
	time.Sleep(1 * time.Second)
}
