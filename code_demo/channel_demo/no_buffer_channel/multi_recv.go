package main

import (
	"fmt"
	"time"
)

func main(){
	ch := make(chan int)
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("one",<-ch)
	}()
	go func() {
		fmt.Println("two",<-ch)
	}()
	time.Sleep(1 * time.Second)
	ch <- 1
	time.Sleep(1 * time.Second)
	fmt.Println("over")
}
