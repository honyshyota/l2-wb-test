package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Db struct {
	user []*User
}

func (db *Db) Save(u *User, dataName string) {
	u.GenID()
	u.SaveData(dataName)
	db.user = append(db.user, u)
}

type User struct {
	Id   int
	Data *Data
}

func (u *User) GenID() {
	rand.Seed(time.Now().Unix())
	u.Id = rand.Intn(100)
}

func (u *User) SaveData(dataName string) {
	data := &Data{}
	data.SetData(dataName)
	u.Data = data.Get()
}

type Data struct {
	Name string
}

func (d *Data) SetData(dataName string) {
	d.Name = dataName
}

func (d *Data) Get() *Data {
	return d
}

func main() {
	u1 := &User{}
	u2 := &User{}

	db := Db{}

	db.Save(u1, "some data 1")
	db.Save(u2, "some data 2")

	fmt.Println(db.user[0].Data)
}
