package main

import "fmt"

func main() {
	// len=0, cap=5, после append четыре элемента, len=4, cap=5
	slice := make([]string, 0, 5)
	slice = append(slice, "0", "1", "2", "3")
	fmt.Println(slice, len(slice), cap(slice)) // [0 1 2 3] 4 5

	// slice передаётся копией заголовка, но массив общий
	// addToSlice1 пишет one в индекс 3, затирает 3, len в main не меняется
	addToSlice1(slice)
	fmt.Println(slice, len(slice), cap(slice)) // [0 1 2 one] 4 5

	// addToSlice2 дописывает two в индекс 4, но len в main всё ещё 4, two не виден
	addToSlice2(slice)
	fmt.Println(slice, len(slice), cap(slice)) // [0 1 2 one] 4 5
}

func addToSlice1(slice []string) {
	slice = append(slice[1:3], "one")
}

func addToSlice2(slice []string) {
	slice = append(slice, "two")
}
