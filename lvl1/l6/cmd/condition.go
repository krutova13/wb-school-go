package main

import (
	"fmt"
	"time"
)

func main() {
	stop := false

	go func() {
		for {
			if stop {
				fmt.Println("Остановлено по флагу")
				return
			}
			fmt.Println("Горутина работает")
			time.Sleep(500 * time.Millisecond)
		}
	}()

	time.Sleep(2 * time.Second)
	stop = true
	time.Sleep(1 * time.Second)
}
