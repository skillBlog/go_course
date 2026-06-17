package main

import "fmt"

func main() {
	arr := []int{1, 2, 3}
	// src подслайс [1], но cap=3, общий массив с arr
	src := arr[:1]
	// foo пишет 5 в индекс 1, затирая 2, len src в main не меняется
	foo(src)
	fmt.Println(src) // [1], len=1, второй элемент не виден
	fmt.Println(arr) // [1 5 3], arr видит 5 на месте бывшей 2
}

func foo(src []int) {
	src = append(src, 5)
}
