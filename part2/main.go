package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddress := "localhost:8080"

	conn, err := net.Dial("tcp", serverAddress)
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
