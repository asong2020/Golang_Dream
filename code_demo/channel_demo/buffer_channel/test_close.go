package main

import (
	"fmt"
)

func main()  {
	ch := make(chan int, 10)
	ch <- 10
	ch <- 20
	close(ch)
	fmt.Println(<-ch) //1
	fmt.Println(<-ch) //2
	fmt.Println(<-ch) //0
	fmt.Println(<-ch) //0
}