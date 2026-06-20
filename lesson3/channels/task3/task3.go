package main

import "fmt"

func main() {
	// исходные числа
	naturals := make(chan int)
	// квадраты чисел
	squares := make(chan int)

	// отправляем числа в канал naturals
	// после отправки всех чисел закрываем канал, чтобы сигнализировать, что данных больше не будет
	go func() {
		for x := 0; x < 100; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	// читаем из naturals по мере поступления,
	// возводим каждое число в квадрат и отправляем результат в squares
	// когда naturals закрыт и опустошён, закрываем squares
	go func() {
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	// читаем готовые квадраты и выводим в консоль
	for x := range squares {
		fmt.Println(x)
	}
}
