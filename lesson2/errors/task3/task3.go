package main

import (
	"errors"
	"fmt"
	"math/rand"
)

var (
	ErrNotFound  = errors.New("ресурс не найден")
	TimeoutError = errors.New("таймаут операции")
)

func SimulateRequest() error {
	// случайная ошибка: 0-4 таймаут, 5-7 not found, 8-9 без обёрток
	switch rand.Intn(10) {
	case 0, 1, 2, 3, 4:
		return fmt.Errorf("запрос не выполнен: %w", TimeoutError)
	case 5, 6, 7:
		return fmt.Errorf("ошибка: %w", ErrNotFound)
	default:
		return errors.New("неизвестная ошибка")
	}
}

func ProcessError(err error) {
	// errors.Is ищет ошибку в цепочке и находит даже обёрнутые через %w
	if errors.Is(err, TimeoutError) {
		fmt.Println("Требуется повторная попытка")
	} else if errors.Is(err, ErrNotFound) {
		fmt.Println("Ресурс не найден")
	} else {
		fmt.Println("Неизвестная ошибка")
	}
}

func main() {
	// выполняем 10 раз и выводим текст ошибки и результат ProcessError
	for i := 0; i < 10; i++ {
		err := SimulateRequest()
		fmt.Println(err)  // текст ошибки, например: запрос не выполнен: таймаут операции
		ProcessError(err) // одно из трёх сообщений в зависимости от типа ошибки
		fmt.Println()
	}
}
