package main

import (
	"fmt"
)

type MyError struct {
	data string
}

func (m *MyError) Error() string {
	return m.data
}

func foo(i int) error {
	var err *MyError
	if i > 5 {
		err = &MyError{data: "i>5"}
	}
	// В интерфейс error попадёт:
	//   type = *MyError
	//   value = nil
	// Поэтому возвращаемый error не равен nil,
	// хотя сам указатель err равен nil
	return err
}

// чтобы это исправить, нужно вернуть не указатель на MyError, а MyError
// func foo(i int) error {
// 	if i > 5 {
// 		return &MyError{data: "i>5"}
// 	}
// 	return nil
// }

func main() {
	err := foo(4)
	// выведится "oops", т.к. возвращаемый error не равен nil
	// но сам указатель err равен nil
	if err != nil {
		fmt.Println("oops")
	} else {
		fmt.Println("ok")
	}
}
