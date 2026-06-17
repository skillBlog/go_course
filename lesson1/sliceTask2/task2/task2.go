package main

import (
	"fmt"
	"strings"
)

func chengeSlice(arr []string) {
	// arr[0] = Goodbye затронул бы общий массив с someSlice
	arr = append([]string{"Goodbye"}, arr[1:]...) // меняется только локальный заголовок arr
}

func appendSomeData(arr []string) []string {
	return append(arr, "!")
}

func main() {
	someSlice := []string{"Hello", "World"}
	chengeSlice(someSlice) // Hello не меняется, в функции другой заголовок слайса
	// добавляем ! в конец слайса
	someSlice = appendSomeData(someSlice)
	fmt.Println(strings.Join(someSlice, "")) // HelloWorld!
}
