package main

import (
	"context"
	"fmt"
	"time"
)

func main()  {
	now := time.Now()
	later,_:=time.ParseDuration("10s")
	
	ctx,cancel := context.WithDeadline(context.Background(),now.Add(later))
	defer cancel()
	go Monitor(ctx)

	time.Sleep(20 * time.Second)

}

func Monitor(ctx context.Context)  {
	select {
	case <- ctx.Done():
		fmt.Println(ctx.Err())
	case <-time.After(20*time.Second):
		fmt.Println("stop monitor")
	}
}
