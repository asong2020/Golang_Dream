package main

import (
	"fmt"
	"time"
)

func main()  {
	ch := make(chan int, 10)
	go func() {
		for i:=0; i<10;i++{
			ch <- i
		}
		close(ch)
	}()
	go func() {
		for val := range ch{
			fmt.Println(val)
		}
		fmt.Println("receive data over")
	}()
	time.Sleep(5* time.Second)
	fmt.Println("program over")
}

