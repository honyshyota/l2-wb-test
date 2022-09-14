package main

import "fmt"

type Command interface {
	Execute()
}

type pString struct {
	receiver *Receiver
}

func (p *pString) Execute() {
	p.receiver.printString()
}

type pInt struct {
	receiver *Receiver
}

func (p *pInt) Execute() {
	p.receiver.printInt()
}

type Invoker struct {
	commands []Command
}

func (i *Invoker) SetCommands(command Command) {
	i.commands = append(i.commands, command)
}

func (i *Invoker) DeleteCommands() {
	if len(i.commands) != 0 {
		i.commands = i.commands[:len(i.commands)-1]
	}
}

func (i *Invoker) StartCommand() {
	for _, command := range i.commands {
		command.Execute()
	}
}

type Receiver struct{}

func (r *Receiver) printString() {
	fmt.Println("some string")
}

func (r *Receiver) printInt() {
	fmt.Println(12345)
}

func main() {
	i := &Invoker{}
	r := &Receiver{}
	comString := &pString{receiver: r}
	comInt := &pInt{receiver: r}

	i.SetCommands(comString)
	i.SetCommands(comInt)
	i.StartCommand()
}
