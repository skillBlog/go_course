package main

import (
	"fmt"
)

const totalCount = 100

type squareResult struct {
	value int
	err   error
}

func square(x int) (int, error) {
	if x < 0 {
		return 0, fmt.Errorf("отрицательное число: %d", x)
	}
	return x * x, nil
}

func main() {
	// буфер = totalCount: продюсер и обработчик не блокируют друг друга на каждой отправке
	naturals := make(chan int, totalCount)
	squares := make(chan squareResult, totalCount)

	go func() {
		defer close(naturals)
		for x := 0; x < totalCount; x++ {
			naturals <- x
		}
	}()

	go func() {
		defer close(squares)
		for x := range naturals {
			v, err := square(x)
			squares <- squareResult{value: v, err: err}
		}
	}()

	for r := range squares {
		if r.err != nil {
			fmt.Println("ошибка:", r.err)
			continue
		}
		fmt.Println(r.value)
	}
}
