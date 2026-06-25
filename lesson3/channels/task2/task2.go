package main

import (
	"sync"
	"time"
)

func worker() chan int {
	ch := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- 42
		close(ch)
	}()
	return ch
}

// fanIn объединяет несколько каналов в один (паттерн fan-in)
func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, ch := range channels {
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
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

	// оба воркера стартуют сразу, fan-in сливает их каналы в один
	merged := fanIn(worker(), worker())
	<-merged
	<-merged

	println(int(time.Since(timeStart).Seconds())) // 3
}
