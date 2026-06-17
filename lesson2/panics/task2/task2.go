package main

import "fmt"

func SafeDivide(a, b int) (result int) {
	// Если внутри функции случится panic, то эта отложенная функция всё равно выполнится
	defer func() {
		// recover ловит panic; result = 0 возвращаем 0 вместо падения программы
		if r := recover(); r != nil {
			result = 0
		}
	}()

	if b == 0 {
		// Намеренно вызываем panic, она поймается defer выше
		panic("деление на ноль")
	}

	return a / b
}

func main() {
	// Не будет падения, вернёт 5
	fmt.Println(SafeDivide(10, 2)) // 5
	// Будет падение, вернёт 0
	fmt.Println(SafeDivide(10, 0)) // 0
}
