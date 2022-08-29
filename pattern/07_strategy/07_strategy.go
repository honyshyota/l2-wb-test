package main

import (
	"fmt"
	"log"
)

// Паттерн "Стратегия"

// Strategy Интерфейс определяющий поведение типов различных стратегий
type Strategy interface {
	work() string
}

// Определяем несколько пользовательских типов и имплементируем интерфейс
type constructorRobot struct{}

func (r *constructorRobot) work() string {
	return "build robot"
}

type removerRobot struct{}

func (r *removerRobot) work() string {
	return "robot remove"
}

// Определяем контекст с полем интерфейс
type context struct {
	strategy Strategy
}

// Присваиваем полю переданный аргумент
func (c *context) algorithm(a Strategy) {
	c.strategy = a
}

// Вызываем через контекст нужный нам метод переданного типа в контекст
func (c *context) work() string {
	return c.strategy.work()
}

// Определяем в этой функции требуемый алгоритм
func algorithmFind(s string) *context {
	if s == "build" {
		ctx := &context{}
		ctx.algorithm(&constructorRobot{})
		return ctx
	} else if s == "remove" {
		ctx := &context{}
		ctx.algorithm(&removerRobot{})
		return ctx
	} else {
		return nil
	}
}

func main() {
	ctx := algorithmFind("build") // Инициализируем контекст и передаем в него запрос на нужную стратегию
	if ctx == nil {
		log.Fatal("can't find strategy for this task")
		return
	}
	fmt.Println(ctx.work()) // Вывод результата в stdout

	ctx = algorithmFind("remove")
	if ctx == nil {
		log.Fatal("can't find strategy for this task")
		return
	}
	fmt.Println(ctx.work())

	ctx = algorithmFind("incorrect")
	if ctx == nil {
		log.Fatal("can't find strategy for this task")
		return
	}
	fmt.Println(ctx.work())
}
