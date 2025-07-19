package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.After(5 * time.Second)

	go func() {
		for {
			select {
			case <-timer:
				fmt.Println("Остановка по таймеру")
				return
			default:
				fmt.Println("Горутина работает")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	time.Sleep(2 * time.Second)
}
