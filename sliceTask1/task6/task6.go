package main

import "fmt"

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	bar := arr[1:3]
	bar = append(bar, 10, 11, 12, 13)
	fmt.Println(arr, bar) // [1 2 3 4 5] [2 3 10 11 12 13]
}
