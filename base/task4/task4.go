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
	ms1 := func1()
	fmt.Println(ms1.MyInt) // 1

	ms2 := func2()
	fmt.Println(ms2.MyInt) // 0

	func3(ms2)
	fmt.Println(ms2.MyInt) // 333

	func4(ms1)
	fmt.Println(ms1.MyInt) // 1

	ms5 := func5()
	fmt.Println(ms5.MyInt) // panic
}
