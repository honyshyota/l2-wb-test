package main

import (
	"fmt"
)

// Паттерн "Строитель"

// Интерфейс строителя
type RobotBuilder interface {
	setWeight(int) RobotBuilder
	setTask(string) RobotBuilder
	setPower(string) RobotBuilder
	build() Robot
}

// Интерфейс продукта
type Robot interface {
	walk()
	stop()
	work()
	isDestroy()
	get() *robot
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

func (r *robot) isDestroy() {
	fmt.Println("Robot is destroy? ", r.destroy)
}

func (r *robot) get() *robot {
	return r
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

// И решил создать еще одного строителя
type RobotDestroyer interface {
	setDestroy(bool) RobotDestroyer
	build(Robot) Robot
}

type robotDestroyer struct {
	destroy bool
}

func newDestroyer() RobotDestroyer {
	return &robotDestroyer{}
}

func (rd *robotDestroyer) setDestroy(status bool) RobotDestroyer {
	rd.destroy = status
	return rd
}

// Метод build возращает интерфейс и записывает отличающиеся данные в переменную
// в зависимости от состояния поля destroy
func (rd *robotDestroyer) build(r Robot) Robot {
	robo := r.get()

	if rd.destroy {
		return &robot{
			weight:  0,
			task:    "",
			power:   "",
			destroy: rd.destroy,
		}
	} else {
		return &robot{
			weight:  robo.weight,
			task:    robo.task,
			power:   robo.power,
			destroy: rd.destroy,
		}
	}
}

func main() {
	robotBuilder := newBuilder() // Создаем строителя

	robot := robotBuilder.setPower("50kw").setWeight(140).setTask("shoot").build() // Передаем в него нужные параметры

	robot.walk() // Проверяем что получили
	robot.work()
	robot.stop()

	robotDestroyer := newDestroyer()

	destRobot := robotDestroyer.setDestroy(false).build(robot)

	destRobot.isDestroy()

	destRobot = robotDestroyer.setDestroy(true).build(robot)

	destRobot.isDestroy()
}
