package main

import (
	"context"
	"fmt"
	"time"
)

//func main()  {
//	HttpHandler1()
//}

func NewContextWithTimeout1() (context.Context,context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3 * time.Second)
}

func HttpHandler1()  {
	ctx, cancel := NewContextWithTimeout1()
	defer cancel()
	deal1(ctx, cancel)
}

func deal1(ctx context.Context, cancel context.CancelFunc)  {
	for i:=0; i< 10; i++ {
		time.Sleep(1*time.Second)
		select {
		case <- ctx.Done():
			fmt.Println(ctx.Err())
			return
		default:
			fmt.Printf("deal time is %d\n", i)
			cancel()
		}
	}
}