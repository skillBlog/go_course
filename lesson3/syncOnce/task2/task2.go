package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ConfigManager: потокобезопасный менеджер конфигурации с ленивой загрузкой
type ConfigManager struct {
	once   sync.Once
	config map[string]string
}

// loadConfig: выполняется ровно один раз внутри once.Do
func (cm *ConfigManager) loadConfig() {
	fmt.Println("загрузка конфигурации...")
	time.Sleep(50 * time.Millisecond) // имитация чтения из файла или БД

	// env переопределяет значения по умолчанию (имитация гибких источников)
	cm.config = map[string]string{
		"app_name":  envOrDefault("APP_NAME", "MyApp"),
		"port":      envOrDefault("PORT", "8080"),
		"log_level": envOrDefault("LOG_LEVEL", "debug"),
	}
	fmt.Println("конфигурация загружена")
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// LoadConfig: загружает конфигурацию один раз и сохраняет в памяти
func (cm *ConfigManager) LoadConfig() {
	cm.once.Do(cm.loadConfig)
}

// Get: возвращает значение по ключу и ok (как у map); при первом вызове триггерит загрузку
func (cm *ConfigManager) Get(key string) (string, bool) {
	cm.LoadConfig()
	value, ok := cm.config[key]
	return value, ok
}

// PrintConfig: выводит все загруженные параметры
func (cm *ConfigManager) PrintConfig() {
	cm.LoadConfig()
	for key, value := range cm.config {
		fmt.Printf("%s = %s\n", key, value)
	}
}

func main() {
	var configManager ConfigManager
	var wg sync.WaitGroup

	// несколько горутин одновременно обращаются к конфигу
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			appName, ok := configManager.Get("app_name")
			if !ok {
				fmt.Printf("горутина %d: ключ app_name не найден\n", id)
				return
			}
			fmt.Printf("горутина %d: app_name = %s\n", id, appName)
		}(i)
	}

	wg.Wait()
	fmt.Println("--- все параметры ---")
	configManager.PrintConfig()

	// Пример вывода:
	//   загрузка конфигурации...
	//   конфигурация загружена
	//   горутина 2: app_name = MyApp
	//   горутина 0: app_name = MyApp
	//   ... (ещё 3 строки, порядок горутин может отличаться)
	//   все параметры:
	//   app_name = MyApp
	//   port = 8080
	//   log_level = debug
	//
	// Почему так:
	//   1. Первый вызов Get/LoadConfig/PrintConfig запускает loadConfig через once.Do
	//   2. Остальные горутины ждут внутри Do, пока загрузка не завершится.
	//   3. "загрузка..." и "конфигурация загружена" печатаются ровно один раз.
	//   4. Все горутины читают один и тот же map в памяти, гонок нет.
}
