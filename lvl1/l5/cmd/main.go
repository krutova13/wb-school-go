package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	timeout := flag.Int("timeout", 5, "время работы программы в секундах")
	flag.Parse()
	timer := time.After(time.Duration(*timeout) * time.Second)
	ch := make(chan int)

	go func() {
		for {
			data := rand.Intn(100)
			ch <- data
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for {
		select {
		case data := <-ch:
			fmt.Printf("Получено значение: %d\n", data)
		case <-timer:
			fmt.Println("\nВремя вышло, остановка программы")
			return
		}
	}
}
