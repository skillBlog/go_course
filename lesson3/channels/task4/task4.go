package main

import (
	"fmt"
	"sync"
)

func mergeChannels(channels ...<-chan int) <-chan int {
	// создаем канал для объединения входных каналов
	out := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, channel := range channels {
		// channel передается аргументом, чтобы каждая горутина читала свой канал
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(channel)
	}

	// отдельная горутина закрывает out, когда все пересыльщики завершились
	// т.е. все входные каналы прочитаны до конца
	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	// заполняем входные каналы в отдельных горутинах и закрываем их
	// иначе mergeChannels будет ждать данных бесконечно
	go func() {
		a <- 1
		a <- 2
		close(a) // в a: 1, 2
	}()
	go func() {
		b <- 3
		close(b) // в b: 3
	}()
	go func() {
		c <- 4
		c <- 5
		close(c) // в c: 4, 5
	}()

	ch := mergeChannels(a, b, c)

	// читаем объединенный канал до закрытия
	// вывод: 5 строк с числами 1, 2, 3, 4, 5 - но порядок  фиксирован
	// три горутины пишут в out параллельно, планировщик решает, кто успеет первым
	// например: 1, 3, 4, 2, 5  или  4, 1, 2, 3, 5 - оба варианта возможны
	// гарантировано только то, что каждое число появится ровно один раз
	for v := range ch {
		fmt.Println(v)
	}
}
