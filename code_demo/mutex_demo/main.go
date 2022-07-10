package main

import "sync"

func main()  {
	l := sync.Mutex{}
	l.Lock()
	l.Lock()
	defer l.Unlock()
	defer l.Unlock()
}
