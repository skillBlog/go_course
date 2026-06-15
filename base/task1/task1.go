package task1

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type CircleShape struct {
	Radius float64
}

type RectangleShape struct {
	Width  float64
	Height float64
}

func (c CircleShape) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c CircleShape) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (r RectangleShape) Area() float64 {
	return r.Width * r.Height
}

func (r RectangleShape) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func main() {
	fmt.Println("Hello, World!")
}
