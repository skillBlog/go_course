package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type User struct {
	Name string
}

type CacheBase interface {
	Set(key string, value any, ttl time.Duration)
	Get(key string) (any, bool)
	Delete(key string)
	Clear()
	Exists(key string) bool
	ToJSON() ([]byte, error)
}

type CachePayload struct {
	// RWMutex: несколько читателей (RLock) или один писатель (Lock)
	mu    sync.RWMutex
	items map[string]CacheItem // value и TTL хранятся вместе
}

type CacheItem struct {
	Value any       `json:"value"`
	TTL   time.Time `json:"ttl"`
}

func (c *CachePayload) Set(key string, value any, ttl time.Duration) {
	// Lock эксклюзивная блокировка для записи в map
	c.mu.Lock()
	// Unlock снимает Lock, defer гарантирует разблокировку при любом return
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value: value,
		TTL:   time.Now().Add(ttl),
	}
}

func (c *CachePayload) Get(key string) (any, bool) {
	// Lock при истечении TTL удаляем ключ (delete) это запись
	c.mu.Lock()
	defer c.mu.Unlock()

	// получение значения из кэша по ключу
	item, ok := c.items[key]
	if !ok {
		// ключ не найден, тогда возвращаем nil и false
		return nil, false
	}

	// ленивое удаление просроченных ключей
	if time.Now().After(item.TTL) {
		delete(c.items, key)
		return nil, false
	}

	// ключ существует и не просрочен, тогда возвращаем значение
	return item.Value, true
}

func (c *CachePayload) Delete(key string) {
	// Lock удаление из map это запись, нужен эксклюзивный доступ
	c.mu.Lock()
	defer c.mu.Unlock()

	// удаляем ключ из кэша
	delete(c.items, key)
}

// типизированное получение значения из кэша.
// В Go у методов не может быть type parameters, поэтому это generic-функция
func GetAs[T any](c CacheBase, key string) (T, error) {
	var zero T

	// получаем значение через Get с проверкой TTL и ленивым удалением
	value, ok := c.Get(key)
	if !ok {
		return zero, fmt.Errorf("ключ %q не найден или просрочен", key)
	}

	// приводим any к запрошенному типу T
	typed, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("ключ %q: неверный тип, ожидался %T", key, zero)
	}

	return typed, nil
}

func (c *CachePayload) Clear() {
	// Lock пересоздание map это запись, нужен эксклюзивный доступ
	c.mu.Lock()
	defer c.mu.Unlock()

	// полная очистка кэша
	c.items = make(map[string]CacheItem)
}

func (c *CachePayload) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return false
	}

	return !time.Now().After(item.TTL)
}

func (c *CachePayload) ToJSON() ([]byte, error) {
	// RLock shared-блокировка только для чтения; несколько горутин могут читать параллельно
	c.mu.RLock()
	// RUnlock снимает RLock (не путать с Unlock, который снимает Lock)
	defer c.mu.RUnlock()

	// сериализуем актуальные записи; value и TTL уже хранятся вместе
	data := make(map[string]CacheItem)
	now := time.Now()

	for k, item := range c.items {
		if now.After(item.TTL) {
			continue
		}
		data[k] = item
	}

	return json.Marshal(data)
}

func NewCache() CacheBase {
	return &CachePayload{items: make(map[string]CacheItem)}
}

func main() {
	cache := NewCache()

	// Добавление данных с TTL
	cache.Set("user", User{Name: "Alice"}, time.Hour) // Хранится 1 час
	cache.Set("temp_data", 42, time.Minute)           // Хранится 1 минуту

	// Получение значения с проверкой TTL
	if value, ok := cache.Get("user"); ok {
		fmt.Println("Get(user):", value)
	}

	// Типизированное получение
	user, err := GetAs[User](cache, "user")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("GetAs(user):", user)
	}

	// Сериализация в JSON
	jsonData, _ := cache.ToJSON()
	fmt.Println(string(jsonData))

	// Удаление конкретного ключа
	cache.Delete("temp_data")
	fmt.Println("Exists (temp_data):", cache.Exists("temp_data")) // false

	// Очистка кэша
	cache.Clear()
	fmt.Println("Exists (user):", cache.Exists("user")) // false
}
