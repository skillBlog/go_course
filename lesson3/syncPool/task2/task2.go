package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// RequestData: данные JSON запроса с вложенными map и slice
type RequestData struct {
	UserID string            `json:"user_id"`
	Action string            `json:"action"`
	Items  []string          `json:"items"`
	Labels map[string]string `json:"labels"`
}

// Reset: очищает объект перед возвратом в пул: без утечки данных прошлого запроса
func (rd *RequestData) Reset() {
	rd.UserID = ""
	rd.Action = ""
	rd.Items = rd.Items[:0]
	for k := range rd.Labels {
		delete(rd.Labels, k)
	}
}

// requestDataPool: переиспользуем RequestData вместо new() на каждый запрос
var requestDataPool = sync.Pool{
	New: func() any {
		return &RequestData{
			Items:  make([]string, 0, 8),
			Labels: make(map[string]string),
		}
	},
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := requestDataPool.Get().(*RequestData)
	defer func() {
		data.Reset() // один раз: очищаем перед возвратом в пул
		requestDataPool.Put(data)
	}()

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user_id":     data.UserID,
		"action":      data.Action,
		"items_count": len(data.Items),
		"labels":      data.Labels,
	})
}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server started at :8080")
	fmt.Println(`Пример: curl -X POST localhost:8080 -d '{"user_id":"1","action":"login","items":["a","b"],"labels":{"env":"dev"}}'`)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("server error:", err)
	}

	// как работает sync.Pool здесь:
	//   1. Get() берёт RequestData с уже созданными Items (cap 8) и Labels map.
	//   2. Объект из пула уже чистый: предыдущий запрос вызвал Reset() в defer перед Put().
	//   3. defer Reset() + Put(): очищаем после обработки, буферы map/slice переиспользуются.
	//   4. sync.Pool потокобезопасен при конкурентных HTTP-запросах.
}
