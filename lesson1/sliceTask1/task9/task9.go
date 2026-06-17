package main

import (
	"fmt"
)

func main() {
	// len=3, cap=4, три нуля: [0 0 0]
	slice := make([]int, 3, 4)
	// slice[:1] len=1, но cap=4 и общий массив с slice
	appendingSlice(slice[:1])
	// append записал 1 в индекс 1 общего массива
	fmt.Println(slice) // [0 1 0]
}

func appendingSlice(slice []int) {
	slice = append(slice, 1)
}
