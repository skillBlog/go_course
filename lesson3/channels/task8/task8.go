package main

import (
	"fmt"
	"sync"
)

// parse: добавляет префикс "parsed - " к каждой строке
func parse(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for s := range in {
			out <- "parsed - " + s
		}
	}()
	return out
}

// split: распределяет данные между N каналами по round-robin
func split(in <-chan string, n int) []<-chan string {
	channels := make([]chan string, n)
	readOnly := make([]<-chan string, n)

	// создаем каналы для каждого worker
	for i := range channels {
		channels[i] = make(chan string)
		readOnly[i] = channels[i]
	}

	go func() {
		idx := 0
		for s := range in {
			// idx % n - остаток от деления; индекс ходит по кругу: 0, 1, 2, 0, 1, 2...
			channels[idx%n] <- s
			idx++
		}
		for _, ch := range channels {
			close(ch)
		}
	}()

	return readOnly
}

// send: N горутин читают свои каналы, добавляют "sent - "
// и отправляют результат в общий выходной канал (fan-in)
func send(channels []<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for s := range c {
				out <- "sent - " + s
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// создаем входной канал для сырых данных
	raw := make(chan string)

	// источник сырых данных
	go func() {
		messages := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta"}
		for _, msg := range messages {
			raw <- msg
		}
		close(raw)
	}()

	// конвейер: parse -> split -> send
	const workers = 3
	parsed := parse(raw)
	channels := split(parsed, workers)
	results := send(channels)

	// вывод: 6 строк с полной цепочкой префиксов
	// каждая строка: "sent - parsed - <исходное слово>"
	//
	// round-robin при workers=3:
	//   alpha  -> worker 0    beta   -> worker 1    gamma  -> worker 2
	//   delta  -> worker 0    epsilon-> worker 1    zeta   -> worker 2
	//
	// порядок строк в консоли не фиксирован, worker-горутины работают параллельно
	// все 6 сообщений пройдут через все три этапа
	for result := range results {
		fmt.Println(result)
	}
}
