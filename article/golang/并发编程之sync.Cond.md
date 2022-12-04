## 前言

哈喽，大家好，我是`asong`，这是我并发编程系列的第三篇文章，这一篇我们一起来看看`sync.Cond`的使用与实现。之前写过`java`的朋友对等待/通知(wait/notify)机制一定很熟悉，可以利用等待/通知机制实现阻塞或者唤醒，在`Go`语言使用`Cond`也可以达到同样的效果，接下来我们一起来看看它的使用与实现。



## `sync.Cond`的基本使用

`Go`标准库提供了`Cond`原语，为等待/通知场景下的并发问题提供支持。`Cond`他可以让一组的`Goroutine`都在满足特定条件(这个等待条件有很多，可以是某个时间点或者某个变量或一组变量达到了某个阈值，还可以是某个对象的状态满足了特定的条件)时被唤醒，`Cond`是和某个条件相关，这个条件需要一组`goroutine`协作共同完成，在条件还没有满足的时候，所有等待这个条件的`goroutine`都会被阻塞住，只有这一组`goroutine`通过协作达到了这个条件，等待的`goroutine`才可以继续进行下去。

先看这样一个例子：

```go
var (
	done = false
	topic = "Golang梦工厂"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	go Consumer(topic,cond)
	go Consumer(topic,cond)
	go Consumer(topic,cond)
	Push(topic,cond)
	time.Sleep(5 * time.Second)

}

func Consumer(topic string,cond *sync.Cond)  {
	cond.L.Lock()
	for !done{
		cond.Wait()
	}
	fmt.Println("topic is ",topic," starts Consumer")
	cond.L.Unlock()
}

func Push(topic string,cond *sync.Cond)  {
	fmt.Println(topic,"starts Push")
	cond.L.Lock()
	done = true
	cond.L.Unlock()
	fmt.Println("topic is ",topic," wakes all")
	cond.Broadcast()
}
// 运行结果
Golang梦工厂 starts Push
topic is  Golang梦工厂  wakes all
topic is  Golang梦工厂  starts Consumer
topic is  Golang梦工厂  starts Consumer
topic is  Golang梦工厂  starts Consumer
```

上述代码我们运行了`4`个`Goroutine`，其中三个`Goroutine`分别做了相同的事情，通过调用`cond.Wait()`等特定条件的满足，1个`Goroutine`会调用`cond.Broadcast`唤醒所用陷入等待的`Goroutine`。画个图看一下更清晰：
![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/cond-1.png)

我们看上面这一段代码，`Cond`使用起来并不简单，使用不当就出现不可避免的问题，所以，有的开发者会认为，`Cond`是唯一难以掌握的`Go`并发原语。为了让大家能更好的理解`Cond`，接下来我们一起看看`Cond`的实现原理。



## `Cond`实现原理

`Cond`的实现还是比较简单的，代码量比较少，复杂的逻辑已经被`Locker`或者`runtime`的等待队列实现了，所以我们来看这些源代码也会轻松一些。首先我们来看一下它的结构体：

```go
type Cond struct {
	noCopy noCopy

	// L is held while observing or changing the condition
	L Locker

	notify  notifyList
	checker copyChecker
}
```

主要有`4`个字段：

- `nocopy` ：之前在讲`waitGroup`时介绍过，保证结构体不会在编译器期间拷贝，原因就不在这里说了，想了解的看这篇文章[源码剖析sync.WaitGroup(文末思考题你能解释一下吗?)](https://mp.weixin.qq.com/s/r9g4ZQLTYJ5QGvBmVIM8YA)
- `checker`：用于禁止运行期间发生拷贝，双重检查(`Double check`)
- `L`：可以传入一个读写锁或互斥锁，当修改条件或者调用`wait`方法时需要加锁
- `notify`：通知链表，调用`wait()`方法的`Goroutine`会放到这个链表中，唤醒从这里取。我们可以看一下`notifyList`的结构：

```go
type notifyList struct {
	wait   uint32
	notify uint32
	lock   uintptr // key field of the mutex
	head   unsafe.Pointer
	tail   unsafe.Pointer
}
```

我们简单分析一下`notifyList`的各个字段：

- `wait`：下一个等待唤醒`Goroutine`的索引，他是在锁外自动递增的.
- `notify`：下一个要通知的`Goroutine`的索引，他可以在锁外读取，但是只能在锁持有的情况下写入.
- `head`：指向链表的头部
- `tail`：指向链表的尾部

基本结构我们都知道了，下面我就来看一看`Cond`提供的三种方法是如何实现的～。



### `wait`

我们先来看一下`wait`方法源码部分：

```go
func (c *Cond) Wait() {
	c.checker.check()
	t := runtime_notifyListAdd(&c.notify)
	c.L.Unlock()
	runtime_notifyListWait(&c.notify, t)
	c.L.Lock()
}
```

代码量不多，执行步骤如下：

- 执行运行期间拷贝检查，如果发生了拷贝，则直接`panic`程序
- 调用`runtime_notifyListAdd`将等待计数器加一并解锁；
- 调用`runtime_notifyListWait`等待其他 `Goroutine` 的唤醒并加锁

`runtime_notifyListAdd`的实现：

```go
// See runtime/sema.go for documentation.
func notifyListAdd(l *notifyList) uint32 {
	// This may be called concurrently, for example, when called from
	// sync.Cond.Wait while holding a RWMutex in read mode.
	return atomic.Xadd(&l.wait, 1) - 1
}
```

代码实现比较简单，原子操作将等待计数器加1，因为`wait`代表的是下一个等待唤醒`Goroutine`的索引，所以需要减1操作。

`runtime_notifyListWait`的实现：

```go
// See runtime/sema.go for documentation.
func notifyListWait(l *notifyList, t uint32) {
	lockWithRank(&l.lock, lockRankNotifyList)

	// Return right away if this ticket has already been notified.
	if less(t, l.notify) {
		unlock(&l.lock)
		return
	}

	// Enqueue itself.
	s := acquireSudog()
	s.g = getg()
	s.ticket = t
	s.releasetime = 0
	t0 := int64(0)
	if blockprofilerate > 0 {
		t0 = cputicks()
		s.releasetime = -1
	}
	if l.tail == nil {
		l.head = s
	} else {
		l.tail.next = s
	}
	l.tail = s
	goparkunlock(&l.lock, waitReasonSyncCondWait, traceEvGoBlockCond, 3)
	if t0 != 0 {
		blockevent(s.releasetime-t0, 2)
	}
	releaseSudog(s)
}
```

这里主要执行步骤如下：

- 检查当前`wait`与`notify`索引位置是否匹配，如果已经被通知了，便立即返回.
- 获取当前`Goroutine`，并将当前`Goroutine`追加到链表末端.
- 调用`goparkunlock`方法让当前`Goroutine`进入等待状态，也就是进入睡眠，等待唤醒
- 被唤醒后，调用`releaseSudog`释放当前等待列表中的`Goroutine`

看完源码我们来总结一下注意事项：

`wait`方法会把调用者放入`Cond`的等待队列中并阻塞，直到被唤醒，调用`wait`方法必须要持有`c.L`锁。





### `signal`和`Broadcast`

`signal`和`Broadcast`都会唤醒等待队列，不过`signal`是唤醒链表最前面的`Goroutine`，`Boradcast`会唤醒队列中全部的`Goroutine`。下面我们分别来看一下`signal`和`broadcast`的源码：

- `signal`

```go
func (c *Cond) Signal() {
	c.checker.check()
	runtime_notifyListNotifyOne(&c.notify)
}
func notifyListNotifyOne(l *notifyList) {
	if atomic.Load(&l.wait) == atomic.Load(&l.notify) {
		return
	}
	lockWithRank(&l.lock, lockRankNotifyList)
	t := l.notify
	if t == atomic.Load(&l.wait) {
		unlock(&l.lock)
		return
	}

	atomic.Store(&l.notify, t+1)

	for p, s := (*sudog)(nil), l.head; s != nil; p, s = s, s.next {
		if s.ticket == t {
			n := s.next
			if p != nil {
				p.next = n
			} else {
				l.head = n
			}
			if n == nil {
				l.tail = p
			}
			unlock(&l.lock)
			s.next = nil
			readyWithTime(s, 4)
			return
		}
	}
	unlock(&l.lock)
}
```

上面我们看`wait`源代码时，每次都会调用都会原子递增`wait`，那么这个`wait`就代表当前最大的`wait`值，对应唤醒的时候，也就会对应一个`notify`属性，我们在`notifyList`链表中逐个检查，找到`ticket`对应相等的`notify`属性。这里大家肯定会有疑惑，我们为何不直接取链表头部唤醒呢？

`notifyList`并不是一直有序的，`wait`方法中调用`runtime_notifyListAdd`和`runtime_notifyListWait`完全是两个独立的行为，中间还有释放锁的行为，而当多个 `goroutine` 同时进行时，中间会产生进行并发操作，这样就会出现乱序，所以采用这种操作即使在 `notifyList` 乱序的情况下，也能取到最先` Wait` 的 `goroutine`。



- `broadcast`

```go
func (c *Cond) Broadcast() {
	c.checker.check()
	runtime_notifyListNotifyAll(&c.notify)
}
func notifyListNotifyAll(l *notifyList) {
	if atomic.Load(&l.wait) == atomic.Load(&l.notify) {
		return
	}

	lockWithRank(&l.lock, lockRankNotifyList)
	s := l.head
	l.head = nil
	l.tail = nil

	atomic.Store(&l.notify, atomic.Load(&l.wait))
	unlock(&l.lock)
	for s != nil {
		next := s.next
		s.next = nil
		readyWithTime(s, 4)
		s = next
	}
}
```

全部唤醒实现要简单一些，主要是通过调用`readyWithTime`方法唤醒链表中的`goroutine`，唤醒顺序也是按照加入队列的先后顺序，先加入的会先被唤醒，而后加入的可能 `Goroutine` 需要等待调度器的调度。

最后我们总结一下使用这两个方法要注意的问题：

- `Signal`：允许调用者唤醒一个等待此`Cond`的`Goroutine`，如果此时没有等待的 `goroutine`，显然无需通知` waiter`；如果` Cond` 等待队列中有一个或者多个等待的 `goroutine`，则需要从等待队列中移除第一个 `goroutine` 并把它唤醒。调用 `Signal `方法时，不强求你一定要持有 `c.L` 的锁。
- `broadcast`：允许调用者唤醒所有等待此 `Cond` 的` goroutine`。如果此时没有等待的` goroutine`，显然无需通知 waiter；如果 `Cond` 等待队列中有一个或者多个等待的 `goroutine`，则清空所有等待的 `goroutine`，并全部唤醒，不强求你一定要持有 `c.L` 的锁。

## 注意事项

- 调用`wait`方法的时候一定要加锁，否则会导致程序发生`panic`.
- `wait`调用时需要检查等待条件是否满足，也就说`goroutine`被唤醒了不等于等待条件被满足，等待者被唤醒，只是得到了一次检查的机会而已，推荐写法如下：

```go
//    c.L.Lock()
//    for !condition() {
//        c.Wait()
//    }
//    ... make use of condition ...
//    c.L.Unlock()
```

-  `Signal` 和 `Boardcast` 两个唤醒操作不需要加锁



## 总结

其实`Cond`在实际项目中被使用的机会比较少，`Go`特有的`channel`就可以代替它，暂时只在`Kubernetes`项目中看到了应用，使用场景是每次往队列中成功增加了元素后就需要调用 `Broadcast` 通知所有的等待者，使用`Cond`就很合适，相比`channel`减少了代码复杂性。

**好啦，这篇文章就到这里啦，素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：加我vx拉你入群，或者公众号获取入群二维码**

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**翻译了一份Machinery中文文档，会定期进行维护，有需要的小伙伴们后台回复[machinery]即可获取。**

**我是asong，一名普普通通的程序猿，让我们一起慢慢变强吧。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [Go看源码必会知识之unsafe包](https://mp.weixin.qq.com/s/nPWvqaQiQ6Z0TaPoqg3t2Q)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/mzSCWI8C_ByIPbb07XYFTQ)
- [详解并发编程基础之原子操作(atomic包)](https://mp.weixin.qq.com/s/PQ06eL8kMWoGXodpnyjNcA)
- [详解defer实现机制](https://mp.weixin.qq.com/s/FUmoBB8OHNSfy7STR0GsWw)
- [真的理解interface了嘛](https://mp.weixin.qq.com/s/sO6Phr9C5VwcSTQQjJux3g)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/wURQFRt2ISz66icW7jbHFw)
- [十张动图带你搞懂排序算法(附go实现代码)](https://mp.weixin.qq.com/s/rZBsoKuS-ORvV3kML39jKw)
- [go参数传递类型](https://mp.weixin.qq.com/s/JHbFh2GhoKewlemq7iI59Q)
- [手把手教姐姐写消息队列](https://mp.weixin.qq.com/s/0MykGst1e2pgnXXUjojvhQ)
- [常见面试题之缓存雪崩、缓存穿透、缓存击穿](https://mp.weixin.qq.com/s?__biz=MzIzMDU0MTA3Nw==&mid=2247483988&idx=1&sn=3bd52650907867d65f1c4d5c3cff8f13&chksm=e8b0902edfc71938f7d7a29246d7278ac48e6c104ba27c684e12e840892252b0823de94b94c1&token=1558933779&lang=zh_CN#rd)
- [详解Context包，看这一篇就够了！！！](https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg)
- [go-ElasticSearch入门看这一篇就够了(一)](https://mp.weixin.qq.com/s/mV2hnfctQuRLRKpPPT9XRw)
- [面试官：go中for-range使用过吗？这几个问题你能解释一下原因吗](https://mp.weixin.qq.com/s/G7z80u83LTgLyfHgzgrd9g)

