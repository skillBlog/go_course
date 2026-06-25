package main

import (
	"fmt"
	"sync"
)

// bufferPool: переиспользуемые буферы []byte вместо новой аллокации на каждый вызов
var bufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 64)
		return &b
	},
}

// toUpperInPlace меняет регистр в существующем слайсе, без новой аллокации
func toUpperInPlace(buf []byte) {
	for i, b := range buf {
		if b >= 'a' && b <= 'z' {
			buf[i] = b - ('a' - 'A')
		}
	}
}

// ProcessString: преобразует строку в верхний регистр, буфер берётся из sync.Pool
func ProcessString(s string) string {
	bp := bufferPool.Get().(*[]byte)
	buf := (*bp)[:0]

	defer func() {
		*bp = buf[:0] // сброс длины, capacity сохраняется для следующего вызова
		bufferPool.Put(bp)
	}()

	buf = append(buf, s...)
	toUpperInPlace(buf)

	// string(buf) копирует данные, после return буфер можно безопасно вернуть в пул
	return string(buf)
}

func main() {
	examples := []string{
		"hello, world!",
		"gopher",
		"lorem ipsum dolor sit amet",
	}

	for _, s := range examples {
		processed := ProcessString(s)
		fmt.Printf("Original: %q\nProcessed: %q\n\n", s, processed)
	}

	// Пример вывода:
	//   Original: "hello, world!"
	//   Processed: "HELLO, WORLD!"
	//
	//   Original: "gopher"
	//   Processed: "GOPHER"
	//
	//   Original: "lorem ipsum dolor sit amet"
	//   Processed: "LOREM IPSUM DOLOR SIT AMET"
	//
	// Почему sync.Pool:
	//   1. На каждый ProcessString не создаётся новый []byte: берём из пула.
	//   2. После string(buf) буфер возвращается в пул через defer Put.
	//   3. sync.Pool потокобезопасен: можно вызывать ProcessString из разных горутин.
	//   4. Утечек нет: строка-результат — копия, буфер сбрасывается [:0] перед Put.
	//   5. toUpperInPlace меняет buf на месте; bytes.ToUpper создавал бы новый слайс.
}
