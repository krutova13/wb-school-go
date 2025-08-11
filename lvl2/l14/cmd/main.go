package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

// or объединяет несколько done-каналов в один
// Возвращаемый канал закрывается, как только закроется любой из исходных каналов
var or func(channels ...<-chan interface{}) <-chan interface{}

func init() {
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		if len(channels) == 0 {
			c := make(chan interface{})
			close(c)
			return c
		}

		if len(channels) == 1 {
			return channels[0]
		}

		result := make(chan interface{})

		go func() {
			defer close(result)

			if len(channels) == 2 {
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
				return
			}

			select {
			case <-channels[0]:
			case <-or(channels[1:]...):
			}
		}()

		return result
	}
}
