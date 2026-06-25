package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

const (
	envFilesDir    = "TEXT_FILES_DIR"
	envWorkerCount = "WORKER_COUNT"
	defaultWorkers = 4
)

// результат обработки одного файла, передается из горутины в main через канал
type wordCountResult struct {
	filename string
	count    int
	err      error
}

// потоковое чтение: не загружаем весь файл в память, подходит для больших файлов
func countWords(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024) // увеличенный лимит токена для длинных слов

	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

// fan-out с фиксированным пулом воркеров: N горутин обрабатывают все файлы из очереди jobs
func fanOut(filePaths []string, workers int) <-chan wordCountResult {
	if workers > len(filePaths) {
		workers = len(filePaths)
	}
	if workers < 1 {
		workers = 1
	}

	jobs := make(chan string, len(filePaths))
	results := make(chan wordCountResult, len(filePaths))

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for path := range jobs {
				count, err := countWords(path)
				results <- wordCountResult{
					filename: filepath.Base(path),
					count:    count,
					err:      err,
				}
			}
		}()
	}

	go func() {
		for _, path := range filePaths {
			jobs <- path
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	return results
}

// TEXT_FILES_DIR - явный путь к папке с файлами, иначе ищем относительно cwd
func resolveFilesDir() string {
	if dir := os.Getenv(envFilesDir); dir != "" {
		return dir
	}
	candidates := []string{
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

// WORKER_COUNT - размер пула, по умолчанию runtime.NumCPU()
func resolveWorkerCount() int {
	if s := os.Getenv(envWorkerCount); s != "" {
		n, err := strconv.Atoi(s)
		if err == nil && n > 0 {
			return n
		}
	}
	if n := runtime.NumCPU(); n > 0 {
		return n
	}
	return defaultWorkers
}

func main() {
	filesDir := resolveFilesDir()

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

	// fan-out: фиксированный пул воркеров + агрегация в main
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
	for result := range fanOut(filePaths, resolveWorkerCount()) {
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
