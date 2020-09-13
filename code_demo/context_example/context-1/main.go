package main

import (
	"context"
	"fmt"
	"time"
)

func main()  {
	ctx,cancel := context.WithTimeout(context.Background(),1 * time.Second)
	defer cancel()
	go HelloHandle(ctx,2000*time.Millisecond)
	select {
	case <- ctx.Done():
		fmt.Println("Hello Handle ",ctx.Err())
	}

}

func HelloHandle(ctx context.Context,duration time.Duration)  {

	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}

}