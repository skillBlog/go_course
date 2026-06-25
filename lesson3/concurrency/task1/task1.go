package main

import (
	"fmt"
	"sync"
)

// map не потокобезопасна, поэтому доступ защищён мьютексом
type SafeCache struct {
	mu   sync.RWMutex // RWMutex - много читателей или один писатель
	data map[string]string
}

// чтение из кеша: RLock позволяет нескольким горутинам читать одновременно
func (c *SafeCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	return value, ok // ok=false, если ключ не найден
}

// запись в кеш: Lock блокирует и читателей, и других писателей
func (c *SafeCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// NewSafeCache создаёт потокобезопасный кеш с инициализированной map
func NewSafeCache() *SafeCache {
	return &SafeCache{
		data: make(map[string]string),
	}
}

func main() {
	cache := NewSafeCache()

	cache.Set("key", "value")
	value, ok := cache.Get("key")
	fmt.Println(value, ok)

	// вывод: value true
	//   value - строка, которую положили в Set
	//   true  - ключ key найден в кеше (значение ok из Get)
}
