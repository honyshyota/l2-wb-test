package main

import (
	"fmt"
)

type Service interface {
	Execute(*Data)
	SetNext(Service)
}

type Data struct {
	FirstMessage  string
	SecondMessage string
}

type fStep struct {
	Next Service
}

func (f *fStep) Execute(data *Data) {
	fmt.Println("Данные обработаны первым обработчиком")
	data.FirstMessage = "Первый обработчик"
	f.Next.Execute(data)
}

func (f *fStep) SetNext(s Service) {
	f.Next = s
}

type sStep struct {
	Next Service
}

func (s *sStep) Execute(data *Data) {
	fmt.Println("Данные обработаны вторым обработчиком")
	data.SecondMessage = "Второй обработчик"
	s.Next.Execute(data)
}

func (s *sStep) SetNext(svc Service) {
	s.Next = svc
}

type Save struct {
	data []*Data
	Next Service
}

func (s *Save) Execute(data *Data) {
	s.data = append(s.data, data)
	fmt.Println("Данные сохранены")
}

func (s *Save) SetNext(svc Service) {
	s.Next = svc
}

func (s *Save) Print() {
	for _, data := range s.data {
		fmt.Println(data.FirstMessage, data.SecondMessage)
	}
}

func main() {
	fStep := &fStep{}
	sStep := &sStep{}
	save := &Save{}

	fStep.SetNext(sStep)
	sStep.SetNext(save)

	data := &Data{}

	fStep.Execute(data)
	fStep.Execute(data)

	save.Print()
}
