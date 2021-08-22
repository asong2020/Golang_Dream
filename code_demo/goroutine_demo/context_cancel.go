package main

import (
	"context"
	"fmt"
	"time"
)


func main()  {
	ctx,cancel := context.WithCancel(context.Background())
	go Speak(ctx)
	time.Sleep(10*time.Second)
	cancel()
	time.Sleep(2 * time.Second)
	fmt.Println("bye bye!")
}

func Speak(ctx context.Context)  {
	for range time.Tick(time.Second){
		select {
		case <- ctx.Done():
			fmt.Println("asong哥，我收到信号了，要走了，拜拜！")
			return
		default:
			fmt.Println("asong哥，你好帅呀～balabalabalabala")
		}
	}
}
