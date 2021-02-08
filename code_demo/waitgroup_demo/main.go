package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type User struct {
	//noCopy noCopy
	Name string
	Info *Info
}

type Info struct {
	Age int
	Number int
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}
	wg.Wait()
}


func f(i int, wg sync.WaitGroup) {
	fmt.Println(i)
	wg.Done()
}


//func main()  {
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//	doDeadLock(wg)
//	wg.Wait()
//	//u := User{
//	//	Name: "asong",
//	//	Info: &Info{
//	//		Age: 10,
//	//		Number: 24,
//	//	},
//	//}
//	//u1 := u
//	//u1.Name = "Golang梦工厂"
//	//u1.Info.Age = 30
//	//fmt.Println(u.Info.Age,u.Name)
//	//fmt.Println(u1.Info.Age,u1.Name)
//}

func exampleToNoCopy()  {

}

type httpPkg struct{}

func (httpPkg) Get(url string) {}

var http httpPkg

func exampleDemo()  {
	var wg sync.WaitGroup
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	for _, url := range urls {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Fetch the URL.
			http.Get(url)
		}(url)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
}

func exampleImplWaitGroup()  {
	done := make(chan struct{}) // 收10份保护费
	count := 10 // 10个马仔
	for i:=0;i < count;i++{
		go func(i int) {
			defer func() {
				done <- struct {}{}
			}()
			fmt.Printf("马仔%d号收保护费\n",i)
		}(i)
	}
	for i:=0;i< count;i++{
		<- done
		fmt.Printf("马仔%d号已经收完保护费\n",i)
	}
	fmt.Println("所有马仔已经干完活了，开始酒吧消费～")
}

func doDeadLock(wg sync.WaitGroup)  {
	defer wg.Done()
	fmt.Println("do something")
}


type MyWaitGroup struct {
	state1 [3]uint32
}

// state returns pointers to the state and sema fields stored within wg.state1.
func (wg *MyWaitGroup) state() (statep *uint64, semap *uint32) {
	if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
		return (*uint64)(unsafe.Pointer(&wg.state1)), &wg.state1[2]
	} else {
		return (*uint64)(unsafe.Pointer(&wg.state1[1])), &wg.state1[0]
	}
}

func (wg *MyWaitGroup) Add(delta int) {
	statep, _ := wg.state()
	state := atomic.AddUint64(statep, uint64(delta)<<32)
	v := int32(state >> 32)
	w := uint32(state)
	fmt.Println(v,w)
}

// Done decrements the WaitGroup counter by one.
func (wg *MyWaitGroup) Done() {
	wg.Add(-1)
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *MyWaitGroup) Wait() {
	statep, _ := wg.state()
	for {
		state := atomic.LoadUint64(statep)
		v := int32(state >> 32)
		w := uint32(state)
		fmt.Println(v,w)
		if v == 0 {
			return
		}
		// Increment waiters count.
		if atomic.CompareAndSwapUint64(statep, state, state+1) {
			v = int32(state >> 32)
			w = uint32(state)
			fmt.Println(state,v,w)
			if *statep != 0 {
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
			return
		}
	}
}
