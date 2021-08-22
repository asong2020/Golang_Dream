package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	Limit = 3  // 同时运行的goroutine上限
	Weight = 1 // 信号量的权重
)

func main() {
	names := []string{
		"asong1",
		"asong2",
		"asong3",
		"asong4",
		"asong5",
		"asong6",
		"asong7",
	}

	sem := semaphore.NewWeighted(Limit)
	var w sync.WaitGroup
	for _, name := range names {
		w.Add(1)
		go func(name string) {
			sem.Acquire(context.Background(), Weight)
			fmt.Println(name)
			time.Sleep(2 * time.Second)
			sem.Release(Weight)
			w.Done()
		}(name)
	}
	w.Wait()

	fmt.Println("over--------")
}