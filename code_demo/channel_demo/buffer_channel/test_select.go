package main

import (
	"fmt"
)

func fibonacci(ch chan int, done chan struct{}) {
	x, y := 0, 1
	for {
		select {
		case ch <- x:
			x, y = y, x+y
		case <-done:
			fmt.Println("over")
			return
		}
	}
}
func main() {
	ch := make(chan int)
	done := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-ch)
		}
		done <- struct{}{}
	}()
	fibonacci(ch, done)
}