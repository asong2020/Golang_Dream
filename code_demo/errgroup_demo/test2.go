package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main()  {
	group,ctx := errgroup.WithContext(context.Background())
	for i:=0;i<5;i++{
		index := i
		group.Go(func() error {
			fmt.Printf("start to execute the %d goroutine\n", index)
			select {
			case <- ctx.Done():
				fmt.Printf("goroutine: %d 被取消了\n",index)
			default:
				time.Sleep(time.Duration(index) * time.Second)
				if index % 2 == 0{
					return fmt.Errorf("somthing has failed on goroutine: %d", index)
				}
				fmt.Printf("goroutine: %d end\n",index)
				return nil
			}
			return nil

		})
	}
	if err := group.Wait(); err != nil{
		fmt.Println(err)
	}
}