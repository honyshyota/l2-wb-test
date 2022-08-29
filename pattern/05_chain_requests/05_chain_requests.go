package main

import "fmt"

// Паттерн "Цепочка вызовов"

// Handler Интерфейс обработчиков в цепочке
type Handler interface {
	sendRequest(int) string
}

// Первый обработчик
type handler1 struct {
	next Handler
}

// Метод для имплементации интерфейса и обработки аргумента
func (h *handler1) sendRequest(message int) string {
	if message == 1 {
		return "handler 1"
	} else if h.next != nil {
		return h.next.sendRequest(message)
	}

	return "handler missing"
}

// Второй обработчик
type handler2 struct {
	next Handler
}

// Метод для имплементации интерфейса и обработки аргумента
func (h *handler2) sendRequest(message int) string {
	if message == 2 {
		return "handler 2"
	} else if h.next != nil {
		return h.next.sendRequest(message)
	}

	return "handler missing"
}

// Третий обработчик
type handler3 struct {
	next Handler
}

// Метод для имплементации интерфейса и обработки аргумента
func (h *handler3) sendRequest(message int) string {
	if message == 3 {
		return "handler 3"
	} else if h.next != nil {
		return h.next.sendRequest(message)
	}

	return "handler missing"
}

// Четвертый обработчик
type handler4 struct {
	next Handler
}

// Метод для имплементации интерфейса и обработки аргумента
func (h *handler4) sendRequest(message int) string {
	if message == 4 {
		return "handler 4"
	} else if h.next != nil {
		return h.next.sendRequest(message)
	}

	return "handler missing"
}

func main() {
	handlers := &handler1{&handler2{&handler3{&handler4{}}}} // Создаем цепочку, порядок обработчиков можно присвоить произовльный
	fmt.Println(handlers.sendRequest(3))                     // В stout выводится номер обработчика который обработал число
	fmt.Println(handlers.sendRequest(0))
}
