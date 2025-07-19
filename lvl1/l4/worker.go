package l4

import (
	"context"
	"fmt"
	"sync"
)

func Worker(ctx context.Context, id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case data, ok := <-ch:
			if !ok {
				fmt.Printf("Воркер %d: канал закрыт, выхожу\n", id)
				return
			}
			fmt.Printf("Воркер %d получил значение: %d\n", id, data)
		case <-ctx.Done():
			fmt.Printf("Воркер %d останавливается\n", id)
			return
		}
	}
}
