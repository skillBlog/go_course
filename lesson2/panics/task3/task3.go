package main

import "fmt"

func Level1() {
	// паника из Level3 поднимется через Level2, но поймается на этом уровне
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Паника обработана на уровне 1: %v\n", r)
		}
	}()
	Level2()
}

func Level2() {
	// defer выполнится даже при panic в Level3, до того, как паника дойдёт до Level1
	defer fmt.Println("Завершаем Level2")
	Level3()
}

func Level3() {
	panic("ошибка в Level3")
}

func main() {
	// выведет:
	// Завершаем Level2
	// Паника обработана на уровне 1: ошибка в Level3
	// Программа продолжает работу
	Level1()
	fmt.Println("Программа продолжает работу")
	// Программа продолжает работу, так как паника была перехвачена на уровне 1
}
