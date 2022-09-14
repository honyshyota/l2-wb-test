package main

import (
	"fmt"
)

type Strategy interface {
	Calculate(int, int)
}

type Calculator struct {
	Strategy
}

func (c *Calculator) SetStrategy(s Strategy) {
	c.Strategy = s
}

type Sum struct{}

func (s *Sum) Calculate(a int, b int) {
	fmt.Println("Sum: ", a+b)
}

type Substract struct{}

func (s *Substract) Calculate(a int, b int) {
	fmt.Println("Substract: ", a-b)
}

var (
	strategies = []Strategy{
		&Sum{},
		&Substract{},
	}
)

func main() {
	calc := &Calculator{}

	for _, strat := range strategies {
		calc.SetStrategy(strat)
		calc.Calculate(5, 4)
	}
}
