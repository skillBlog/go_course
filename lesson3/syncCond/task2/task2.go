package main

import (
	"fmt"
	"sync"
	"time"
)

// Restaurant: ресторан с фиксированным числом столиков
type Restaurant struct {
	mu       sync.Mutex
	cond     *sync.Cond
	tables   int // всего столиков
	occupied int // сколько занято сейчас
}

func NewRestaurant(tables int) *Restaurant {
	r := &Restaurant{tables: tables}
	// создаем новый sync.Cond с использованием Mutex
	// это позволяет горутинам блокироваться и разблокироваться
	// при помощи одного и того же Mutex
	r.cond = sync.NewCond(&r.mu)
	// возвращаем новый Restaurant
	return r
}

// OccupyTable: посетитель занимает столик. Если все заняты, ждёт в очереди.
func (r *Restaurant) OccupyTable() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for r.occupied >= r.tables {
		r.cond.Wait() // все столики заняты, спим, пока кто-то не вызовет ReleaseTable
	}
	r.occupied++
}

// ReleaseTable: посетитель уходит, столик свободен. Разбудим одного из очереди
func (r *Restaurant) ReleaseTable() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.occupied > 0 {
		r.occupied--
	}
	r.cond.Signal()
}

func main() {
	const tables = 5
	const visitors = 8 // больше, чем столиков, часть будет ждать

	restaurant := NewRestaurant(tables)
	var wg sync.WaitGroup

	for i := 1; i <= visitors; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			fmt.Printf("посетитель %d: жду столик...\n", id)
			restaurant.OccupyTable()
			fmt.Printf("посетитель %d: сел за столик\n", id)

			time.Sleep(100 * time.Millisecond) // имитация еды
			fmt.Printf("посетитель %d: ушёл\n", id)

			restaurant.ReleaseTable()
		}(i)
	}

	wg.Wait()

	// вывод: первые 5 посетителей садятся сразу, остальные 3 ждут.
	// как только кто-то уходит (ReleaseTable + Signal), следующий из очереди садится.
	// порядок "сел" у 6-8 зависит от того, кто быстрее освободил столик.
}
