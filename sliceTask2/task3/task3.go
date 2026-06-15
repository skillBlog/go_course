package main

import "fmt"

func test(testSlice []string) []string {
	return append(testSlice, "Пока")
}

func main() {
	testSlice := make([]string, 0, 3)
	testSlice = append(testSlice, "Привет")
	testSlice = append(testSlice, "Привет")
	testSlice = test(testSlice)
	fmt.Println(testSlice)
}
