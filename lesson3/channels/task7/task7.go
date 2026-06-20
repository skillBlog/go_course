package main

import (
	"fmt"
	"sync"
)

type ServerMetric struct {
	Name  string  // название метрики (например, "memory_usage")
	Value float64 // значение в байтах
}

// константа для перевода байтов в мегабайты
const bytesInMB = 1024 * 1024

// декоратор (паттерн transformer): читает метрики в байтах
// возвращает новый канал с метриками в мегабайтах
func transformMetrics(in <-chan ServerMetric) <-chan ServerMetric {
	out := make(chan ServerMetric) // новый канал с метриками в мегабайтах

	// читаем метрики в байтах и отправляем в новый канал с метриками в мегабайтах
	go func() {
		defer close(out) // закрываем канал после чтения всех метрик
		for m := range in {
			out <- ServerMetric{ // отправляем метрику в новый канал с метриками в мегабайтах
				Name: m.Name,
				// делим значение на 1024 * 1024, чтобы получить мегабайты
				Value: m.Value / bytesInMB,
			}
		}
	}()

	return out
}

// produceMetrics: отправляет метрики в байтах и закрывает канал
func produceMetrics(out chan<- ServerMetric, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		out <- ServerMetric{Name: "memory_usage", Value: 1_000_000}
	}
	close(out)
}

// sendToAPI: имитирует отправку преобразованных метрик в API
func sendToAPI(in <-chan ServerMetric, wg *sync.WaitGroup) {
	defer wg.Done()
	for m := range in {
		fmt.Printf("API: %s = %.4f MB\n", m.Name, m.Value)
	}
}

func main() {
	metrics := make(chan ServerMetric)

	var wg sync.WaitGroup
	wg.Add(2)

	// источник: метрики в байтах
	go produceMetrics(metrics, &wg)

	// декоратор: переводит байты в мегабайты
	apiMetrics := transformMetrics(metrics)

	// потребитель: читает уже преобразованные метрики
	go sendToAPI(apiMetrics, &wg)

	wg.Wait()

	// вывод: 10 строк вида
	//   API: memory_usage = 0.9537 MB
	// (1_000_000 / (1024*1024) ≈ 0.95367431640625)
}
