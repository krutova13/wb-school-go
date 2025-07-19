package main

import (
	"fmt"
	"wbschoolgo/lvl1/l24"
)

func main() {
	var x1, y1, x2, y2 float64
	fmt.Print("Введите координаты первой точки (x1 y1): ")
	_, err := fmt.Scan(&x1, &y1)
	if err != nil {
		return
	}
	fmt.Print("Введите координаты второй точки (x2 y2): ")
	_, err = fmt.Scan(&x2, &y2)
	if err != nil {
		return
	}

	p1 := l24.NewPoint(x1, y1)
	p2 := l24.NewPoint(x2, y2)

	distance := p1.Distance(p2)
	fmt.Printf("Расстояние между точками = %v", distance)
}
