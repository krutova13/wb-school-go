package main

import (
	"fmt"
	"wbschoolgo/lvl1/l1"
)

func main() {
	a := l1.Action{
		Human: l1.Human{
			Name:    "Анастасия",
			Surname: "Артемова",
			Age:     29,
		},
		ActionName: "Представление: ",
	}
	fmt.Println(a.ActionName)
	fmt.Println(a.Greeting())
}
