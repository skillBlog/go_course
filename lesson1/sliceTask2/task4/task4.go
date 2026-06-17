package main

import "fmt"

func main() {
	first := []int{10, 20, 30, 40}
	// second слайс указателей на элементы first
	second := make([]*int, len(first))
	for i := range first {
		// &first[i] — адрес i-го элемента в массиве first (у каждого i свой адрес)
		second[i] = &first[i]
	}
	// разыменовываем указатели, получим значения из first
	fmt.Println(*second[0], *second[1]) // 10 20
}
