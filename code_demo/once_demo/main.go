package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type MyOnce struct {
	flag uint32
	lock sync.Mutex
}

func (m *MyOnce)Do(f func())  {
	if atomic.LoadUint32(&m.flag) == 0{
		m.lock.Lock()
		defer m.lock.Unlock()
		if atomic.CompareAndSwapUint32(&m.flag,0,1){
			f()
		}
	}
}

func testDo()  {
	mOnce := MyOnce{}
	for i := 0;i<10;i++{
		go func() {
			mOnce.Do(func() {
				fmt.Println("test my once only run once")
			})
		}()
	}
}

func main()  {
	nestedDo()
	//testDo()
	////Concurrent()
	//time.Sleep(10 * time.Second)
	//Deadlock()
}


// 并发测试 sync.once是否在程序运行期间只执行一次
func Concurrent()  {
	once := &sync.Once{}
	for i:= 0;i<10;i++{
		go func() {
			once.Do(func() {
				fmt.Printf("run in goroutine")
			})
		}()
	}
}

func nestedDo()  {
	once1 := &sync.Once{}
	once2 := &sync.Once{}
	once1.Do(func() {
		once2.Do(func() {
			fmt.Println("test nestedDo")
		})
	})
}

func panicDo()  {
	once := &sync.Once{}
	defer func() {
		if err := recover();err != nil{
			once.Do(func() {
				fmt.Println("run in recover")
			})
		}
	}()
	once.Do(func() {
		panic("panic i=0")
	})

}

