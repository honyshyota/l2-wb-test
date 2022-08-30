package main

import (
	"bufio"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Println("Server is running...")

	l, err := net.Listen("tcp", ":8181") // слушаем порт
	if err != nil {
		logrus.Fatalln("Server error", err)
		return
	}
	defer l.Close()

	for {
		logrus.Println("Connecting wait...")

		conn, err := l.Accept() // ждем соединения
		if err != nil {
			logrus.Fatalln("Error accepting", err)
			return
		}
		defer conn.Close()

		logrus.Println("New connection: ", conn.RemoteAddr())
		logrus.Println("Waiting message...")

		reader := bufio.NewReader(conn) // создаем ридер для соединения

		for { // в цикле принимаем сообщения
			message, err := reader.ReadString('\n')
			if err != nil {
				logrus.Println("Error reading", err)
			}

			message = strings.TrimSpace(message) // так же тримим их

			logrus.Println("Message received: ", conn.RemoteAddr(), string(message))

			responseMessage := message + " from server" // немного редактируем входящее сообщение чтоб была видна разница

			_, err = conn.Write([]byte(responseMessage + "\n")) // ну и отдаем обратно отправителю
			if err != nil {
				logrus.Fatalln("Error writing: ", err)
				break
			}
		}

		logrus.Println("Client disconnect: ", conn.RemoteAddr())
	}
}
