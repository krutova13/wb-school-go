package l7

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Worker(id int, mu *sync.Mutex, data map[int]int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		key := rand.Intn(10)
		value := rand.Intn(100)
		mu.Lock()
		data[key] = value
		fmt.Printf("Воркер %d записал: [%d] = %d\n", id, key, value)
		mu.Unlock()
		time.Sleep(time.Millisecond * 100)
	}
}
