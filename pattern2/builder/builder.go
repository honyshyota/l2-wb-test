package main

import (
	"fmt"
)

type Builder interface {
	SetBody() Builder
	SetEngine() Builder
	SetTransmission() Builder
	GetCar() *Car
}

type Collector struct {
	builder Builder
}

func (c *Collector) Build() *Car {
	return c.builder.SetBody().SetEngine().SetTransmission().GetCar()
}

type CarBuilder struct {
	Body         string
	Engine       string
	Transmission string
}

func (c *CarBuilder) SetBody() Builder {
	c.Body = "sedan"
	return c
}

func (c *CarBuilder) SetEngine() Builder {
	c.Engine = "inline 4"
	return c
}

func (c *CarBuilder) SetTransmission() Builder {
	c.Transmission = "MT"
	return c
}

func (c *CarBuilder) GetCar() *Car {
	return &Car{
		Body:         c.Body,
		Engine:       c.Engine,
		Transmission: c.Transmission,
	}
}

type Car struct {
	Body         string
	Engine       string
	Transmission string
}

func (c *Car) Print() {
	fmt.Printf("Car body - [%s], engine - [%s], transmission - [%s]\n", c.Body, c.Engine, c.Transmission)
}

func main() {
	carBuilder := &Collector{&CarBuilder{}}

	car := carBuilder.Build()

	car.Print()
}
