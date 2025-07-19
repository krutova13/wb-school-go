package l2

func Square(numbers []int) []int {
	ch := make(chan int, len(numbers))

	for _, n := range numbers {
		go func(num int) {
			ch <- num * num
		}(n)
	}

	results := make([]int, 0, len(numbers))
	for i := 0; i < len(numbers); i++ {
		results = append(results, <-ch)
	}
	return results
}
