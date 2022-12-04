package main

import (
	"fmt"
	"sync"
	"time"
)

func main()  {
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000 ;i++{
		time.Sleep(10 * time.Millisecond)
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}
	wg.Wait()
}
