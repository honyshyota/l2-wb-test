package main

import "fmt"

// Паттер "Посетитель"

// Visitor Интерфейс посетителя
type Visitor interface {
	VisitRepairModule(*repairModule) string
	VisitWorkModule(*workModule) string
	VisitSoftwareModule(*softwareModule) string
}

// Module Интерфейс элементов для обхода
type Module interface {
	Accept(v Visitor) string
}

// Структура посетителя
type robot struct {
}

// Методы для обхода сторонних эллементов
func (v *robot) VisitRepairModule(r *repairModule) string {
	return r.Repair()
}

func (v *robot) VisitWorkModule(r *workModule) string {
	return r.Work()
}

func (v *robot) VisitSoftwareModule(r *softwareModule) string {
	return r.UpdateSoftware()
}

// Коллекция элементов для обхода
type prodaction struct {
	modules []Module
}

// Метод для добавления элеменета в список обхода
func (p *prodaction) Add(m Module) {
	p.modules = append(p.modules, m)
}

// Метод для обхода элементов из списка
func (p *prodaction) Accept(v Visitor) string {
	var result string
	for _, p := range p.modules {
		result += p.Accept(v)
	}
	return result
}

// Реализация элементов для обхода
type repairModule struct {
}

func (r *repairModule) Accept(v Visitor) string {
	return v.VisitRepairModule(r)
}

func (r *repairModule) Repair() string {
	return "Robot repaired..."
}

type workModule struct {
}

func (w *workModule) Accept(v Visitor) string {
	return v.VisitWorkModule(w)
}

func (w *workModule) Work() string {
	return "Robot worked..."
}

type softwareModule struct {
}

func (s *softwareModule) Accept(v Visitor) string {
	return v.VisitSoftwareModule(s)
}

func (s *softwareModule) UpdateSoftware() string {
	return "Software updated..."
}

func main() {
	prodaction := &prodaction{} // Создаем коллекцию

	prodaction.Add(&repairModule{}) // Добавляем в коллекцию элементы
	prodaction.Add(&workModule{})
	prodaction.Add(&softwareModule{})

	fmt.Println(prodaction.Accept(&robot{})) // Обходим их нашим посетителем и для примера результат выводим в stdout

}
