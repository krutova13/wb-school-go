package main

import (
	"flag"
	"math/rand"
	"time"
	"wbschoolgo/lvl1/l3"
)

func main() {
	n := flag.Int("workers", 3, "количество воркеров")
	flag.Parse()

	ch := make(chan int)

	for i := 0; i <= *n; i++ {
		go l3.Worker(i, ch)
	}

	for {
		data := rand.Intn(100)
		ch <- data
		time.Sleep(500 * time.Millisecond)
	}
}
