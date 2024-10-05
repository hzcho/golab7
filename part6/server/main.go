package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Обновление соединений через безопасную горутину
var clients = make(map[*websocket.Conn]bool) // все активные клиенты
var broadcast = make(chan Message)           // канал для передачи сообщений
var mutex = sync.Mutex{}                     // защита для изменения карты клиентов

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Структура сообщения
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	// Запуск HTTP-сервера
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Println("Запуск сервера на :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Обновление HTTP-соединения до веб-сокета
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Добавление нового клиента
	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	// Прослушивание сообщений от клиента
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Ошибка чтения сообщения: %v", err)
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			break
		}
		// Отправка полученного сообщения в канал
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Получение нового сообщения
		msg := <-broadcast

		// Отправка сообщения всем клиентам
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
