package main

import (
	"fmt"
)

type req string

const (
	Sedan req = "sedan"
	Wagon req = "wagon"
)

type Selector interface {
	NewFactory(req) Factory
}

type Factory interface {
	Print()
}

type Creator struct{}

func NewCreator() Selector {
	return &Creator{}
}

func (c *Creator) NewFactory(typeName req) Factory {
	switch typeName {
	default:
		fmt.Println("Не существует такого кейса")
		return nil
	case Sedan:
		return NewSedan()
	case Wagon:
		return NewWagon()
	}
}

type CarSedan struct{}

func NewSedan() Factory {
	return &CarSedan{}
}

func (c *CarSedan) Print() {
	fmt.Println("sedan build")
}

type CarWagon struct{}

func NewWagon() Factory {
	return &CarWagon{}
}

func (c *CarWagon) Print() {
	fmt.Println("Wagon build")
}

func main() {
	creator := NewCreator()
	sedan := creator.NewFactory(Sedan)
	wagon := creator.NewFactory(Wagon)

	sedan.Print()
	wagon.Print()
}
