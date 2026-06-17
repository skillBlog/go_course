package main

import (
	"errors"
	"fmt"
)

var ErrInvalidAge = errors.New("возраст недопустим")

func main() {
	fmt.Println(SimpleError())      // простая ошибка
	fmt.Println(FormattedError(15)) // ошибка: возраст 15 недопустим: возраст недопустим
	fmt.Println(StructError())      // не найдено
}

// возвращает простую ошибку
func SimpleError() error {
	return errors.New("простая ошибка")
}

func FormattedError(age int) error {
	return fmt.Errorf("ошибка: возраст %d недопустим: %w", age, ErrInvalidAge)
}

type MyError struct {
	Code int
	Msg  string
}

func (e MyError) Error() string {
	return e.Msg
}

func StructError() error {
	return MyError{Code: 404, Msg: "не найдено"}
}
