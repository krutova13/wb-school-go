package l13

// Exchange обменивает значения двух переменных без использования временной переменной
func Exchange(a, b int) (int, int) {
	a = a + b
	b = a - b
	a = a - b
	return a, b
}
