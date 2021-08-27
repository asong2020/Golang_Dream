package main


import (
	"container/list"
	"context"
	"sync"
)

type waiter struct {
	n     int64 // 等待调用者权重值
	ready chan<- struct{} // close channel就是唤醒
}


// NewWeighted为并发访问创建一个新的加权信号量，该信号量具有给定的最大权值。
func NewWeighted(n int64) *Weighted {
	w := &Weighted{size: n}
	return w
}

// Weighted provides a way to bound concurrent access to a resource.
// The callers can request access with a given weight.
type Weighted struct {
	size    int64 // 设置一个最大权值
	cur     int64 // 标识当前已被使用的权值
	mu      sync.Mutex // 提供临界区保护
	waiters list.List // 阻塞等待的调用者列表
}


//获取权值为n的信号量，阻塞直到资源可用或ctx Done。成功时返回nil。如果失败，返回ctx.Err()并保持信号量不变。
//如果ctx已经DONE，Acquire仍然可以成功而不被阻塞
func (s *Weighted) Acquire(ctx context.Context, n int64) error {
	s.mu.Lock() // 加锁保护临界区
	// 有资源可用并且没有等待获取权值的goroutine
	if s.size-s.cur >= n && s.waiters.Len() == 0 {
		s.cur += n // 加权
		s.mu.Unlock() // 释放锁
		return nil
	}
	// 要获取的权值n大于最大的权值了
	if n > s.size {
		// 先释放锁，确保其他goroutine调用Acquire的地方不被阻塞
		s.mu.Unlock()
		// 阻塞等待context的返回
		<-ctx.Done()
		return ctx.Err()
	}
	// 走到这里就说明现在没有资源可用了
	// 创建一个channel用来做通知唤醒
	ready := make(chan struct{})
	// 创建waiter对象
	w := waiter{n: n, ready: ready}
	// waiter按顺序入队
	elem := s.waiters.PushBack(w)
	// 释放锁，等待唤醒，别阻塞其他goroutine
	s.mu.Unlock()

	// 阻塞等待唤醒
	select {
	// context关闭
	case <-ctx.Done():
		err := ctx.Err() // 先获取context的错误信息
		s.mu.Lock()
		select {
		case <-ready:
			// 在context被关闭后被唤醒了，那么试图修复队列，假装我们没有取消
			err = nil
		default:
			// 判断是否是第一个元素
			isFront := s.waiters.Front() == elem
			// 移除第一个元素
			s.waiters.Remove(elem)
			// 如果是第一个元素且有资源可用通知其他waiter
			if isFront && s.size > s.cur {
				s.notifyWaiters()
			}
		}
		s.mu.Unlock()
		return err
	// 被唤醒了
	case <-ready:
		return nil
	}
}

// TryAcquire 不阻塞地获取权重为 n 的信号量。
// 成功时返回true，失败时返回false并保持信号量不变
func (s *Weighted) TryAcquire(n int64) bool {
	s.mu.Lock() // 加锁
	// 有资源可用并且没有等待获取资源的goroutine
	success := s.size-s.cur >= n && s.waiters.Len() == 0
	if success {
		s.cur += n
	}
	s.mu.Unlock()
	return success
}

// Release 释放权重为 n 的信号量。
func (s *Weighted) Release(n int64) {
	s.mu.Lock()
	// 释放资源
	s.cur -= n
	// 释放资源大于持有的资源，则会发生panic
	if s.cur < 0 {
		s.mu.Unlock()
		panic("semaphore: released more than held")
	}
	// 通知其他等待的调用者
	s.notifyWaiters()
	s.mu.Unlock()
}

// 通知其他调用者
func (s *Weighted) notifyWaiters() {
	for {
		// 获取等待调用者队列中的队员
		next := s.waiters.Front()
		// 没有要通知的调用者了
		if next == nil {
			break // No more waiters blocked.
		}

		// 断言出waiter信息
		w := next.Value.(waiter)
		if s.size-s.cur < w.n {
			// 没有足够资源为下一个调用者使用时，继续阻塞该调用者，遵循先进先出的原则，
			// 避免需要资源数比较大的waiter被饿死
			//
			// 考虑一个场景，使用信号量作为读写锁，现有N个令牌，N个reader和一个writer
			// 每个reader都可以通过Acquire（1）获取读锁，writer写入可以通过Acquire（N）获得写锁定
			// 但不包括所有的reader，如果我们允许reader在队列中前进，writer将会饿死-总是有一个令牌可供每个reader
			break
		}

		// 获取资源
		s.cur += w.n
		// 从waiter列表中移除
		s.waiters.Remove(next)
		// 使用channel的close机制唤醒waiter
		close(w.ready)
	}
}

