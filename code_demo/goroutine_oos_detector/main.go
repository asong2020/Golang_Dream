package main

import (
	"fmt"
	"runtime"
	"time"
)

func GetData() {
	ch := make(chan struct{})
	go func() {
		<- ch
	}()
	ch <- struct{}{}
}

func main()  {
	defer func() {
		fmt.Println("goroutines: ", runtime.NumGoroutine())
	}()
	GetData()
	time.Sleep(2 * time.Second)
}
