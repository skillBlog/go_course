package main

import "fmt"

func main() {
	first := []int{10, 20, 30, 40}
	second := make([]*int, len(first))
	for i := range first {
		// &first[i] - адрес i-го элемента первого слайса
		second[i] = &first[i]
	}
	fmt.Println(*second[0], *second[1])
}
