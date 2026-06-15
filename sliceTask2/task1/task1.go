package main

import (
	"fmt"
)

// тут и так правильно выводится, т.к версия го новая
// С Go 1.22 (февраль 2024) изменили семантику range:
// на каждой итерации создаётся новая переменная value. Поэтому:
// numbers = append(numbers, &value)
// даёт 4 разных адреса - вывод 10, 20, 30, 40.

func main() {
	var numbers []*int

	for _, value := range []int{10, 20, 30, 40} {
		numbers = append(numbers, &value)
	}

	for _, number := range numbers {
		fmt.Println("d", *number)
	}
}
