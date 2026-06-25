package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	replicaCount       = 3
	totalValues        = 10
	envTeeWorkers      = "TEE_WORKERS"
	defaultTeeWorkers  = 2
	replicaChannelSize = totalValues // буфер: tee не блокируется, пока реплики догоняют
)

// dbReplica имитирует запись в одну реплику БД.
func dbReplica(name string, in <-chan int) {
	for data := range in {
		fmt.Printf("Запись в %s: %d\n", name, data)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("Реплика %s закрыта\n", name)
}

// tee читает input и дублирует каждое значение во все реплики.
// пул воркеров параллельно рассылает одно значение по репликам,
// следующее значение отправляется только после подтверждения записи во все реплики.
func tee(input <-chan int, replicas []chan int, sendWorkers int) {
	if sendWorkers > len(replicas) {
		sendWorkers = len(replicas)
	}
	if sendWorkers < 1 {
		sendWorkers = 1
	}

	for data := range input {
		jobs := make(chan int, len(replicas))
		var wg sync.WaitGroup
		wg.Add(sendWorkers)

		for range sendWorkers {
			go func() {
				defer wg.Done()
				for replicaIdx := range jobs {
					replicas[replicaIdx] <- data
				}
			}()
		}

		for i := range replicas {
			jobs <- i
		}
		close(jobs)
		wg.Wait()
	}

	for _, r := range replicas {
		close(r)
	}
}

func resolveTeeWorkers() int {
	if s := os.Getenv(envTeeWorkers); s != "" {
		n, err := strconv.Atoi(s)
		if err == nil && n > 0 {
			return n
		}
	}
	if n := runtime.NumCPU(); n > 0 && n < defaultTeeWorkers {
		return n
	}
	return defaultTeeWorkers
}

func main() {
	input := make(chan int, totalValues)
	replicas := make([]chan int, replicaCount)
	for i := range replicas {
		replicas[i] = make(chan int, replicaChannelSize)
	}

	var wg sync.WaitGroup

	for i, ch := range replicas {
		wg.Add(1)
		go func(name string, in <-chan int) {
			defer wg.Done()
			dbReplica(name, in)
		}(strconv.Itoa(i), ch)
	}

	go tee(input, replicas, resolveTeeWorkers())

	for i := 0; i < totalValues; i++ {
		input <- i
	}
	close(input)

	wg.Wait()

	// вывод: 30 строк «Запись в X: N» (3 реплики × 10 значений) + 3 строки «Реплика X закрыта»
	// порядок строк «Запись» не фиксирован, реплики работают параллельно и sleep(100ms)
	// гарантировано: каждая реплика (0, 1, 2) получит все числа от 0 до 9
}
