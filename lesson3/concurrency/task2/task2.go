package main

import (
	"fmt"
	"sync"
	"time"
)

// комментарий из БД: UserID нужен, чтобы потом загрузить автора
type Comment struct {
	ID     int
	UserID int
	Text   string
}

// данные сессии пользователя: SessionID нужен для загрузки вложений
type Session struct {
	SessionID string // если пусто, вложения не грузим
}

type User struct {
	ID   int
	Name string
}

// loadComments: имитирует запрос комментариев к БД (sleep = задержка сети)
func loadComments() []Comment {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("[этап 1] комментарии загружены")
	return []Comment{
		{1, 101, "первый комментарий"},
		{2, 102, "второй комментарий"},
	}
}

// loadSession: имитирует загрузку сессии: идёт параллельно с комментариями
func loadSession() Session {
	time.Sleep(80 * time.Millisecond) // чуть быстрее комментариев, для наглядности
	fmt.Println("[этап 1] сессия загружена")
	return Session{SessionID: "sess-abc123"}
}

// loadUsers: грузит пользователей по списку комментариев: вызывается только когда комментарии уже получены
func loadUsers(comments []Comment) []User {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("[этап 2] пользователи загружены (после комментариев)")
	users := make([]User, 0, len(comments))
	seen := make(map[int]struct{})
	for _, c := range comments {
		if _, ok := seen[c.UserID]; ok {
			continue // один пользователь один раз, даже если несколько комментариев
		}
		seen[c.UserID] = struct{}{}
		users = append(users, User{ID: c.UserID, Name: fmt.Sprintf("user-%d", c.UserID)})
	}
	return users
}

// loadAttachments имитирует загрузку вложений по sessionID
func loadAttachments(sessionID string) []string {
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("[этап 3] вложения загружены для session-id: %s\n", sessionID)
	return []string{"report.pdf", "photo.png"}
}

func main() {
	var wg sync.WaitGroup

	commentsCh := make(chan []Comment, 1)
	sessionCh := make(chan Session, 1)

	// этап 1a: комментарии
	// стартует сразу, не ждёт сессию, они независимы
	wg.Add(1)
	go func() {
		defer wg.Done()
		commentsCh <- loadComments()
	}()

	// этап 1b: сессия
	// стартует сразу, работает параллельно с комментариями
	wg.Add(1)
	go func() {
		defer wg.Done()
		sessionCh <- loadSession()
	}()

	// этап 2: пользователи
	// блокируется на <-commentsCh: не начнёт, пока комментарии не придут
	wg.Add(1)
	go func() {
		defer wg.Done()
		comments := <-commentsCh
		users := loadUsers(comments)
		for _, u := range users {
			fmt.Printf("  -> user %d (%s)\n", u.ID, u.Name)
		}
	}()

	// этап 3: вложения
	// ждёт сессию: если sessionID пустой, выходим; иначе грузим один раз в этой горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		session := <-sessionCh
		if session.SessionID == "" {
			fmt.Println("[этап 3] вложения пропущены — нет session-id")
			return
		}
		attachments := loadAttachments(session.SessionID)
		for _, a := range attachments {
			fmt.Printf("  → вложение: %s\n", a)
		}
	}()

	// main ждёт завершения всех 4 горутин, программа не выйдет раньше времени
	wg.Wait()

	// пример реального вывода и почему так
	//
	//   [этап 1] сессия загружена
	//   [этап 1] комментарии загружены
	//   [этап 3] вложения загружены для session-id: sess-abc123
	//     -> вложение: report.pdf
	//     -> вложение: photo.png
	//   [этап 2] пользователи загружены (после комментариев)
	//     -> user 101 (user-101)
	//     -> user 102 (user-102)
	//
	// Почему сессия раньше комментариев:
	//   обе горутины стартуют сразу, но sleep у сессии 80ms, у комментариев 100ms.
	//
	// что гарантировано всегда:
	//   - пользователи стартуют только после комментариев (<-commentsCh)
	//   - вложения грузятся только при непустом sessionID
	//   - main не завершится, пока все 4 горутины не вызовут wg.Done()
	//
	// что может меняться при каждом запуске:
	//   - порядок строк этапа 1 (зависит от sleep)
	//   - порядок этапов 2 и 3 (они параллельны друг другу)
}
