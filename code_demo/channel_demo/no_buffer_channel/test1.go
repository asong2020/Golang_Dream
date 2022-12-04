package main

import (
	"fmt"
)

// 可出面试题
// 测试chanel出队
func main()  {
	ch := make(chan string,10)
	go func() {
		ch <- "asong"
		close(ch)
	}()
	val, ok := <- ch
	if !ok{
		fmt.Println("over")
		return
	}
	fmt.Printf("print %s\n",val)
}