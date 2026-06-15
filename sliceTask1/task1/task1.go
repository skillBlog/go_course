package main

import "fmt"

type account struct {
	value int
}

func main() {
	s1 := make([]account, 0, 2)
	s1 = append(s1, account{})
	s2 := append(s1, account{})
	acc := &s2[0]
	acc.value = 100
	fmt.Println(s1, s2) // [{100}] [{100} {0}]
	s1 = append(s1, account{})
	acc.value += 100
	fmt.Println(s1, s2) // [{200} {0}] [{200} {0}]
}
