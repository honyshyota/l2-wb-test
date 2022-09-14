package main

import "fmt"

type Visitor interface {
	VisitFactory(p *WorkShop)
}

type Place interface {
	Accept(v Visitor)
}

type Worker struct {
	content string
}

func (w *Worker) VisitFactory(p *WorkShop) {
	w.content = p.BuildingCar()
}

func (w *Worker) Print() {
	fmt.Printf("Workshop - %s\n", w.content)
}

type Places struct {
	place []Place
}

func (p *Places) SetPlace(place Place) {
	p.place = append(p.place, place)
}

func (p *Places) Accepts(v Visitor) {
	for _, place := range p.place {
		place.Accept(v)
	}
}

type WorkShop struct{}

func (w *WorkShop) Accept(v Visitor) {
	v.VisitFactory(w)
}

func (w *WorkShop) BuildingCar() string {
	return "Build car"
}

func main() {
	places := &Places{}
	places.SetPlace(&WorkShop{})

	visitor := &Worker{}

	places.Accepts(visitor)

	visitor.Print()
}
