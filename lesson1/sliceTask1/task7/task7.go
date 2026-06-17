package main

import "fmt"

func main() {
	a := []string{"a", "b", "c"}
	b := a[1:2]
	fmt.Println(b, cap(b), len(b)) // [b] 2 1
	b[0] = "q"
	fmt.Println(a) // [a q c]
}
