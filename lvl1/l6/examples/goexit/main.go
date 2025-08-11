package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	go func() {
		for {
			fmt.Println("Горутина работает")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Остановлено через Goexit")
			runtime.Goexit()
		}
	}()
	time.Sleep(2 * time.Second)
}
