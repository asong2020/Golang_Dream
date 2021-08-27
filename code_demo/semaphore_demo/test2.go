package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/semaphore"
)

func main()  {
	s := semaphore.NewWeighted(3)
	ctx,cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()

	for i :=0; i < 3; i++{
			if i != 0{
				go func(num int) {
					if err := s.Acquire(ctx,3); err != nil{
						fmt.Printf("goroutine： %d, err is %s\n", num, err.Error())
						return
					}
					time.Sleep(2 * time.Second)
					fmt.Printf("goroutine： %d run over\n",num)
					s.Release(3)

				}(i)
			}else {
				go func(num int) {
					ct,cancel := context.WithTimeout(context.Background(), time.Second * 3)
					defer cancel()
					if err := s.Acquire(ct,3); err != nil{
						fmt.Printf("goroutine： %d, err is %s\n", num, err.Error())
						return
					}
					time.Sleep(3 * time.Second)
					fmt.Printf("goroutine： %d run over\n",num)
					s.Release(3)
				}(i)
			}

	}
	time.Sleep(10 * time.Second)
}