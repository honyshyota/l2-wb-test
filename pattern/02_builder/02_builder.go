package main

import (
	"fmt"
)

// Паттерн "Строитель"

// RobotBuilder Интерфейс строителя
type RobotBuilder interface {
	setWeight(int) RobotBuilder
	setTask(string) RobotBuilder
	setPower(string) RobotBuilder
	build() Robot
}

// Robot Интерфейс продукта
type Robot interface {
	walk()
	stop()
	work()
}

// Структура продукта
type robot struct {
	weight  int
	task    string
	power   string
	destroy bool
}

func (r *robot) walk() {
	fmt.Println("Robot walk with power: ", r.power)
}

func (r *robot) stop() {
	fmt.Println("Robot stop with weight: ", r.weight, " kg")
}

func (r *robot) work() {
	fmt.Println("Robot complete task: ", r.task)
}

// Структура строителя
type robotBuilder struct {
	weightOption int
	taskOption   string
	powerOption  string
}

// Конструктор возращающий интерфейс строителя
func newBuilder() RobotBuilder {
	return &robotBuilder{}
}

// Реализация методовод строителя

func (rb *robotBuilder) setWeight(weight int) RobotBuilder {
	rb.weightOption = weight
	return rb
}

func (rb *robotBuilder) setTask(task string) RobotBuilder {
	rb.taskOption = task
	return rb
}

func (rb *robotBuilder) setPower(power string) RobotBuilder {
	rb.powerOption = power
	return rb
}

func (rb *robotBuilder) build() Robot {
	return &robot{
		weight: rb.weightOption,
		task:   rb.taskOption,
		power:  rb.powerOption,
	}
}

func main() {
	robotBuilder := newBuilder() // Создаем строителя

	robot := robotBuilder.setPower("50kw").setWeight(140).setTask("shoot").build() // Передаем в него нужные параметры

	robot.walk() // Проверяем что получили
	robot.work()
	robot.stop()
}
