package main

import (
	"fmt"
)

type errorString struct {
	s string
}

func (e errorString) Error() string {
	return e.s
}

// внутри интерфейса error хранится ссылка на конкретный тип (*errorString) и само значение
// Это не часть объявления интерфейса, а часть его реализации в рантайме
// Официальная документация описывает это так:
// "An interface value can be thought of as a pair: a concrete value and that value's type."
func checkErr(err error) {
	// Интерфейс error внутри хранит две вещи:
	//   1. конкретный тип
	//   2. конкретное значение
	//
	// err == nil вернёт true только если оба поля nil:
	//   type  = nil
	//   value = nil
	//
	// Если type уже известен, то интерфейс считается ненулевым,
	// даже если value == nil
	fmt.Println(err == nil)
}

func main() {
	var e1 error

	// e1 уже имеет тип error и равен nil.
	// В checkErr попадёт nil интерфейс и поэтому выведется true
	checkErr(e1) // true

	var e *errorString
	// e == nil, но checkErr принимает error
	// Поэтому Go создаёт интерфейс error и кладёт в него:
	//   тип: *errorString
	//   значение: nil
	//
	// Такой интерфейс уже не равен nil,
	// потому что внутри есть информация о типе и поэтому выведется false
	checkErr(e) // false

	e = &errorString{}
	// В интерфейсе есть и тип, и значение.
	// Поэтому выведется false
	checkErr(e) // false

	e = nil
	// Снова e == nil, но при передаче в error
	// создаётся интерфейс:
	//   тип: *errorString
	//   значение: nil
	//
	// Поэтому результат опять false
	checkErr(e) // false
}
