package main

import "fmt"

type account struct {
	value int
}

func main() {
	// len=0, cap=2, место под 2 элемента, но пока пусто
	s1 := make([]account, 0, 2)
	s1 = append(s1, account{})
	// cap хватает, realloc не нужен, s2 делит с s1 один и тот же массив в памяти
	s2 := append(s1, account{})
	// указатель на s2[0], это тот же элемент, что и s1[0] в общем массиве
	acc := &s2[0]
	acc.value = 100
	// s1 видит только 1 элемент [{100}], s2 видит оба [{100} {0}]
	fmt.Println(s1, s2)

	// снова append в s1: cap=2, len был 1, вторая ячейка уже занята, realloc не происходит
	s1 = append(s1, account{})
	// acc всё ещё указывает на тот же s1[0]/s2[0] в общем массиве
	acc.value += 100
	// меняем общий элемент → оба слайса показывают 200 в первой позиции
	fmt.Println(s1, s2) // [{200} {0}] [{200} {0}]
}
