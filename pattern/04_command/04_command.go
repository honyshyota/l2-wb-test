package main

import "fmt"

// Паттерн "Команда"

// Command Базовый интерфейс описывающий поведение команд
type Command interface {
	execute() bool
}

// Команда включения
type onCommand struct {
	prodaction *robot
}

// Имплементация интерфейса Command
func (c *onCommand) execute() bool {
	return c.prodaction.on()
}

// Команда выключения
type offCommand struct {
	prodaction *robot
}

// Имплементация интерфейса Command
func (c *offCommand) execute() bool {
	return c.prodaction.off()
}

// Инициатор, который будет давать сигнал на выполнение команд
type invoker struct {
	command Command
}

// Сеттер устанавливающий команду переданную аргументом
func (i *invoker) com(command Command) {
	i.command = command
}

// Метод выполнения команды
func (i *invoker) goComm() bool {
	return i.command.execute()
}

// Тип принимающий команды
type robot struct {
	isOn bool
}

// Метод включения
func (r *robot) on() bool {
	r.isOn = true
	return r.isOn
}

// Метод выключения
func (r *robot) off() bool {
	r.isOn = false
	return r.isOn
}

func main() {
	robot := &robot{}     // Инициализируем рессивер
	invoker := &invoker{} // Инициализируем инициатора

	invoker.com(&onCommand{prodaction: robot}) // Устанавливаем команду инициатору
	invoker.goComm()                           // Запрос на выполнение команды
	fmt.Println(robot)                         // Вывод результата в stdout

	invoker.com(&offCommand{prodaction: robot}) // Устанавливаем команду инициатору
	invoker.goComm()                            // Запрос на выполнение команды
	fmt.Println(robot)                          // Вывод результата в stdout
}
