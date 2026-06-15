package main

import (
	"fmt"
	"slices"
)

// RemoveUnordered удаляет элемент по индексу без сохранения порядка (swap с последним).
func RemoveUnordered[T any](s []T, i int) []T {
	if i < 0 || i >= len(s) {
		return s
	}
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// RemoveOrdered удаляет элемент по индексу с сохранением порядка.
func RemoveOrdered[T any](s []T, i int) []T {
	if i < 0 || i >= len(s) {
		return s
	}
	return slices.Delete(s, i, i+1)
}

// RemoveAllByValue удаляет все вхождения value из слайса.
func RemoveAllByValue[T comparable](s []T, value T) []T {
	return slices.DeleteFunc(s, func(t T) bool {
		return t == value
	})
}

// RemoveDuplicates оставляет только уникальные элементы, сохраняя порядок
// первых вхождений. Например: [1, 2, 1, 3] -> [1, 2, 3].
// slices.Compact здесь не подходит: он удаляет только соседние дубликаты.
func RemoveDuplicates[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	result := s[:0]

	for _, v := range s {
		if _, exists := seen[v]; exists {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}

	return result
}

// RemoveIf удаляет элементы, для которых predicate возвращает true.
func RemoveIf[T any](s []T, predicate func(T) bool) []T {
	result := s[:0]
	for _, v := range s {
		if predicate(v) {
			continue
		}
		result = append(result, v)
	}
	return result
}

// RemoveOrderedWithNil удаляет элемент из слайса указателей с сохранением порядка.
// s[i] обнуляется до удаления, иначе ссылка в backing array может удерживать объект в памяти.
func RemoveOrderedWithNil[T any](s []*T, i int) []*T {
	if i < 0 || i >= len(s) {
		return s
	}
	s[i] = nil
	return slices.Delete(s, i, i+1)
}

// ShrinkCapacity уменьшает capacity, если cap > 2*len.
// Иначе возвращает слайс без изменений. Лишняя ёмкость не оправдана.
func ShrinkCapacity[T any](s []T) []T {
	if cap(s) <= 2*len(s) {
		return s
	}
	// cap = len: копируем данные в компактный массив, старый может быть собран GC.
	trimmed := make([]T, len(s))
	copy(trimmed, s)
	return trimmed
}

func main() {
	slice := []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}
	fmt.Println(slice)
	fmt.Println(RemoveUnordered(slice, 2))
	fmt.Println(RemoveOrdered(slice, 2))
	fmt.Println(RemoveAllByValue(slice, 2))
	fmt.Println(RemoveDuplicates(slice))
}
