package main

import (
	"context"
	"fmt"
	"time"
)

func main()  {
	HttpHandler()
}

func NewContextWithTimeout() (context.Context,context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3 * time.Second)
}

func HttpHandler()  {
	ctx, cancel := NewContextWithTimeout()
	defer cancel()
	deal(ctx)
}

func deal(ctx context.Context)  {
	for i:=0; i< 10; i++ {
		time.Sleep(1*time.Second)
		select {
		case <- ctx.Done():
			fmt.Println(ctx.Err())
			return
		default:
			fmt.Printf("deal time is %d\n", i)
		}
	}
}