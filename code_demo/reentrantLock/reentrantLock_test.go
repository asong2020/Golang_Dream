package reentrantLock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type TestReentrantLock struct {
	mu sync.Locker
	id int64
}

func NewTestReentrantLock() *TestReentrantLock{
	return &TestReentrantLock{
		mu: NewReentrantLock(),
	}
}

func (t *TestReentrantLock)SetID()  {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.getID()
	t.id = 1
}

func (t *TestReentrantLock) getID() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.id

}

func TestSingleGoroutine(t *testing.T) {
	rt := NewTestReentrantLock()
	rt.SetID()
	fmt.Println(rt.getID())
}

func TestMuliteGoroutine(t *testing.T)  {
	rt := NewTestReentrantLock()
	signel := func() {
		rt.SetID()
		fmt.Println(rt.getID())
	}
	for i:=0 ; i < 5 ; i++{
		go signel()
	}
	time.Sleep(10 * time.Second)
}