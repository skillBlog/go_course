package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// результат обработки одного файла, передается из горутины в main через канал
type wordCountResult struct {
	filename string
	count    int
	err      error
}

// читает файл и считает слова через strings.Fields
func countWords(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return len(strings.Fields(string(data))), nil
}

// паттерн fan-out: одна задача (список файлов) распределяется
// на несколько горутин воркеров. Каждый файл обрабатывается в отдельной горутине
// результаты собираются в main через for range по каналу results
func fanOut(filePaths []string) <-chan wordCountResult {
	// буфер = len(filePaths): горутины не блокируются на send, пока main не читает
	results := make(chan wordCountResult, len(filePaths))
	var wg sync.WaitGroup

	for _, path := range filePaths {
		wg.Add(1)
		// path передается аргументом, чтобы каждая горутина получила свой файл
		go func(p string) {
			defer wg.Done()
			count, err := countWords(p)
			results <- wordCountResult{
				filename: filepath.Base(p),
				count:    count,
				err:      err,
			}
		}(path)
	}

	// отдельная горутина закрывает results, когда все воркеры завершились
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// ищет папку textFiles относительно расположения task5.go,
// иначе при запуске из корня репозитория или из редактора путь textFiles не найдется
func resolveFilesDir() string {
	candidates := []string{
		filepath.Join(filepath.Dir(sourceFileDir()), "textFiles"),
		"textFiles",
		filepath.Join("lesson3", "channels", "task5", "textFiles"),
	}
	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}
	return "textFiles"
}

func sourceFileDir() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}

func main() {
	filesDir := resolveFilesDir()

	// читаем список файлов из директории (не содержимое, только имена)
	entries, err := os.ReadDir(filesDir)
	if err != nil {
		fmt.Println("ошибка чтения директории:", err)
		return
	}

	var filePaths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filePaths = append(filePaths, filepath.Join(filesDir, entry.Name()))
	}

	if len(filePaths) == 0 {
		fmt.Println("файлы не найдены")
		return
	}

	// fan-out: параллельная обработка + агрегация в main
	// содержимое textFiles/ и ожидаемый подсчет слов:
	//   file1.txt - "hello world from file one" + "two words here"           8 слов
	//   file2.txt - "go is awesome for concurrency patterns"                 6 слов
	//   file3.txt - "fan out distributes work across goroutines"             6 слов
	//
	// вывод (порядок строк с файлами не фиксирован, горутины завершаются параллельно):
	//
	//   file2.txt: 6 слов
	//   file1.txt: 8 слов
	//   file3.txt: 6 слов
	//   ---
	//   файлов обработано: 3
	//   всего слов: 20
	//
	// Итог всего слов: 20 всегда одинаковый; порядок строк file1/file2/file3 может меняться
	totalWords := 0
	for result := range fanOut(filePaths) {
		if result.err != nil {
			fmt.Printf("%s: ошибка — %v\n", result.filename, result.err)
			continue
		}
		fmt.Printf("%s: %d слов\n", result.filename, result.count)
		totalWords += result.count
	}

	fmt.Println("---")
	fmt.Printf("файлов обработано: %d\n", len(filePaths))
	fmt.Printf("всего слов: %d\n", totalWords)
}
