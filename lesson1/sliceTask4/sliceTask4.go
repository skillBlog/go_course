package main

import (
	"fmt"
	"unsafe"
)

func main() {
	//1
	// будет выведено: first: [] : 0 : 0 потомучто мы присваиваем nil слайсу, а nil это пустой слайс
	first := []int{1, 2, 3, 4, 5}
	first = nil
	fmt.Println("first:", first, ":", len(first), ":", cap(first)) // first: [] : 0 : 0

	//2
	// будет выведено: second: [] : 0 : 5 потомучто мы присваиваем слайсу пустой слайс
	second := []int{1, 2, 3, 4, 5}
	second = second[:0]
	fmt.Println("second:", second, ":", len(second), ":", cap(second)) // second: [] : 0 : 5

	//3
	// будет выведено: third: [0 0 0 0 0] : 5 : 5 потомучто мы очищаем слайс через clear, т.е. устанавливаем все элементы в 0
	third := []int{1, 2, 3, 4, 5}
	clear(third)
	fmt.Println("third:", third, ":", len(third), ":", cap(third)) // third: [0 0 0 0 0] : 5 : 5

	//4
	// будет выведено: fourth: [1 0 0 4 5] : 5 : 5 потомучто мы очищаем слайс через clear, т.е. устанавливаем 2 и 3 элементы в 0
	fourth := []int{1, 2, 3, 4, 5}
	clear(fourth[1:3])
	fmt.Println("fourth:", fourth, ":", len(fourth), ":", cap(fourth)) // fourth: [1 0 0 4 5] : 5 : 5

	//5
	// будет выведено: slice = [10 0 0] 3 6 потомучто мы создаем слайс длиной 3 и вместимостью 6, и присваиваем его массиву
	// и устанавливаем 0 элемент в 10
	// и массив тоже будет [0 0 0] 3 3 потомучто мы присваиваем слайсу массив, а массив это статический массив
	slice := make([]int, 3, 6)
	array := [3]int(slice[:3])
	slice[0] = 10

	fmt.Println("slice = ", slice, len(slice), cap(slice)) // slice = [10 0 0] 3 6
	fmt.Println("array =", array, len(array), cap(array))  // array =[0 0 0] 3 3

	//6 В каких случаях Slice пустой или нулевой
	//1
	// будет выведено: var data []string: empty=true nil=false size=24 data=0xc000018120 потомучто мы создаем пустой слайс
	var data []string
	fmt.Println("var data []string:")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data))
	//2
	// будет выведено: data = []string(nil): empty=true nil=true size=24 data=0xc000018120 потомучто мы присваиваем nil слайсу
	data = []string(nil)
	fmt.Println("data = []string(nil):")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data))
	//3
	// будет выведено: data = []string{} empty=true nil=false size=24 data=0xc000018120 потомучто мы присваиваем пустой слайс
	data = []string{}
	fmt.Println("data = []string{}")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data))
	//4
	// будет выведено: data =make([]string,0) empty=true nil=false size=24 data=0xc000018120 потомучто мы создаем слайс длиной 0
	data = make([]string, 0)
	fmt.Println("data =make([]string,0)")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data))

	empty := struct{}{}
	fmt.Println("empty struct address ", unsafe.Pointer(&empty))
}
