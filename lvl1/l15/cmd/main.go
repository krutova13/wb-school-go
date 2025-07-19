package main

import "fmt"

func main() {
	fmt.Print(someFunc())
}

func someFunc() string {
	v := createHugeString(1 << 10)
	justString := make([]byte, 100)
	copy(justString, v[:100])
	return string(justString)
}

func createHugeString(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}
