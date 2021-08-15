package main

import (
	"fmt"
)

type User struct {
	Name string
	Age int
}

func GetName(u *User)  {
	fmt.Println(u.Name)
}

func SetName(u *User)  {
	u.Name = "asong真帅"
}

func main()  {
	u := User{
		Name: "asong",
		Age: 18,
	}
	GetName(&u)
	SetName(&u)
}