package main

import "fmt"

type Stack[T any] struct {
	elements []T
}

func (s *Stack[T]) Push(item T) {
	s.elements = append(s.elements, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	// пустой стек, возвращаем нулевое значение типа T и false
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	last := len(s.elements) - 1
	item := s.elements[last]
	s.elements = s.elements[:last] // убираем последний элемент
	return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
	// смотрим верхушку, но не удаляем
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	return s.elements[len(s.elements)-1], true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{elements: make([]T, 0)}
}

func main() {
	s := NewStack[int]() // стек для int

	s.Push(1)
	s.Push(2)
	s.Push(3) // сверху лежит 3 (LIFO)

	fmt.Println(s.Peek())    // 3 true, только смотрим
	fmt.Println(s.Pop())     // 3 true, сняли верхний
	fmt.Println(s.Pop())     // 2 true
	fmt.Println(s.IsEmpty()) // false, внутри ещё 1

	s.Pop()                  // сняли последний 1
	fmt.Println(s.IsEmpty()) // true
	fmt.Println(s.Pop())     // 0 false, стек пуст
}
