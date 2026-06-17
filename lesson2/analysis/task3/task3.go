package main

import "fmt"

type CustomError struct {
	message string
}

func (e *CustomError) Error() string {
	return e.message
}

func returnError(flag bool) error {
	if flag {
		return &CustomError{"Что-то пошло не так"}
	}
	// err == nil, но при возврате как error
	// значение заворачивается в интерфейс.
	// Интерфейс хранит (тип, значение),
	// поэтому (*CustomError, nil) != nil.
	var err *CustomError
	return err
}

// Чтобы исправить, нужно вернуть nil интерфейс error,
// а не nil указатель *CustomError.
// func returnError(flag bool) error {
// 	if flag {
// 		return &CustomError{"Что-то пошло не так"}
// 	}
// 	return nil
// }

func main() {
	err1 := returnError(true)
	err2 := returnError(false)

	// err1 содержит реальную ошибку (*CustomError),
	// поэтому интерфейс error не равен nil, поэтому выведется false
	fmt.Println("err1 == nil:", err1 == nil) // false

	// err2 был создан из nil указателя *CustomError.
	// Но после обертки в интерфейс error внутри остаётся тип
	// CustomError, поэтому интерфейс тоже не равен nil, поэтому выведется false
	fmt.Println("err2 == nil:", err2 == nil) // false
}
