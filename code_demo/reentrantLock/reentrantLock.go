package reentrantLock

import (
	"fmt"
	"sync"
)

type ReentrantLock struct {
	lock *sync.Mutex
	cond *sync.Cond
	recursion int32
	host     int64
}

func NewReentrantLock()  sync.Locker{
	res := &ReentrantLock{
		lock: new(sync.Mutex),
		recursion: 0,
		host: 0,
	}
	res.cond = sync.NewCond(res.lock)
	return res
}

func (rt *ReentrantLock) Lock()  {
	id := GetGoroutineID()
	rt.lock.Lock()
	defer rt.lock.Unlock()

	if rt.host == id{
		rt.recursion++
		return
	}

	for rt.recursion != 0{
		rt.cond.Wait()
	}
	rt.host = id
	rt.recursion = 1
}

func (rt *ReentrantLock) Unlock()  {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	if rt.recursion == 0 || rt.host != GetGoroutineID() {
		panic(fmt.Sprintf("the wrong call host: (%d); current_id: %d; recursion: %d", rt.host,GetGoroutineID(),rt.recursion))
	}

	rt.recursion--
	if rt.recursion == 0{
		rt.cond.Signal()
	}
}

