package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "golang is great and golang is fast"
	mp := WordFrequency(text)
	PrintWordFrequency(mp)
}

func WordFrequency(text string) map[string]int {
	// strings.Fields разбивает текст на слова по пробелам и переносам строк, возвращает слайс слов
	words := strings.Fields(text)
	mp := make(map[string]int)
	for _, word := range words {
		// приводим к нижнему регистру
		word = strings.ToLower(word)
		// если слова ещё нет в мапе, то будет 0, затем увеличиваем на 1
		mp[word]++
	}
	return mp
}

func PrintWordFrequency(mp map[string]int) {
	// проходим по мапе и выводим ключ + слово, значение - сколько раз встретилось
	for word, count := range mp {
		fmt.Println(word, count)
	}
}
