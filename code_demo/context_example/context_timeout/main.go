package main

import (
	"context"
	"fmt"
	"time"
)

//func main()  {
//
//	ctx,cancel := context.WithTimeout(context.Background(),10 * time.Second)
//	defer cancel()
//	go Monitor(ctx)
//
//	time.Sleep(20 * time.Second)
//
//}

func Monitor(ctx context.Context)  {
	select {
	case <- ctx.Done():
		fmt.Println(ctx.Err())
	case <-time.After(20*time.Second):
		fmt.Println("stop monitor")
	}
}
