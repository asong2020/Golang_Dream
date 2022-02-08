package main

import (
	"context"
	"fmt"
	"time"
)

func main()  {
	ctx := context.WithValue(context.Background(), "asong", "test01")
	go func() {
		for {
			_ = context.WithValue(ctx, "asong", "test02")
		}
	}()
	go func() {
		for {
			_ = context.WithValue(ctx, "asong", "test03")
		}
	}()
	go func() {
		for {
			fmt.Println(ctx.Value("asong"))
		}
	}()
	go func() {
		for {
			fmt.Println(ctx.Value("asong"))
		}
	}()
	time.Sleep(10 * time.Second)
}
