package main

import (
	"fmt"
	"sync"
	"time"
)

// BoundedQueue: очередь с лимитом мест (capacity)
// sync.Cond: нужен, чтобы горутины ждали "не пусто" / "не полно", а не крутились в цикле
type BoundedQueue struct {
	mu       sync.Mutex
	notEmpty *sync.Cond    // будит Get, когда появилась задача (не пусто)
	notFull  *sync.Cond    // будит Put, когда освободилось место (не полно)
	queue    []interface{} // очередь задач
	capacity int           // лимит мест в очереди
	closed   bool          // true после Shutdown: очередь закрыта
}

// NewBoundedQueue: создаёт очередь с лимитом мест (capacity)
func NewBoundedQueue(capacity int) *BoundedQueue {
	q := &BoundedQueue{
		queue:    make([]interface{}, 0, capacity),
		capacity: capacity,
	}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

// Put: кладёт задачу. Если мест нет, спит на notFull, пока кто-то не сделает Get
func (q *BoundedQueue) Put(task interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// for, а не if: после пробуждения нужно перепроверить условие
	for len(q.queue) >= q.capacity && !q.closed {
		q.notFull.Wait()
	}
	if q.closed {
		return // после Shutdown: новые задачи не принимаем
	}

	q.queue = append(q.queue, task)
	q.notEmpty.Signal() // разбудить одного ожидающего Get
}

// Get: забирает задачу. Если очередь пуста, спит на notEmpty.
// после Shutdown: сначала отдаёт остаток, потом nil
func (q *BoundedQueue) Get() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.queue) == 0 && !q.closed {
		q.notEmpty.Wait()
	}
	if len(q.queue) == 0 {
		return nil // очередь закрыта и пуста, консьюмер может выйти
	}

	task := q.queue[0]
	q.queue = q.queue[1:]
	q.notFull.Signal() // освободилось место, разбудить Put
	return task
}

// Shutdown: закрывает очередь. Broadcast будит ВСЕХ ждущих (и Put, и Get)
func (q *BoundedQueue) Shutdown() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.closed = true
	q.notEmpty.Broadcast()
	q.notFull.Broadcast()
}

func main() {
	// capacity=2: в очереди максимум 2 задачи, 3-й Put будет ждать пробуждения
	queue := NewBoundedQueue(2)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			task := fmt.Sprintf("task-%d", i)
			fmt.Printf("Put: %s\n", task)
			queue.Put(task)
		}
		fmt.Println("продюсер завершил отправку")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			task := queue.Get()
			if task == nil {
				fmt.Println("консьюмер: очередь закрыта, выход")
				return
			}
			fmt.Printf("Get: %s\n", task)
			time.Sleep(50 * time.Millisecond) // консьюмер медленнее, очередь заполняется
		}
	}()

	time.Sleep(300 * time.Millisecond)
	queue.Shutdown()
	wg.Wait()

	// Пример вывода:
	//   Put: task-1 -> Get: task-1 -> Put: task-2 -> Put: task-3 (ждал места) -> ...
	//   производитель завершил отправку
	//   потребитель: очередь закрыта, выход
	//
	// Соответствие требованиям:
	//   Put блокируется при полной очереди, Get при пустой
	//   sync.Cond + Mutex: Shutdown завершает работу без утечек горутин
}
