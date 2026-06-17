package main

import "fmt"

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	// bar слайс поверх массива: [2 3], len=2, cap=4
	bar := arr[1:3]
	// нужно 6 элементов, cap=4 не влезает, append выделяет новый массив
	bar = append(bar, 10, 11, 12, 13)
	// arr не меняется, bar уже на другом массиве
	fmt.Println(arr, bar) // [1 2 3 4 5] [2 3 10 11 12 13]
}
