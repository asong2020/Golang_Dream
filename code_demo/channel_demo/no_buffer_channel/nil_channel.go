package main

import (
	"fmt"
)

func main()  {
	var ch chan string
	ch <- "asong真帅"
	fmt.Println(<- ch)
}
