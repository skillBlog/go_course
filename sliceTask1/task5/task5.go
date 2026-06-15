package main

import "fmt"

func main() {
	arr := []int{1, 2, 3}
	src := arr[:1]
	foo(src)
	fmt.Println(src) // [1]
	fmt.Println(arr) // [1 5 3]
}

func foo(src []int) {
	src = append(src, 5)
}
