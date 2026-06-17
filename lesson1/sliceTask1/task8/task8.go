package main

import (
	"fmt"
)

func main() {
	// len=1, cap=3, один элемент со значением по умолчанию 0
	nums := make([]int, 1, 3)
	fmt.Println(nums) // [0]

	// slice передаётся копией, append внутри функции не меняет nums
	appendSlice(nums, 1)
	fmt.Println(nums) // [0]

	// copy копирует min(len(nums), len(src)) = 1 элемент, nums[0]=2
	copySlice(nums, []int{2, 3})
	fmt.Println(nums) // [2]

	// len(nums)=1, индекс 1 не существует, panic
	mutateSlice(nums, 1, 4)
	fmt.Println(nums) // panic
}

func appendSlice(sl []int, val int) {
	sl = append(sl, val)
}

func copySlice(sl, cp []int) {
	copy(sl, cp)
}

func mutateSlice(sl []int, idx, val int) {
	sl[idx] = val
}
