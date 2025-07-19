package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	in := make(chan int, len(arr))
	out := make(chan int, len(arr))

	go func() {
		for _, v := range arr {
			in <- v
		}
		close(in)
	}()

	go func() {
		for v := range in {
			out <- v * 2
		}
		close(out)
	}()

	for v := range out {
		fmt.Println(v)
	}
}
