package main

import (
	"fmt"
	"slices"
)

func FilterByValue(m map[int]string, allowedValues []string) map[int]string {
	newMap := make(map[int]string)
	for key, value := range m {
		if slices.Contains(allowedValues, value) {
			newMap[key] = value
		}
	}
	return newMap
}

func InvertMap(m map[int]string) map[string]int {
	newMap := make(map[string]int)
	for key, value := range m {
		newMap[value] = key
	}
	return newMap
}

func main() {
	m := map[int]string{1: "a", 2: "b", 3: "c", 4: "d", 5: "e"}
	allowedValues := []string{"a", "b", "e"}
	// результат: map[1:a 2:b 5:e], c и d отфильтрованы
	fmt.Println(FilterByValue(m, allowedValues))
	// меняем ключи и значения местами, результат: map[a:1 b:2 e:5]
	fmt.Println(InvertMap(FilterByValue(m, allowedValues)))
}
