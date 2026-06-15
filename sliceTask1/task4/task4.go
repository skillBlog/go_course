package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	b := a[:2]
	b = append(b, 4)
	fmt.Println(b) // [1 2 4]
	fmt.Println(a) // [1 2 4]
}
