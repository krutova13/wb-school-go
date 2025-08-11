package l24

import "math"

// Point представляет точку в двумерном пространстве
type Point struct {
	x float64
	y float64
}

// NewPoint создает новую точку с заданными координатами
func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

// Distance вычисляет расстояние между двумя точками
func (p Point) Distance(other Point) float64 {
	dx := p.x - other.x
	dy := p.y - other.y
	return math.Sqrt(dx*dx + dy*dy)
}
