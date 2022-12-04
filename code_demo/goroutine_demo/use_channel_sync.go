package main

import (
	"fmt"
	"sync"
	"time"
)

func main()  {
	count := 9 // 要运行的goroutine数量
	limit := 3 // 同时运行的goroutine为3个
	ch := make(chan bool, limit)
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i:=0; i < count; i++{
		go func(num int) {
			defer wg.Done()
			ch <- true // 发送信号
			fmt.Printf("%d 我在干活 at time %d\n",num,time.Now().Unix())
			time.Sleep(2 * time.Second)
			<- ch // 接受数据代表退出了
		}(i)
	}
	wg.Wait()
}
