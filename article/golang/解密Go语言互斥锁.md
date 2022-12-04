## 前言

> 哈喽，大家好，我是`asong`。
>
> 当提到并发编程、多线程编程时，都会在第一时间想到锁，锁是并发编程中的同步原语，他可以保证多线程在访问同一片内存时不会出现竞争来保证并发安全；在`Go`语言中更推崇由`channel`通过通信的方式实现共享内存，这个设计点与许多主流编程语言不一致，但是`Go`语言也在`sync`包中提供了互斥锁、读写锁，毕竟`channel`也不能满足所有场景，互斥锁、读写锁的使用与我们是分不开的，所以接下来我会分两篇来分享互斥锁、读写锁是怎么实现的，本文我们先来看看互斥锁的实现。

本文基于`Golang`版本：**1.18**





## Go语言互斥锁设计实现

###  mutex介绍

`sync` 包下的`mutex`就是互斥锁，其提供了三个公开方法：调用`Lock()`获得锁，调用`Unlock()`释放锁，在`Go1.18`新提供了`TryLock()`方法可以非阻塞式的取锁操作：

- `Lock()`：调用`Lock`方法进行加锁操作，使用时应注意在同一个`goroutine`中必须在锁释放时才能再次上锁，否则会导致程序`panic`。
- `Unlock()`：调用`UnLock`方法进行解锁操作，使用时应注意未加锁的时候释放锁会引起程序`panic`，已经锁定的 Mutex 并不与特定的 goroutine 相关联，这样可以利用一个 goroutine 对其加锁，再利用其他 goroutine 对其解锁。
- `tryLock()`：调用`TryLock`方法尝试获取锁，当锁被其他 goroutine 占有，或者当前锁正处于饥饿模式，它将立即返回 false，当锁可用时尝试获取锁，获取失败不会自旋/阻塞，也会立即返回false；

`mutex`的结构比较简单只有两个字段：

```go
type Mutex struct {
	state int32
	sema  uint32
}
```

- `state`：表示当前互斥锁的状态，复合型字段；
- `sema`：信号量变量，用来控制等待`goroutine`的阻塞休眠和唤醒

初看结构你可能有点懵逼，互斥锁应该是一个复杂东西，怎么就两个字段就可以实现？那是因为设计使用了位的方式来做标志，`state`的不同位分别表示了不同的状态，使用最小的内存来表示更多的意义，其中低三位由低到高分别表示`mutexed`、`mutexWoken` 和 `mutexStarving`，剩下的位则用来表示当前共有多少个`goroutine`在等待锁：

```go
const (
   mutexLocked = 1 << iota // 表示互斥锁的锁定状态
   mutexWoken // 表示从正常模式被从唤醒
   mutexStarving // 当前的互斥锁进入饥饿状态
   mutexWaiterShift = iota // 当前互斥锁上等待者的数量
)
```

![截屏2022-06-26 下午12.08.46](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-06-26%20%E4%B8%8B%E5%8D%8812.08.46.png)

`mutex`最开始的实现只有正常模式，在正常模式下等待的线程按照先进先出的方式获取锁，但是新创建的`gouroutine`会与刚被唤起的 `goroutine`竞争，会导致刚被唤起的 `goroutine`获取不到锁，这种情况的出现会导致线程长时间被阻塞下去，所以`Go`语言在`1.9`中进行了优化，引入了饥饿模式，当`goroutine`超过`1ms`没有获取到锁，就会将当前互斥锁切换到饥饿模式，在饥饿模式中，互斥锁会直接交给等待队列最前面的`goroutine`，新的 goroutine 在该状态下不能获取锁、也不会进入自旋状态，它们只会在队列的末尾等待。如果一个 goroutine 获得了互斥锁并且它在队列的末尾或者它等待的时间少于 1ms，那么当前的互斥锁就会切换回正常模式。

`mutex`的基本情况大家都已经掌握了，接下来我们从加锁到解锁来分析`mutex`是如何实现的；



## Lock加锁

从`Lock`方法入手：

```go
func (m *Mutex) Lock() {
	// 判断当前锁的状态，如果锁是完全空闲的，即m.state为0，则对其加锁，将m.state的值赋为1
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}
```

上面的代码主要两部分逻辑：

- 通过`CAS`判断当前锁的状态，也就是`state`字段的低1位，如果锁是完全空闲的，即m.state为0，则对其加锁，将m.state的值赋为1
- 若当前锁已经被其他`goroutine`加锁，则进行`lockSlow`方法尝试通过自旋或饥饿状态下饥饿`goroutine`竞争方式等待锁的释放，我们在下面介绍`lockSlow`方法；

`lockSlow`代码段有点长，主体是一个`for`循环，其主要逻辑可以分为以下三部分：
- 状态初始化
- 判断是否符合自旋条件，符合条件进行自旋操作
- 抢锁准备期望状态
- 通过`CAS`操作更新期望状态

### 初始化状态

在`locakSlow`方法内会先初始化5个字段：

```go
func (m *Mutex) lockSlow() {
	var waitStartTime int64 
	starving := false
	awoke := false
	iter := 0
	old := m.state
	........
}
```

- `waitStartTime`用来计算`waiter`的等待时间
- `starving`是饥饿模式标志，如果等待时长超过1ms，starving置为true，后续操作会把Mutex也标记为饥饿状态。
- `awoke`表示协程是否唤醒，当`goroutine`在自旋时，相当于CPU上已经有在等锁的协程。为避免Mutex解锁时再唤醒其他协程，自旋时要尝试把Mutex置为唤醒状态，Mutex处于唤醒状态后 要把本协程的 awoke 也置为true。
- `iter`用于记录协程的自旋次数，
- `old`记录当前锁的状态



### 自旋

自旋的判断条件非常苛刻：

```go
for {
    // 判断是否允许进入自旋 两个条件，条件1是当前锁不能处于饥饿状态
    // 条件2是在runtime_canSpin内实现，其逻辑是在多核CPU运行，自旋的次数小于4
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
      // !awoke 判断当前goroutine不是在唤醒状态
      // old&mutexWoken == 0 表示没有其他正在唤醒的goroutine
      // old>>mutexWaiterShift != 0 表示等待队列中有正在等待的goroutine
      // atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) 尝试将当前锁的低2位的Woken状态位设置为1，表示已被唤醒, 这是为了通知在解锁Unlock()中不要再唤醒其他的waiter了
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
					// 设置当前goroutine唤醒成功
          awoke = true
			}
      // 进行自旋
			runtime_doSpin()
      // 自旋次数
			iter++
      // 记录当前锁的状态
			old = m.state
			continue
		}
}
```

自旋这里的条件还是很复杂的，我们想让当前`goroutine`进入自旋转的原因是我们乐观的认为**当前正在持有锁的goroutine能在较短的时间内归还锁**，所以我们需要一些条件来判断，`mutex`的判断条件我们在文字描述一下：

`old&(mutexLocked|mutexStarving) == mutexLocked` 用来判断锁是否处于正常模式且加锁，为什么要这么判断呢？ 

`mutexLocked` 二进制表示为 0001

`mutexStarving` 二进制表示为 0100

`mutexLocked|mutexStarving` 二进制为 0101. 使用0101在当前状态做 `&`操作，如果当前处于饥饿模式，低三位一定会是1，如果当前处于加锁模式，低1位一定会是1，所以使用该方法就可以判断出当前锁是否处于正常模式且加锁；

`runtime_canSpin()`方法用来判断是否符合自旋条件：

```go
// / go/go1.18/src/runtime/proc.go
const active_spin     = 4
func sync_runtime_canSpin(i int) bool {
	if i >= active_spin || ncpu <= 1 || gomaxprocs <= int32(sched.npidle+sched.nmspinning)+1 {
		return false
	}
	if p := getg().m.p.ptr(); !runqempty(p) {
		return false
	}
	return true
}
```

自旋条件如下：

- 自旋的次数要在4次以内
- `CPU`必须为多核
- `GOMAXPROCS>1`
- 当前机器上至少存在一个正在运行的处理器 P 并且处理的运行队列为空；

判断当前`goroutine`可以进自旋后，调用`runtime_doSpin`方法进行自旋：

```go
const active_spin_cnt = 30
func sync_runtime_doSpin() {
	procyield(active_spin_cnt)
}
// asm_amd64.s
TEXT runtime·procyield(SB),NOSPLIT,$0-0
	MOVL	cycles+0(FP), AX
again:
	PAUSE
	SUBL	$1, AX
	JNZ	again
	RET
```

循环次数被设置为`30`次，自旋操作就是执行30次`PAUSE`指令，通过该指令占用`CPU`并消费`CPU`时间，进行忙等待；

这就是整个自旋操作的逻辑，这个就是为了优化 等待阻塞->唤醒->参与抢占锁这个过程不高效，所以使用自旋进行优化，在期望在这个过程中锁被释放。



### 抢锁准备期望状态

自旋逻辑处理好后开始根据上下文计算当前互斥锁最新的状态，根据不同的条件来计算`mutexLocked`、`mutexStarving`、`mutexWoken` 和 `mutexWaiterShift`：

首先计算`mutexLocked`的值：

```go
    // 基于old状态声明到一个新状态
		new := old
		// 新状态处于非饥饿的条件下才可以加锁
		if old&mutexStarving == 0 {
			new |= mutexLocked
		}
```

计算`mutexWaiterShift`的值：

```go
//如果old已经处于加锁或者饥饿状态，则等待者按照FIFO的顺序排队
if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}
```

计算`mutexStarving`的值：

```go
// 如果当前锁处于饥饿模式，并且已被加锁，则将低3位的Starving状态位设置为1，表示饥饿
if starving && old&mutexLocked != 0 {
			new |= mutexStarving
		}
```

计算`mutexWoken`的值：

```go
// 当前goroutine的waiter被唤醒,则重置flag
if awoke {
			// 唤醒状态不一致，直接抛出异常
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
     // 新状态清除唤醒标记，因为后面的goroutine只会阻塞或者抢锁成功
     // 如果是挂起状态，那就需要等待其他释放锁的goroutine来唤醒。
     // 假如其他goroutine在unlock的时候发现Woken的位置不是0，则就不会去唤醒，那该goroutine就无法在被唤醒后加锁
			new &^= mutexWoken
}
```



### 通过`CAS`操作更新期望状态

上面我们已经得到了锁的期望状态，接下来通过`CAS`将锁的状态进行更新：

```go
// 这里尝试将锁的状态更新为期望状态
if atomic.CompareAndSwapInt32(&m.state, old, new) {
  // 如果原来锁的状态是没有加锁的并且不处于饥饿状态，则表示当前goroutine已经获取到锁了，直接推出即可
			if old&(mutexLocked|mutexStarving) == 0 {
				break // locked the mutex with CAS
			}
			// 到这里就表示goroutine还没有获取到锁，waitStartTime是goroutine开始等待的时间，waitStartTime != 0就表示当前goroutine已经等待过了，则需要将其放置在等待队列队头，否则就排到队列队尾
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
      // 阻塞等待
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)
      // 被信号量唤醒后检查当前goroutine是否应该表示为饥饿
     // 1. 当前goroutine已经饥饿
     // 2. goroutine已经等待了1ms以上
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
  // 再次获取当前锁的状态
			old = m.state
   // 如果当前处于饥饿模式，
			if old&mutexStarving != 0 {
        // 如果当前锁既不是被获取也不是被唤醒状态，或者等待队列为空 这代表锁状态产生了不一致的问题
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
        // 当前goroutine已经获取了锁，等待队列-1
				delta := int32(mutexLocked - 1<<mutexWaiterShift
         // 当前goroutine非饥饿状态 或者 等待队列只剩下一个waiter，则退出饥饿模式(清除饥饿标识位)              
				if !starving || old>>mutexWaiterShift == 1 {
					delta -= mutexStarving
				}
        // 更新状态值并中止for循环，拿到锁退出
				atomic.AddInt32(&m.state, delta)
				break
			}
      // 设置当前goroutine为唤醒状态，且重置自璇次数
			awoke = true
			iter = 0
		} else {
      // 锁被其他goroutine占用了，还原状态继续for循环
			old = m.state
		}
```

这块的逻辑很复杂，通过`CAS`来判断是否获取到锁，没有通过 CAS 获得锁，会调用 `runtime.sync_runtime_SemacquireMutex`通过信号量保证资源不会被两个 `goroutine` 获取，`runtime.sync_runtime_SemacquireMutex`会在方法中不断尝试获取锁并陷入休眠等待信号量的释放，一旦当前 `goroutine` 可以获取信号量，它就会立刻返回，如果是新来的`goroutine`，就需要放在队尾；如果是被唤醒的等待锁的`goroutine`，就放在队头，整个过程还需要啃代码来加深理解。



## 解锁

相对于加锁操作，解锁的逻辑就没有那么复杂了，接下来我们来看一看`UnLock`的逻辑：

```go
func (m *Mutex) Unlock() {
	// Fast path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if new != 0 {
		// Outlined slow path to allow inlining the fast path.
		// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
		m.unlockSlow(new)
	}
}
```

使用`AddInt32`方法快速进行解锁，将m.state的低1位置为0，然后判断新的m.state值，如果值为0，则代表当前锁已经完全空闲了，结束解锁，不等于`0`说明当前锁没有被占用，会有等待的`goroutine`还未被唤醒，需要进行一系列唤醒操作，这部分逻辑就在`unlockSlow`方法内：

```go
func (m *Mutex) unlockSlow(new int32) {
  // 这里表示解锁了一个没有上锁的锁，则直接发生panic
	if (new+mutexLocked)&mutexLocked == 0 {
		throw("sync: unlock of unlocked mutex")
	}
  // 正常模式的释放锁逻辑
	if new&mutexStarving == 0 {
		old := new
		for {
      // 如果没有等待者则直接返回即可
      // 如果锁处于加锁的状态，表示已经有goroutine获取到了锁，可以返回
      // 如果锁处于唤醒状态，这表明有等待的goroutine被唤醒了，不用尝试获取其他goroutine了
      // 如果锁处于饥饿模式，锁之后会直接给等待队头goroutine
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				return
			}
			// 抢占唤醒标志位，这里是想要把锁的状态设置为被唤醒，然后waiter队列-1
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
        // 抢占成功唤醒一个goroutine
				runtime_Semrelease(&m.sema, false, 1)
				return
			}
      // 执行抢占不成功时重新更新一下状态信息，下次for循环继续处理
			old = m.state
		}
	} else {
    // 饥饿模式释放锁逻辑，直接唤醒等待队列goroutine
		runtime_Semrelease(&m.sema, true, 1)
	}
}
```

我们在唤醒`goroutine`时正常模式/饥饿模式都调用`func runtime_Semrelease(s *uint32, handoff bool, skipframes int)`，这两种模式在第二个参数的传参上不同，如果`handoff is true, pass count directly to the first waiter.`。



## 非阻塞加锁

`Go`语言在`1.18`版本中引入了非阻塞加锁的方法`TryLock()`，其实现就很简洁：

```go
func (m *Mutex) TryLock() bool {
  // 记录当前状态
	old := m.state
  //  处于加锁状态/饥饿状态直接获取锁失败
	if old&(mutexLocked|mutexStarving) != 0 {
		return false
	}
	// 尝试获取锁，获取失败直接获取失败
	if !atomic.CompareAndSwapInt32(&m.state, old, old|mutexLocked) {
		return false
	}


	return true
}
```

`TryLock`的实现就比较简单了，主要就是两个判断逻辑：

- 判断当前锁的状态，如果锁处于加锁状态或饥饿状态直接获取锁失败
- 尝试获取锁，获取失败直接获取锁失败

`TryLock`并不被鼓励使用，至少我还没想到有什么场景可以使用到它。



## 总结

通读源码后你会发现互斥锁的逻辑真的十分复杂，代码量虽然不多，但是很难以理解，一些细节点还需要大家多看看几遍才能理解其为什么这样做，文末我们再总结一下互斥锁的知识点：

- 互斥锁有两种模式：正常模式、饥饿模式，饥饿模式的出现是为了优化正常模式下刚被唤起的`goroutine`与新创建的`goroutine`竞争时长时间获取不到锁，在`Go1.9`时引入饥饿模式，如果一个`goroutine`获取锁失败超过`1ms`,则会将`Mutex`切换为饥饿模式，如果一个`goroutine`获得了锁，并且他在等待队列队尾 或者 他等待小于`1ms`，则会将`Mutex`的模式切换回正常模式
- 加锁的过程：
  - 锁处于完全空闲状态，通过CAS直接加锁
  - 当锁处于正常模式、加锁状态下，并且符合自旋条件，则会尝试最多4次的自旋
  - 若当前`goroutine`不满足自旋条件时，计算当前goroutine的锁期望状态
  - 尝试使用CAS更新锁状态，若更新锁状态成功判断当前`goroutine`是否可以获取到锁，获取到锁直接退出即可，若获取不到锁则陷入睡眠，等待被唤醒
  - goroutine被唤醒后，如果锁处于饥饿模式，则直接拿到锁，否则重置自旋次数、标志唤醒位，重新走for循环自旋、获取锁逻辑；
- 解锁的过程
  - 原子操作mutexLocked，如果锁为完全空闲状态，直接解锁成功
  - 如果锁不是完全空闲状态，，那么进入`unlockedslow`逻辑
  - 如果解锁一个未上锁的锁直接panic，因为没加锁`mutexLocked`的值为0，解锁时进行mutexLocked - 1操作，这个操作会让整个互斥锁混乱，所以需要有这个判断
  - 如果锁处于饥饿模式直接唤醒等待队列队头的waiter
  - 如果锁处于正常模式下，没有等待的goroutine可以直接退出，如果锁已经处于锁定状态、唤醒状态、饥饿模式则可以直接退出，因为已经有被唤醒的 `goroutine` 获得了锁.
- 使用互斥锁时切记拷贝`Mutex`，因为拷贝`Mutex`时会连带状态一起拷贝，因为`Lock`时只有锁在完全空闲时才会获取锁成功，拷贝时连带状态一起拷贝后，会造成死锁
- TryLock的实现逻辑很简单，主要判断当前锁处于加锁状态、饥饿模式就会直接获取锁失败，尝试获取锁失败直接返回；

本文之后你对互斥锁有什么不理解的吗？欢迎评论区批评指正～；

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**


![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)











