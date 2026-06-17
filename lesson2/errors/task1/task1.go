package main

import (
	"fmt"
)

// выведет: Default error
func main() {
	// helper возвращает error, fmt.Println выведет текст ошибки
	fmt.Println(helper())
}

func helper() error {
	// fmt.Errorf создаёт ошибку с сообщением (как errors.New, но с форматированием)
	return fmt.Errorf("Default error")
}
