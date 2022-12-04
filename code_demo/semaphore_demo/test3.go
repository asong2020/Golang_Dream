package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const (
	limit = 2
) 

func main()  {
	serviceName := []string{
		"cart",
		"order",
		"account",
		"item",
		"menu",
	}
	eg,ctx := errgroup.WithContext(context.Background())
	s := semaphore.NewWeighted(limit)
	for index := range serviceName{
		name := serviceName[index]
		if err := s.Acquire(ctx,1); err != nil{
			fmt.Printf("Acquire failed and err is %s\n", err.Error())
			break
		}
		eg.Go(func() error {
			defer s.Release(1)
			return callService(name)
		})
	}

	if err := eg.Wait(); err != nil{
		fmt.Printf("err is %s\n", err.Error())
		return
	}
	fmt.Printf("run success\n")
}

func callService(name string) error {
	fmt.Println("call ",name)
	time.Sleep(1 * time.Second)
	return nil
}