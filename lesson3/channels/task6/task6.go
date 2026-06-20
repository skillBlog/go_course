package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// dbReplica имитирует запись в одну реплику БД.
func dbReplica(name string, in <-chan int) {
	for data := range in {
		fmt.Printf("Запись в %s: %d\n", name, data)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("Реплика %s закрыта\n", name)
}

// разветвитель: читает из input и отправляет каждое значение во все каналы реплик
// когда input закрыт и опустошен, закрывает каналы реплик
func tee(input <-chan int, replicas []chan int) {
	// читаем из input, пока канал не закроют (close в main)
	for data := range input {
		// каждое значение отправляем во все реплики, суть паттерна tee
		for _, r := range replicas {
			r <- data // блокируемся, пока dbReplica не примет значение
		}
	}
	// input исчерпан, закрываем каналы реплик, чтобы dbReplica вышла из for range
	for _, r := range replicas {
		close(r)
	}
}

func main() {
	input := make(chan int) // Канал для входящих данных
	replicas := []chan int{ // Реплики БД (каналы)
		make(chan int),
		make(chan int),
		make(chan int),
	}

	var wg sync.WaitGroup

	// каждая реплика читает из своего канала (не из общего input)
	for i, ch := range replicas {
		wg.Add(1)
		go func(name string, in <-chan int) {
			defer wg.Done()
			dbReplica(name, in)
		}(strconv.Itoa(i), ch)
	}

	// tee-горутина дублирует данные из input во все реплики
	go tee(input, replicas)

	// отправляем 10 значений и закрываем входной канал
	for i := 0; i < 10; i++ {
		input <- i
	}
	close(input)

	// ждем, пока все реплики обработают данные и завершатся
	wg.Wait()

	// вывод: 30 строк «Запись в X: N» (3 реплики × 10 значений) + 3 строки «Реплика X закрыта»
	// порядок строк «Запись» не фиксирован, реплики работают параллельно и sleep(100ms)
	// гарантировано: каждая реплика (0, 1, 2) получит все числа от 0 до 9
}
