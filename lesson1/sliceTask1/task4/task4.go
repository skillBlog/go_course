package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	// b подслайс a[:2], но cap=3, общий массив, место под третий элемент есть
	b := a[:2]
	// append пишет 4 в индекс 2, затирает 3 в массиве a, realloc не нужен
	b = append(b, 4)
	fmt.Println(b) // [1 2 4]
	fmt.Println(a) // [1 2 4], a видит то же изменение в общем массиве
}
