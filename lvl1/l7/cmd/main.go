package main

import (
	"fmt"
	"sync"
	"wbschoolgo/lvl1/l7"
)

func main() {
	var (
		mu   sync.Mutex
		data = make(map[int]int)
		wg   sync.WaitGroup
	)

	numWorkers := 3
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go l7.Worker(i, &mu, data, &wg)
	}

	wg.Wait()
	fmt.Println("Результат:", data)
}
