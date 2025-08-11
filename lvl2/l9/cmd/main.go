package main

import (
	"fmt"
)

func main() {
	str := "a4bc2d5e\\45"
	result, err := UnpackString(str)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
	fmt.Println(result)
}
