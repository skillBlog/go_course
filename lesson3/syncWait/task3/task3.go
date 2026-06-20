package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type logic struct{}

var Logic logic

func (l *logic) UpdateDB(ctx context.Context, item *Item) error {
	return nil // заглушка
}

func (l *logic) FetchItems(ctx context.Context) ([]*Item, error) {
	return []*Item{
		{Value: 5},
		{Value: 15},
		{Value: 7},
	}, nil // заглушка
}

type Item struct {
	Value int
}

func processItem(item *Item) {
	time.Sleep(time.Second)
	if item.Value > 10 {
		fmt.Printf("ERROR: item %d can't be more than 10\n", item.Value)
		return
	}

	err := Logic.UpdateDB(context.Background(), item)
	if err != nil {
		fmt.Println("ERROR: can't process item")
	}
}

// Сломанный вариант (main не ждёт горутины):
//
//	func DoBusinessLogic() error {
//		items, err := Logic.FetchItems(context.Background())
//		if err != nil {
//			return err
//		}
//
//		for _, item := range items {
//			go processItem(item)
//		}
//
//		return nil
//	}

func DoBusinessLogic() error {
	items, err := Logic.FetchItems(context.Background())
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		go func(it *Item) {
			defer wg.Done()
			processItem(it)
		}(item)
	}
	wg.Wait()

	return nil
}

func main() {
	err := DoBusinessLogic()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("All items processed")

	// Исправление: WaitGroup + wg.Wait() в DoBusinessLogic (см. сломанный вариант выше).
	//
	// Пример вывода (около 1 сек, не около 0):
	//   ERROR: item 15 can't be more than 10
	//   All items processed
	//
	// "All items processed" всегда после завершения всех трёх processItem.
}
