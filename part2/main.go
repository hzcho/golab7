package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
)

func main() {
	serverAddress := "localhost:8080"

	// Настройка TLS
	config := &tls.Config{
		InsecureSkipVerify: true, // Используйте только для тестирования
	}

	conn, err := tls.Dial("tcp", serverAddress, config)
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Подключено к серверу", serverAddress)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите сообщение: ")
	message, _ := reader.ReadString('\n')

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка при отправке сообщения:", err)
		return
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при получении ответа:", err)
		return
	}
	fmt.Println("Ответ от сервера:", response)
}
