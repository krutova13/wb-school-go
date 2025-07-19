package main

import "wbschoolgo/lvl1/l14"

func main() {
	l14.DefineType(true)
	l14.DefineType("строка")
	l14.DefineType(5)
	l14.DefineType(make(chan int))
}
