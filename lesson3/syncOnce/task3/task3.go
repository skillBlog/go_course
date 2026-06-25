package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Plugin: интерфейс для всех плагинов
type Plugin interface {
	Execute() string
}

// PluginManager: управляет инициализацией и доступом к плагинам
type PluginManager struct {
	plugins map[string]*pluginEntry
	mu      sync.RWMutex
}

type pluginEntry struct {
	once   sync.Once
	plugin Plugin
	err    error
	initFn func() (Plugin, error)
}

// NewPluginManager создаёт новый менеджер плагинов
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*pluginEntry),
	}
}

// RegisterPlugin: регистрирует новый плагин; ошибка, если имя уже занято
func (pm *PluginManager) RegisterPlugin(name string, initFn func() (Plugin, error)) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %q already registered", name)
	}

	pm.plugins[name] = &pluginEntry{
		initFn: initFn,
	}
	return nil
}

// get: потокобезопасная однократная инициализация одного плагина
func (e *pluginEntry) get() (Plugin, error) {
	e.once.Do(func() {
		e.plugin, e.err = e.initFn()
	})
	return e.plugin, e.err
}

// GetPlugin: возвращает инициализированный плагин
func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	entry, ok := pm.plugins[name]
	pm.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("plugin %q not found", name)
	}

	return entry.get()
}

// DemoPlugin: реализация плагина
type DemoPlugin struct{}

func (p *DemoPlugin) Execute() string {
	return "DemoPlugin executed successfully!"
}

func initDemo() (Plugin, error) {
	time.Sleep(100 * time.Millisecond) // имитация длительной инициализации
	return &DemoPlugin{}, nil
}

func main() {
	pm := NewPluginManager()

	if err := pm.RegisterPlugin("demo", initDemo); err != nil {
		log.Fatal(err)
	}
	if err := pm.RegisterPlugin("broken", func() (Plugin, error) {
		return nil, fmt.Errorf("simulated error")
	}); err != nil {
		log.Fatal(err)
	}
	if err := pm.RegisterPlugin("demo", initDemo); err != nil {
		log.Printf("повторная регистрация: %v", err)
	}

	var wg sync.WaitGroup

	// тестирование рабочего плагина: 5 горутин, initDemo вызовется один раз
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			p, err := pm.GetPlugin("demo")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
				return
			}
			log.Printf("Goroutine %d: %s", id, p.Execute())
		}(i)
	}

	// тестирование плагина с ошибкой: initFn упадёт один раз, ошибка закэшируется
	for i := 5; i < 7; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			_, err := pm.GetPlugin("broken")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// Пример вывода:
	//   Goroutine 2: DemoPlugin executed successfully!
	//   Goroutine 0: DemoPlugin executed successfully!
	//   ... (ещё 3 строки для demo, порядок может отличаться)
	//   Goroutine 5 error: simulated error
	//   Goroutine 6 error: simulated error
	//
	// Почему так:
	//   1. У каждого pluginEntry свой sync.Once: demo и broken инициализируются независимо.
	//   2. 5 горутин ждут в once.Do, пока initDemo завершится; initFn вызывается ровно один раз.
	//   3. Ошибка из broken тоже кэшируется в entry.err: повторных вызовов initFn не будет.
	//   4. RWMutex защищает map при RegisterPlugin и чтении entry по имени.
	//   5. Повторный RegisterPlugin с тем же name вернёт ошибку, запись не перезапишется.
}
