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
// in <-chan ServerMetric — только чтение: функция не может писать в in и закрыть его
// возвращает <-chan ServerMetric — потребитель тоже только читает, закрывает transformMetrics
func transformMetrics(in <-chan ServerMetric) <-chan ServerMetric {
	out := make(chan ServerMetric)

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
// out chan<- ServerMetric - только запись: нельзя читать из канала, зато можно close(out)
func produceMetrics(out chan<- ServerMetric, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		out <- ServerMetric{Name: "memory_usage", Value: 1_000_000}
	}
	close(out)
}

// sendToAPI: имитирует отправку преобразованных метрик в API
// in <-chan ServerMetric - только чтение: закрывать канал может только тот, кто в него пишет
func sendToAPI(in <-chan ServerMetric, wg *sync.WaitGroup) {
	defer wg.Done()
	for m := range in {
		fmt.Printf("API: %s = %.4f MB\n", m.Name, m.Value)
	}
}

func main() {
	// chan ServerMetric - двунаправленный канал, и чтение, и запись.
	// make создаёт именно chan T, при передаче в функции Go неявно сужает тип:
	//   metrics       -> produceMetrics(out chan<- ...)  - только запись
	//   metrics       -> transformMetrics(in <-chan ...)  - только чтение
	// Сужение направления на этапе компиляции защищает от случайной записи в канал-приёмник
	// или чтения из канала-источника не той стороной конвейера.
	metrics := make(chan ServerMetric)

	var wg sync.WaitGroup
	wg.Add(2)

	// produceMetrics получает send-only view: пишет 10 метрик и закрывает канал
	go produceMetrics(metrics, &wg)

	// transformMetrics читает из metrics (<-chan) и отдаёт новый read-only канал
	apiMetrics := transformMetrics(metrics)

	// sendToAPI читает только из apiMetrics; закрытие делает transformMetrics после range in
	go sendToAPI(apiMetrics, &wg)

	wg.Wait()

	// вывод: 10 строк вида
	//   API: memory_usage = 0.9537 MB
	// (1_000_000 / (1024*1024) ≈ 0.95367431640625)
}
