package main

import (
	"fmt"
	"sync"
)

// 三个协程交替打印1、2、3
func main() {
	print123()
}

func print123() {
	chanA, chanB, chanC := make(chan bool), make(chan bool), make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go printA(chanA, chanC, &wg)
	go printB(chanA, chanB, &wg)
	go printC(chanB, chanC, &wg)
	wg.Wait()
}

func printA(chanA chan bool, chanC chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		fmt.Println("1")
		chanA <- true
		// 阻塞等待C
		<-chanC
	}
}

func printB(chanA chan bool, chanB chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		<-chanA
		fmt.Println("2")
		chanB <- true
	}
}

func printC(chanB chan bool, chanC chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		<-chanB
		fmt.Println("3")
		chanC <- true
	}
}
