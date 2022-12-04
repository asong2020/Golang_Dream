package main

import (
	"fmt"
)

func main()  {
	ch := make(chan int,10)
	go func() {
		for i:=0;i<10;i++{
			ch <- i
		}
		close(ch)
	}()
	for val := range ch{
		fmt.Println(val)
	}
	fmt.Println("over")
}
