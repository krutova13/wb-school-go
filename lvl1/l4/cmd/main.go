package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"wbschoolgo/lvl1/l4"
)

func main() {

	n := flag.Int("workers", 3, "количество воркеров")
	flag.Parse()

	ch := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Получено прерывание, остановка воркеров...")
		cancel()
	}()

	var wg sync.WaitGroup

	for i := 0; i < *n; i++ {
		wg.Add(1)
		go l4.Worker(ctx, i, ch, &wg)
	}

	for {
		select {
		case <-ctx.Done():
			close(ch)
			wg.Wait()
			fmt.Println("Остановка главной горутины")
			return
		default:
			data := rand.Intn(100)
			ch <- data
			time.Sleep(500 * time.Millisecond)
		}
	}
}
