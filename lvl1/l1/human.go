package l1

import "fmt"

type Human struct {
	Name    string
	Surname string
	Age     int
}

func (h *Human) Greeting() string {
	return fmt.Sprintf("Привет, меня зовут %v %v", h.Name, h.Surname)
}

type Action struct {
	Human
	ActionName string
}
