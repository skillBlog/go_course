package main

import "fmt"

func CausePanic() {
	panic("что-то пошло не так!")
}

func HandlePanic() {
	defer func() {
		// recover работает только внутри defer и возвращает значение panic или nil, если паники не было
		if r := recover(); r != nil {
			fmt.Printf("Паника перехвачена: %v\n", r)
		}
	}()
	// Паника поднимается по стеку, но recover в defer выше перехватывает её, поэтому программа не завершается.
	CausePanic()
}

func main() {
	// CausePanic() // без recover программа аварийно завершится
	HandlePanic()
	fmt.Println("Программа продолжает работу")
}
