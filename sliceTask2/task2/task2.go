package main

import (
	"fmt"
	"slices"
	"strings"
)

func chengeSlice(arr []string) {
	arr[0] = "Goodbye"
}

func appendSomeData(arr []string) []string {
	return append(arr, "!")
}

func main() {
	someSlice := []string{"Hello", "World"}
	// создаем копию слайса, чтобы не изменить исходный
	chengeSlice(slices.Clone(someSlice))
	someSlice = appendSomeData(someSlice)
	fmt.Println(strings.Join(someSlice, ""))
}
