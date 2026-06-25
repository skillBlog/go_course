package main

import (
	"fmt"
	"sync"
	"time"
)

// Connection: имитация подключения к БД
type Connection struct {
	ID int
}

// ConnectionPool: пул с ограничением на число одновременных подключений
type ConnectionPool struct {
	mu        sync.Mutex
	cond      *sync.Cond
	available []*Connection // свободные подключения
	maxSize   int
}

func NewConnectionPool(size int) *ConnectionPool {
	p := &ConnectionPool{
		available: make([]*Connection, 0, size),
		maxSize:   size,
	}
	p.cond = sync.NewCond(&p.mu)

	for i := 1; i <= size; i++ {
		p.available = append(p.available, &Connection{ID: i})
	}
	return p
}

// Get: возвращает свободное подключение. Если все заняты, ждёт на cond
func (p *ConnectionPool) Get() *Connection {
	p.mu.Lock()
	defer p.mu.Unlock()

	for len(p.available) == 0 {
		p.cond.Wait()
	}

	conn := p.available[len(p.available)-1]
	p.available = p.available[:len(p.available)-1]
	return conn
}

// Release: возвращает подключение в пул и разбудит одного ожидающего
func (p *ConnectionPool) Release(conn *Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.available = append(p.available, conn)
	p.cond.Signal()
}

func main() {
	pool := NewConnectionPool(3) // Пул на 3 подключения

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()
			conn := pool.Get()
			defer pool.Release(conn)

			fmt.Printf("Горутина %d: подключение %d получено\n", id, conn.ID)
			time.Sleep(2 * time.Second) // Имитация работы
		}(i)
	}

	wg.Wait()

	// Пример вывода (порядок строк может отличаться, но логика одна):
	//   Горутина 0: подключение 3 получено
	//   Горутина 1: подключение 2 получено
	//   Горутина 2: подключение 1 получено
	//   ... через примерно 2 сек, когда кто-то сделает Release:
	//   Горутина 7: подключение 1 получено
	//   Горутина 3: подключение 2 получено
	//   ...
	//
	// Почему так:
	//   1. Сразу стартуют 10 горутин, но пул даёт только 3 подключения (ID 1, 2, 3).
	//      Первые три, кто успели вызвать Get(), печатают и уходят в Sleep(2 сек).
	//      Остальные 7 блокируются на cond.Wait(): все подключения заняты.
	//   2. Get() забирает подключение с конца слайса available, поэтому у первых трёх
	//      часто будут ID 3, 2, 1 (зависит от порядка планировщика).
	//   3. Через 2 сек первая горутина просыпается, defer Release() возвращает
	//      подключение в пул и Signal() будит одну из ждущих: она печатает и снова Sleep.
	//   4. wg.Wait() дожидается всех 10 горутин: main не завершится, пока каждая
	//      не отработает и не вернёт подключение через Release.
}
