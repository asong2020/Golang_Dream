package main

import (
	"fmt"
)

func Hello(ch chan int)  {
	fmt.Println("hello everybody , I'm asong")
	ch <- 2
}

func main()  {
	ch := make(chan int)
	go Hello(ch)
	<-ch
	fmt.Println("Golang梦工厂")
}
