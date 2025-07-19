package main

import (
	"fmt"
	"sync"
	"wbschoolgo/lvl1/l18"
)

func main() {
	var wg sync.WaitGroup
	counter := &l18.Counter{}

	goroutines := 10
	increments := 1000

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				counter.Inc()
			}
		}()
	}
	wg.Wait()
	fmt.Println("Итоговое значение счётчика:", counter.Value())
}
