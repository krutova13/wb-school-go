package l13

func Exchange(a, b int) (int, int) {
	a = a + b
	b = a - b
	a = a - b
	return a, b
}
