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

	flag.StringVar(&f, "t", "10s", "set timeout") // парсим флаги таймаута
	flag.Parse()

	// создаем канал для получения сигнала ОС, ввиду того что в моем случае стоит последняя версия Ubuntu
	// в ней не работает сигнал SIGQUIT, поэтому собрал с сигналом SIGINT
	exitChOs := make(chan os.Signal, 1)
	signal.Notify(exitChOs, syscall.SIGINT)

	go func(exitCh chan os.Signal) {
		for {
			switch <-exitCh {
			case syscall.SIGINT: // слушаем канал, и если в него приходит сигнал завершаем работу
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

	host := flag.Arg(0) // парсим аргументы хоста и порта
	port := flag.Arg(1)

	connString := host + ":" + port

	conn, err := net.DialTimeout("tcp", connString, timeout) // подключаемся к сокету
	if err != nil {
		logrus.Fatalln("Не удалось соединиться: ", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			time.Sleep(3 * time.Second)           // придумал эти костыли чтоб пинговать сокет
			_, err := net.Dial("tcp", connString) // и в случае его не ответа завершать работу
			if err != nil {
				logrus.Println("Соединение с сокетом разорвано")
				os.Exit(0)
			}
		}
	}()

	stdinReader := bufio.NewReader(os.Stdin) // создаем ридеры для стандартного ввода и ридер соединения
	connReader := bufio.NewReader(conn)

	for { // в цикле читаем из стандартного ввода
		logrus.Println("Введите сообщение: ")

		message, err := stdinReader.ReadString('\n')
		if err != nil {
			logrus.Println("Ошибка ввода")
			continue
		}

		fmt.Fprint(conn, message) // передаем в наш сокет

		responseMessage, err := connReader.ReadString('\n') // получаем ответ
		if err != nil {
			logrus.Error("Ошибка получения сообщения")
			continue
		}

		// тут тоже у меня проблема ибо при переносе строки добавляется не только символ переноса,
		// но и символ возврата каретки
		responseMessage = strings.TrimSpace(responseMessage)

		logrus.Println("Сообщение от сервера: ", responseMessage) // ну и печатаем сообщение
	}
}
