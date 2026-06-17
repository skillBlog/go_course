package main

import "fmt"

func main() {
	a := []string{"a", "b", "c"}
	// b один элемент a[1], len=1, cap=2 (от индекса 1 до конца массива a)
	b := a[1:2]
	fmt.Println(b, cap(b), len(b)) // [b] 2 1

	// b[0] это тот же a[1] в общем массиве
	b[0] = "q"
	fmt.Println(a) // [a q c]
}
