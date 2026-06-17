package main

import "fmt"

func main() {
	value := 123

	// Замыкание читает value в момент выполнения defer,
	// поэтому будет выведено уже изменённое значение 456.
	defer func() {
		fmt.Println(value)
	}()

	changeValue(&value)
}

func changeValue(value *int) {
	*value = 456
}
