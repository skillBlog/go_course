package main

import (
	"fmt"
	"sync"
	"time"
)

// Connection: имитация подключения к БД
type Connection struct {
	ID string
}

// Database: хранит подключение и sync.Once для однократной инициализации
type Database struct {
	once sync.Once
	conn *Connection
}

// initConnection: выполняется ровно один раз внутри once.Do
func (d *Database) initConnection() {
	fmt.Println("инициализация подключения к БД...")
	time.Sleep(100 * time.Millisecond) // имитация долгого подключения
	d.conn = &Connection{ID: "db-main"}
	fmt.Println("подключение создано")
}

// GetConnection: при первом вызове создаёт conn, дальше возвращает тот же указатель
func (d *Database) GetConnection() *Connection {
	d.once.Do(d.initConnection)
	return d.conn
}

func main() {
	var db Database
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			conn := db.GetConnection()
			fmt.Printf("горутина %d: получила подключение %s\n", id, conn.ID)
		}(i)
	}

	wg.Wait()

	// Пример вывода:
	//   инициализация подключения к БД...
	//   подключение создано
	//   горутина 3: получила подключение db-main
	//   горутина 7: получила подключение db-main
	//   горутина 1: получила подключение db-main
	//   ... (ещё 7 строк, порядок горутин может отличаться)
	//
	// Почему так:
	//   1. 10 горутин одновременно вызывают GetConnection().
	//   2. sync.Once.Do гарантирует: initConnection выполнится только один раз.
	//      Остальные горутины заблокируются внутри Do, пока инициализация не завершится.
	//   3. Все 10 получат один и тот же d.conn (один указатель, ID "db-main").
	//   4. Строки "инициализация..." и "подключение создано" напечатаются ровно по одному разу.
}
