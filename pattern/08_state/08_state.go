package main

import (
	"fmt"
)

func main() {
	robot := newRobotStatus() // Инициализируем переменную в которой будем менять состояние

	robot.setState(&statusWait{}) // Передаем состояние ожидания

	fmt.Println(robot.status()) // Выводим в stdout результат

	robot.setState(&statusWork{}) // Передаем состояние работы

	fmt.Println(robot.status()) // Выводим в stdout результат
}

// Интерфейс для различных состояний
type RobotState interface {
	status() string
}

// Пользовательский зависимый от состояния
type robot struct {
	state RobotState
}

// Метод возращающий нынешнее состояние переменной
func (r *robot) status() string {
	return r.state.status()
}

// Метод устанавливающий состояние
func (r *robot) setState(state RobotState) {
	r.state = state
}

// Конструктор робота с дефолтным статусом
func newRobotStatus() *robot {
	return &robot{state: nil}
}

// Тип реализующий интерфейс
type statusWait struct{}

// Метод возвращающий состояние
func (s *statusWait) status() string {
	return "robot is waiting..."
}

// Тип реализующий интерфейс
type statusWork struct{}

// Метод возвращающий состояние
func (s *statusWork) status() string {
	return "robot is working..."
}
