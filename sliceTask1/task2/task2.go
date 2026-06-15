package main

import "fmt"

func main() {
	slice := make([]string, 0, 5)
	slice = append(slice, "0", "1", "2", "3")
	fmt.Println(slice, len(slice), cap(slice)) // [0 1 2 3] 4 5
	addToSlice1(slice)
	fmt.Println(slice, len(slice), cap(slice)) // [0 1 2 one] 4 5
	addToSlice2(slice)
	fmt.Println(slice, len(slice), cap(slice)) // [0 1 2 one] 4 5
}

func addToSlice1(slice []string) {
	slice = append(slice[1:3], "one")
}

func addToSlice2(slice []string) {
	slice = append(slice, "two")
}
