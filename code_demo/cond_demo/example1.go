package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	done = false
	topic = "Golang梦工厂"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})

	go Consumer1(topic,cond)
	time.Sleep(time.Second)
	go Consumer2(topic,cond)
	go Push(topic,cond)
	time.Sleep(5 * time.Second)

}

func Consumer1(topic string,cond *sync.Cond)  {
	cond.L.Lock()
	for !done{
		cond.Wait()
	}
	fmt.Println("topic is ",topic," starts Consumer1")
	cond.L.Unlock()
}
func Consumer2(topic string,cond *sync.Cond)  {
	cond.L.Lock()
	for !done{
		cond.Wait()
	}
	fmt.Println("topic is ",topic," starts Consumer2")
	cond.L.Unlock()
}

func Push(topic string,cond *sync.Cond)  {
	fmt.Println(topic,"starts Push")
	cond.L.Lock()
	done = true
	cond.L.Unlock()
	fmt.Println("topic is ",topic," wakes all")
	cond.Signal()
}



