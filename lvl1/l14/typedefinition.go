package l14

import "fmt"

// DefineType определяет и выводит тип переданного значения
func DefineType(v interface{}) {
	switch v := v.(type) {
	case int:
		fmt.Printf("%v - это тип int\n", v)
	case string:
		fmt.Printf("%v - это тип string\n", v)
	case bool:
		fmt.Printf("%v - это тип bool\n", v)
	case chan int:
		fmt.Printf("%v - это канал int\n", v)
	case chan string:
		fmt.Printf("%v - это канал string\n", v)
	case chan bool:
		fmt.Printf("%v - это канал bool\n", v)
	case chan interface{}:
		fmt.Printf("%v - это канал interface{}\n", v)
	default:
		fmt.Println("Неизвестный тип данных")
	}
}
