package main

import (
	"errors"
	"fmt"
)

func main() {
	println("Case 1")
	case1()
	println()
	println()

	println("Case 2")
	case2()
	println()
	println()

	println("Case 3")
	case3()
	println()
	println()

}

// вывод будет такой:
// without:
// nil
// Default error
// with:
// nil
// Default error
// defer меняет локальную переменную,
// но не возвращаемое значение (результат не именованный)
func case1() {
	helperWithDefer := func(isError bool) error {
		var retVal error
		// defer выполнится после вычисления return-значения, но до возврата
		defer func() {
			retVal = errors.New("Extra error")
		}()

		if isError {
			retVal = errors.New("Default error")
		}

		// возвращается текущее значение retVal,
		// defer уже не может изменить результат
		return retVal
	}

	helperWithoutDefer := func(isError bool) error {
		var retVal error

		if isError {
			retVal = errors.New("Default error")
		}

		return retVal
	}

	fmt.Println("\twithout:")
	fmt.Println(helperWithoutDefer(false))
	fmt.Println(helperWithoutDefer(true))
	fmt.Println("\twith:")
	fmt.Println(helperWithDefer(false))
	fmt.Println(helperWithDefer(true))
}

// вывод будет такой:
// without:
// nil
// Default error
// with:
// Extra error
// Extra error
// retVal именованный результат,
// defer может изменить его перед возвратом
func case2() {
	helperWithDefer := func(isError bool) (retVal error) {
		defer func() {
			retVal = errors.New("Extra error")
		}()

		if isError {
			retVal = errors.New("Default error")
		}

		return
	}

	helperWithoutDefer := func(isError bool) (retVal error) {
		if isError {
			retVal = errors.New("Default error")
		}

		return
	}

	fmt.Println("\twithout:")
	fmt.Println(helperWithoutDefer(false))
	fmt.Println(helperWithoutDefer(true))
	fmt.Println("\twith:")
	fmt.Println(helperWithDefer(false))
	fmt.Println(helperWithDefer(true))
}

// вывод будет такой:
// without:
// nil
// Default error
// with:
// First Error
// First Error
// потому что defer выполняются в порядке стека (LIFO)
func case3() {
	helperWithDefer := func(isError bool) (retVal error) {
		defer func() {
			// выполнится последним, изменит итоговое значение
			retVal = errors.New("First Error")
		}()

		defer func() {
			// выполнится первым
			retVal = errors.New("Second Error")
		}()

		if isError {
			retVal = errors.New("Default error")
		}

		return
	}

	helperWithoutDefer := func(isError bool) (retVal error) {
		if isError {
			retVal = errors.New("Default error")
		}

		return
	}

	fmt.Println("\twithout:")
	fmt.Println(helperWithoutDefer(false))
	fmt.Println(helperWithoutDefer(true))
	fmt.Println("\twith:")
	fmt.Println(helperWithDefer(false))
	fmt.Println(helperWithDefer(true))
}
