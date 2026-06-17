package main

import (
	"fmt"
)

func main() {
	// len=3, cap=4, три пустые строки
	slice := make([]string, 3, 4)
	fmt.Println(slice) // [   ]

	slice = appendSlice(slice) // len=4, влезает в cap, присваиваем результат
	fmt.Println(slice)         // [   privet]

	mutateSlice(slice) // slice[0] меняет общий массив, присваивание не нужно
	fmt.Println(slice) // [vasya  privet]
}

func appendSlice(slice []string) []string {
	return append(slice, "privet")
}

func mutateSlice(slice []string) {
	slice[0] = "vasya"
}
