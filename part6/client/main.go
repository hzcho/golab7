package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	url := "ws://localhost:8080/ws"
	fmt.Println("Подключение к серверу:", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Ошибка подключения к серверу:", err)
	}
	defer conn.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("Ошибка чтения сообщения:", err)
				return
			}
			fmt.Printf("[%s]: %s\n", msg.Username, msg.Content)
		}
	}()

	fmt.Print("Введите ваше имя: ")
	var username string
	fmt.Scanln(&username)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Отключение клиента...")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Ошибка закрытия соединения:", err)
				return
			}
			time.Sleep(time.Second)
			return
		default:
			fmt.Print("\nВведите сообщение: ")
			var message string
			fmt.Scanln(&message)

			msg := Message{
				Username: username,
				Content:  message,
			}

			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println("Ошибка отправки сообщения:", err)
				return
			}
		}
	}
}
