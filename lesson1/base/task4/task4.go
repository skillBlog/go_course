package task4

import "fmt"

type MyStruct struct {
	MyInt int
}

func func1() MyStruct {
	return MyStruct{MyInt: 1}
}

func func2() *MyStruct {
	return &MyStruct{}
}

func func3(s *MyStruct) {
	s.MyInt = 333
}

func func4(s MyStruct) {
	s.MyInt = 923
}

func func5() *MyStruct {
	return nil
}

func main() {
	// func1 возвращает структуру по значению, копия с {MyInt: 1}
	ms1 := func1()
	fmt.Println(ms1.MyInt) // 1

	// func2 возвращает указатель на новую структуру; поля по умолчанию — нули
	ms2 := func2()
	fmt.Println(ms2.MyInt) // 0

	// func3 принимает указатель и меняет оригинал, на который ссылается ms2
	func3(ms2)
	fmt.Println(ms2.MyInt) // 333

	// func4 принимает копию структуры и внутри меняется только копия, ms1 не трогается
	func4(ms1)
	fmt.Println(ms1.MyInt) // 1

	// func5 возвращает nil; обращение к полю nil указателя вызывает panic
	ms5 := func5()
	fmt.Println(ms5.MyInt) // panic
}
