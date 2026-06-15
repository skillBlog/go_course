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
	words := strings.Fields(text)
	mp := make(map[string]int)
	for _, word := range words {
		mp[word]++
	}
	return mp
}

func PrintWordFrequency(mp map[string]int) {
	for word, count := range mp {
		fmt.Println(word, count)
	}
}
