package main

import (
	"fmt"
	"time"
)

func GoroutineOne(ch chan <-string)  {
	fmt.Println("GoroutineOne running")
	ch <- "asong真帅"
	fmt.Println("GoroutineOne end of the run")
}

func GoroutineTwo(ch <- chan string)  {
	fmt.Println("GoroutineTwo running")
	fmt.Printf("Two女朋友说：%s\n",<-ch)
	fmt.Println("GoroutineTwo end of the run")
}

func GoroutineThree(ch <- chan string)  {
	fmt.Println("GoroutineThree running")
	fmt.Printf("Three女朋友说: %s\n",<-ch)
	fmt.Println("GoroutineThree end of the run")
}

func main()  {
	ch := make(chan string)
	go GoroutineOne(ch)
	go GoroutineTwo(ch)
	go GoroutineThree(ch)
	time.Sleep(3 * time.Second)
}

