package main

import (
	"fmt"
)

func main() {
	slice := make([]string, 3, 4)
	fmt.Println(slice)

	slice = appendSlice(slice)
	fmt.Println(slice)

	mutareSlice(slice)
	fmt.Println(slice)
}

func appendSlice(slice []string) []string {
	return append(slice, "privet")
}

func mutareSlice(slice []string) {
	slice[0] = "vasya"
}
