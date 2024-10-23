package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := ":8080"

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println("Ошибка загрузки сертификата:", err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", port, config)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Сервер запущен на порту", port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Завершение работы слушателя...")
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Ошибка при принятии соединения:", err)
					continue
				}
				go handleConnection(ctx, conn)
			}
		}
	}()

	<-stop
	fmt.Println("Получен сигнал завершения работы, закрываем сервер...")
	cancel()

	time.Sleep(time.Second * 2)
	fmt.Println("Сервер остановлен")
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	fmt.Println("Клиент подключен:", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Завершение соединения с:", conn.RemoteAddr())
			return
		default:
			if scanner.Scan() {
				message := scanner.Text()
				fmt.Println("Получено сообщение:", message)

				_, err := conn.Write([]byte("Сообщение получено\n"))
				if err != nil {
					fmt.Println("Ошибка при отправке ответа:", err)
					return
				}
			} else {
				if err := scanner.Err(); err != nil {
					fmt.Println("Ошибка при чтении данных:", err)
				}
				return
			}
		}
	}
}
