package main

import (
	"fmt"
)

// выведет: start, end, 3, 2, 1, потому что defer работает как стек и выполняется в обратном порядке
func main() {
	fmt.Println("start")
	for i := 1; i < 4; i++ {
		defer fmt.Println(i)
	}
	fmt.Println("end")
}
