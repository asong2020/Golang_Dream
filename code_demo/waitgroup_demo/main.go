package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main()  {
	atomic.Value{}
	testNoCopy()
	ctx,cancel := context.WithTimeout(context.Background(),1*time.Second)
	defer cancel()
	fmt.Println(ctx)
}

func testNoCopy()  {
	wgPointer := &sync.WaitGroup{}
	newWgPointer := wgPointer
	fmt.Println(wgPointer,newWgPointer)
	wgValue := sync.WaitGroup{}
	newWgValue := wgValue
	fmt.Println(wgValue,newWgValue)
}