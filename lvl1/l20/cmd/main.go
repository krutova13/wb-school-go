package main

import (
	"fmt"
	"wbschoolgo/lvl1/l20"
)

func main() {
	sentence := "snow dog sun"
	reversedSentence := l20.ReverseSentence(sentence)
	fmt.Printf("Перевернутое предложение: %v", reversedSentence)
}
