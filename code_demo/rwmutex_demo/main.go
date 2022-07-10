package main

import "sync"

const maxValue = 3

type test struct {
	rw sync.RWMutex
	index int
}

func (t *test) Get() int {
	return t.index
}

func (t *test)Set() {
	t.rw.Lock()
	t.index++
	if t.index >= maxValue{
		t.index =0
	}
	t.rw.Unlock()
}

func main()  {
	print(1 << 30)
	t := test{}
	sw := sync.WaitGroup{}
	for i:=0; i < 100000; i++{
		sw.Add(2)
		go func() {
			t.Set()
			sw.Done()
		}()
		go func() {
			val := t.Get()
			if val >= maxValue{
				print("get value error| value=", val, "\n")
			}
			sw.Done()
		}()
	}
	sw.Wait()
}
