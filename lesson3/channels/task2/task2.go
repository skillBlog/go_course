package main

import (
	"time"
)

func worker() chan int {
	ch := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- 42
	}()
	return ch
}

// func main() {
// 	timeStart := time.Now()
// 	// Второй worker() не запустится, пока не завершится первый worker
// 	// Поэтому ожидание идёт последовательно, 3 сек + 3 сек = 6 сек
// 	_, _ = <-worker(), <-worker()
// 	println(int(time.Since(timeStart).Seconds())) // 6
// }

func main() {
	timeStart := time.Now()

	// сначала запускаем горутины
	ch1 := worker()
	ch2 := worker()

	// теперь обе горутины уже работают одновременно
	// и обе спят свои 3 секунды параллельно
	_, _ = <-ch1, <-ch2
	println(int(time.Since(timeStart).Seconds())) // 3
}
