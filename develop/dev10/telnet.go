package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

/*
Реализовать простейший telnet-клиент.

Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Требования:
Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s)
При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера,
программа должна также завершаться. При подключении к несуществующему сервер, программа должна завершаться через timeout
*/

func main() {
	var f string

	flag.StringVar(&f, "t", "10s", "set timeout")
	flag.Parse()

	exitChOs := make(chan os.Signal, 1)
	signal.Notify(exitChOs, syscall.SIGINT)

	go func(exitCh chan os.Signal) {
		for {
			switch <-exitCh {
			case syscall.SIGINT:
				logrus.Println("System shutdown...")
				os.Exit(0)
			default:
			}
		}
	}(exitChOs)

	timeout, err := time.ParseDuration(f)
	if err != nil {
		logrus.Error("Ошибка парсинга таймаута: ", err)
		return
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	connString := host + ":" + port

	conn, err := net.DialTimeout("tcp", connString, timeout)
	if err != nil {
		logrus.Fatalln("Не удалось соединиться: ", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			time.Sleep(3 * time.Second)
			_, err := net.Dial("tcp", connString)
			if err != nil {
				logrus.Println("Соединение с сокетом разорвано")
				os.Exit(0)
			}
		}
	}()

	stdinReader := bufio.NewReader(os.Stdin)
	connReader := bufio.NewReader(conn)

	for {
		logrus.Println("Введите сообщение: ")

		message, err := stdinReader.ReadString('\n')
		if err != nil {
			logrus.Println("Ошибка ввода")
			continue
		}

		fmt.Fprint(conn, message)

		responseMessage, err := connReader.ReadString('\n')
		if err != nil {
			logrus.Error("Ошибка получения сообщения")
			continue
		}

		responseMessage = strings.TrimSpace(responseMessage)

		logrus.Println("Сообщение от сервера: ", responseMessage)
	}
}
