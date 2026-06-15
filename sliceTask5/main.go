package main

// Stacker — интерфейс стека с операциями LIFO.
type Stacker interface {
	Push(v int)
	Pop() int
}

// stack — реализация стека на слайсе.
type stack struct {
	items []int
}

// Push добавляет значение на вершину стека.
func (s *stack) Push(v int) {
	s.items = append(s.items, v)
}

// Pop извлекает и возвращает верхний элемент (последний добавленный).
// LIFO: последним пришёл — первым ушёл.
// Если стек пуст, вызывает panic.
func (s *stack) Pop() int {
	if len(s.items) == 0 {
		panic("stack is empty")
	}

	// Верх стека — последний элемент слайса.
	item := s.items[len(s.items)-1]

	// Укорачиваем слайс на 1: элемент «снимается» с вершины.
	// Ёмкость (cap) не меняется — память под массив остаётся для следующих Push.
	s.items = s.items[:len(s.items)-1]

	return item
}

// New создаёт новый пустой стек с начальной ёмкостью 10.
func New() *stack {
	return &stack{items: make([]int, 0, 10)}
}
