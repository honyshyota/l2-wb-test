package main

import (
	"fmt"
)

type req string

const (
	Sedan req = "sedan"
	Wagon req = "wagon"
)

type Factory interface {
	NewFactory(req) Product
}

type Product interface {
	Print()
}

type Creator struct{}

func NewCreator() Factory {
	return &Creator{}
}

func (c *Creator) NewFactory(typeName req) Product {
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

func NewSedan() Product {
	return &CarSedan{}
}

func (c *CarSedan) Print() {
	fmt.Println("sedan build")
}

type CarWagon struct{}

func NewWagon() Product {
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
