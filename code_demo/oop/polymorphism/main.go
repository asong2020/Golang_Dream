package main

import "fmt"

type SendEmail interface {
	send()
}

func Send(s SendEmail)  {
	s.send()
}

type user struct {
	name string
	email string
}

func (u *user) send()  {
	fmt.Println(u.name + " email is " + u.email + "already send")
}

type admin struct {
	name string
	email string
}

func (a *admin) send()  {
	fmt.Println(a.name + " email is " + a.email + "already send")
}

func main()  {
	u := &user{
		name: "asong",
		email: "你猜",
	}
	a := &admin{
		name: "asong1",
		email: "就不告诉你",
	}
	Send(u)
	Send(a)
}
