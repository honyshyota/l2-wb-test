package main

import "fmt"

// Паттерн «фасад»

// В качестве демонстрации использовал родительские и дочернюю структуру
// которые по сути имплементируют основную систему и подсистемы 

// подсистема sms
type sms struct{} 

func (s *sms) sendSMS() string {
	return "sms sending"
}

// подсистема mms
type mms struct{}

func (m *mms) sendMMS() string {
	return "mms sending"
}

// подсистема email
type email struct{}

func (e *email) sendEmail() string {
	return "email sending"
}

// основная система
type notifySys struct {
	sms   *sms
	mms   *mms
	email *email
}

// конструктор
func newNotifySys() *notifySys {
	return &notifySys{
		sms:   &sms{},
		mms:   &mms{},
		email: &email{},
	}
}

// и реализация
func main() {
	notifySys := newNotifySys()

	fmt.Println(notifySys.sms.sendSMS())
	fmt.Println(notifySys.mms.sendMMS())
	fmt.Println(notifySys.email.sendEmail())
}
