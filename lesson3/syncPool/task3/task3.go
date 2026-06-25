package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type cacheItem struct {
	value     any
	expiresAt time.Time
}

// ObjectCache: потокобезопасный кэш с TTL и JSON сериализацией через sync.Pool
type ObjectCache struct {
	mu       sync.RWMutex
	items    map[string]cacheItem
	ttl      time.Duration
	jsonPool sync.Pool
	stopCh   chan struct{}
}

func NewObjectCache(ttl time.Duration) *ObjectCache {
	c := &ObjectCache{
		items:  make(map[string]cacheItem),
		ttl:    ttl,
		stopCh: make(chan struct{}),
	}
	c.jsonPool.New = func() any {
		b := make([]byte, 0, 256)
		return &b
	}
	go c.cleanupLoop()
	return c
}

// cleanupLoop: фоновая горутина удаляет устаревшие записи
func (c *ObjectCache) cleanupLoop() {
	interval := c.ttl / 2
	if interval < time.Second {
		interval = time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stopCh:
			return
		}
	}
}

func (c *ObjectCache) removeExpired() {
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if now.After(item.expiresAt) {
			delete(c.items, key)
		}
	}
}

// Set: добавляет объект в кэш с TTL
func (c *ObjectCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Get: возвращает объект, если ключ есть и TTL не истёк
func (c *ObjectCache) Get(key string) (any, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}
	return item.value, true
}

// Delete: удаляет объект по ключу
func (c *ObjectCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// putJSONBuffer возвращает слайс в пул.
// bytes.Buffer при превышении cap(*bp) выделяет новый backing array, pooled слайс
// на него уже не ссылается, сохраняем увеличенную ёмкость, иначе пул всегда остаётся 256 байт.
func (c *ObjectCache) putJSONBuffer(bp *[]byte, buf *bytes.Buffer) {
	if buf.Cap() > cap(*bp) {
		*bp = make([]byte, 0, buf.Cap())
	} else {
		*bp = (*bp)[:0]
	}
	c.jsonPool.Put(bp)
}

// ToJSON: сериализует актуальные записи кэша, буфер берётся из sync.Pool
func (c *ObjectCache) ToJSON() ([]byte, error) {
	c.mu.RLock()
	snapshot := make(map[string]any, len(c.items))
	now := time.Now()
	for key, item := range c.items {
		if now.Before(item.expiresAt) {
			snapshot[key] = item.value
		}
	}
	c.mu.RUnlock()

	bp := c.jsonPool.Get().(*[]byte)
	buf := bytes.NewBuffer((*bp)[:0])

	if err := json.NewEncoder(buf).Encode(snapshot); err != nil {
		c.putJSONBuffer(bp, buf)
		return nil, err
	}

	out := append([]byte(nil), buf.Bytes()...)
	c.putJSONBuffer(bp, buf)
	return out, nil
}

func main() {
	cache := NewObjectCache(5 * time.Second)

	cache.Set("user:1", map[string]string{"name": "Alice", "role": "admin"})
	cache.Set("user:2", map[string]string{"name": "Bob", "role": "user"})

	if user, found := cache.Get("user:1"); found {
		fmt.Println("Найден:", user)
	}

	jsonData, err := cache.ToJSON()
	if err != nil {
		fmt.Println("ошибка JSON:", err)
		return
	}
	fmt.Println("Кэш в JSON:", string(jsonData))

	time.Sleep(6 * time.Second)
	_, found := cache.Get("user:1")
	fmt.Println("После TTL, user:1 найден?", found)

	// Пример вывода:
	//   Найден: map[name:Alice role:admin]
	//   Кэш в JSON: {"user:1":{"name":"Alice","role":"admin"},"user:2":{"name":"Bob","role":"user"}}
	//   После TTL, user:1 найден? false
	//
	// Почему так:
	//   1. map + RWMutex: Get читает под RLock, Set/Delete под Lock.
	//   2. cleanupLoop в фоне периодически удаляет просроченные ключи.
	//   3. Get тоже проверяет TTL: даже без фоновой очистки вернёт false.
	//   4. ToJSON использует sync.Pool для буфера сериализации вместо новой аллокации каждый раз.
	//   5. Если JSON больше cap pooled слайса, Buffer реллоцируется; putJSONBuffer сохраняет
	//      новую ёмкость в пуле, иначе рост буфера терялся бы при каждом Put.
}
