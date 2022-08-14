package main

import "fmt"

// Паттерн "Фабрика"

// Интерфейс задачи
type Task interface {
	taskDescription(int)
}

// Интерфейс описывающий поведение фабричного типа
type IRobot interface {
	getName() string
	getPower() int
}

// Фабричный тип
type robot struct {
	name  string
	power int
}

func (r *robot) getName() string {
	return r.name
}

func (r *robot) getPower() int {
	return r.power
}

// Реализуем несколько фабричных типов
type fabricatorRobot struct {
	*robot
}

func newFabricatorRobot() IRobot {
	return &fabricatorRobot{
		&robot{
			name:  "fabricator robot",
			power: 50,
		},
	}
}

type cleanerRobot struct {
	*robot
}

func newCleanerRobot() IRobot {
	return &cleanerRobot{
		&robot{
			name:  "cleaner robot",
			power: 10,
		},
	}
}

// Реализуем описания типа задач которые будет использовать для опредения создания фабрикой нужного нам типа
type task struct{}

func newTask() Task {
	return &task{}
}

func (t *task) taskDescription(power int) {
	if power <= 10 {
		robot := newCleanerRobot()
		fmt.Println("task is performed", robot.getName(), "it has power", robot.getPower(), "kw.")
	} else if power > 10 && power <= 50 {
		robot := newFabricatorRobot()
		fmt.Println("task is performed", robot.getName(), "it has power", robot.getPower(), "kw.")
	} else {
		fmt.Println("there is no robot for this task")
	}
}

func main() {
	task := newTask() // Инициализируем нашу задачу

	task.taskDescription(55) // Передаем в нее аргументы при которых создаются нужные пользовательские типы
	task.taskDescription(5)
	task.taskDescription(11)
}
