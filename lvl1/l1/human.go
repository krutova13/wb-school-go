package l1

import "fmt"

// Human представляет человека с именем, фамилией и возрастом
type Human struct {
	Name    string
	Surname string
	Age     int
}

// Greeting возвращает приветствие человека
func (h *Human) Greeting() string {
	return fmt.Sprintf("Привет, меня зовут %v %v", h.Name, h.Surname)
}

// Action представляет действие, выполняемое человеком
type Action struct {
	Human
	ActionName string
}
