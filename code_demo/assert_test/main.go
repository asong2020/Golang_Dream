package main

import (
	"fmt"
)

type Basic interface {
	GetName() string
	SetName(name string) error
}

type User struct {
	Name string
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) SetName(name string) error {
	u.Name = name
	return nil
}

type MoreMethod interface {
	Set() string
	Get() string
	One() string
	Two() string
	Three() string
	Four() string
	Five() string
	Six() string
	Seven() string
	Eight() string
	Nine() string
	Ten() string
}

type More struct {
	Name string
}

func (m *More) Set() string   { return m.Name }
func (m *More) Get() string   { return m.Name }
func (m *More) One() string   { return m.Name }
func (m *More) Two() string   { return m.Name }
func (m *More) Three() string { return m.Name }
func (m *More) Four() string  { return m.Name }
func (m *More) Five() string  { return m.Name }
func (m *More) Six() string   { return m.Name }
func (m *More) Seven() string { return m.Name }
func (m *More) Eight() string { return m.Name }
func (m *More) Nine() string  { return m.Name }
func (m *More) Ten() string   { return m.Name }

func main() {
	var u Basic = &User{Name: "asong"}
	v, ok := u.(Basic)
	if !ok {
		fmt.Printf("%v\n", v)
	}
}

//switch u.(type) {
//case *User:
//	u1 := u.(*User)
//	fmt.Println(u1.GetName())
//default:
//	fmt.Println("failed to match")
//}

