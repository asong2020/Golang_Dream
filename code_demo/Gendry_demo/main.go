package main

import (
	"fmt"
)

type Animal interface {
	Walk()
}

type Dog struct{}

func (d *Dog) Walk() {
	fmt.Println("walk")
}

func NewAnimal() Animal {
	var d *Dog
	return d
}

func main() {
	if NewAnimal() == nil {
		fmt.Println("this is empty interface")
	} else {
		fmt.Println("this is non-empty interface")
	}
}