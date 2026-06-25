package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// Исходный код с ошибками
//
// func main() {
// 	alreadyStored := make(map[int]struct{})
// 	capacity := 1000
// 	doubles := make([]int, 0, capacity)
// 	for i := 0; i < capacity; i++ {
// 		doubles = append(doubles, rand.Intn(10))
// 	}
// 	uniqueIDs := make(chan int, capacity)
// 	wg := sync.WaitGroup{}
// 	for i := 0; i < capacity; i++ {
// 		i := i
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			// data race — несколько горутин одновременно читают и пишут
// 			// в map alreadyStored, map в Go не потокобезопасна
// 			if _, ok := alreadyStored[doubles[i]]; !ok {
// 				alreadyStored[doubles[i]] = struct{}{}
// 				uniqueIDs <- doubles[i]
// 			}
// 		}()
// 	}
// 	wg.Wait()
// 	// канал uniqueIDs не закрывается — for range зависнет
// 	for val := range uniqueIDs {
// 		fmt.Println(val)
// 	}
// 	// fmt.Println(uniqueIDs) выводит адрес канала, а не данные
// 	fmt.Println(uniqueIDs)
// }

const (
	totalCount    = 1000       // сколько чисел генерируем
	valueRange    = 10         // rand.Intn(valueRange) -> 0..9
	maxUniqueVals = valueRange // не больше valueRange уникальных значений
)

// параллельная дедупликация чисел, чтобы из 1000 чисел получить уникальные
func main() {
	// мьютекс защищает map от одновременного доступа из горутин
	var mu sync.Mutex
	alreadyStored := make(map[int]struct{}, maxUniqueVals)
	doubles := make([]int, 0, totalCount)
	for i := 0; i < totalCount; i++ {
		doubles = append(doubles, rand.Intn(valueRange))
	}
	// буфер = maxUniqueVals: читатель стартует только после wg.Wait(),
	// иначе отправка уникального значения заблокируется и будет deadlock
	uniqueIDs := make(chan int, maxUniqueVals)
	wg := sync.WaitGroup{}
	for i := 0; i < totalCount; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			if _, ok := alreadyStored[doubles[i]]; !ok {
				alreadyStored[doubles[i]] = struct{}{}
				mu.Unlock()
				uniqueIDs <- doubles[i]
				return
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	// закрываем канал, чтобы for range завершился
	close(uniqueIDs)

	unique := make([]int, 0, maxUniqueVals)
	for val := range uniqueIDs {
		fmt.Println(val)
		unique = append(unique, val)
	}
	// выводим собранные значения, а не адрес канала
	fmt.Println(unique)
}
