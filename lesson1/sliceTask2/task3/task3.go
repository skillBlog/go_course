package main

import "fmt"

func test(testSlice []string) []string {
	// append может увеличить len, результат нужно вернуть и присвоить в main
	return append(testSlice, "Пока")
}

func main() {
	testSlice := make([]string, 0, 3)       // len=0, cap=3
	testSlice = append(testSlice, "Привет") // [Привет], len=1, cap=3
	testSlice = append(testSlice, "Привет") // [Привет Привет], len=2, cap=3
	testSlice = test(testSlice)             // добавили Пока, len=3, cap=3, влезает в cap
	fmt.Println(testSlice)                  // [Привет Привет Пока], выводим результат
}
