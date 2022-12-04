package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// ifaceWords is interface{} internal representation.
type ifaceWords struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}

type Value struct {
	v interface{}
}

func main()  {
	//checkABA()
	v := &Value{}
	vp := (*ifaceWords)(unsafe.Pointer(v))
	fmt.Println(atomic.LoadPointer(&vp.typ))
	fmt.Printf("&vp.typ: %v\n",&vp.typ)
	fmt.Printf("%v\n",unsafe.Pointer(^uintptr(0)))
	fmt.Println(atomic.CompareAndSwapPointer(&vp.typ, nil, unsafe.Pointer(^uintptr(0))))
	fmt.Println(uintptr(0))
}

func checkABA()  {
	var share uint64 = 1
	wg := sync.WaitGroup{}
	wg.Add(3)
	// 协程1，期望值是1,欲更新的值是2
	go func() {
		defer wg.Done()
		swapped := atomic.CompareAndSwapUint64(&share,1,2)
		fmt.Println("goroutine 1",swapped)
	}()
	// 协程2，期望值是1，欲更新的值是2
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond)
		swapped := atomic.CompareAndSwapUint64(&share,1,2)
		fmt.Println("goroutine 2",swapped)
	}()
	// 协程3，期望值是2，欲更新的值是1
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Millisecond)
		swapped := atomic.CompareAndSwapUint64(&share,2,1)
		fmt.Println("goroutine 3",swapped)
	}()
	wg.Wait()
	fmt.Println("main exit")
}
