package main

import (
	"fmt"
)

//go:generate stringer -type=OrderStatus
type OrderStatus int

const (
	CREATE OrderStatus = iota + 1
	PAID
	DELIVERING
	COMPLETED
	CANCELLED
)

func main()  {
	var a = CREATE
	fmt.Println(a)
}