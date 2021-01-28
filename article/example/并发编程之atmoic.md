## 前言

> 嗨，大家好，我是asong。最近想写一个并发编程系列的文章，使用`Go`也有一段时间了，但是对并发的理解不是很透彻，借着这次总结，希望能更进一步。我们以"原子操作"开篇，对于并发操作而言，原子操作是个非常现实的问题，比较典型的应用的就是`i++`操作，并发情况下，同时对内存中的`i`进行读取，就会产生与预期不符的结果，所以`Go`语言中的`sync/atomic`就是解决这个问题的，接下来我们一起来看一看`Go`的原子操作。



## 什么是原子性、原子操作

原子(atomic)本意是"不能被进一步分割的最小粒子"，而原子操作(atomic operation)意为"不可中断的一个或一系列操作"。其实用大白话说出来就是让多个线程对同一块内存的操作是串行的，不会因为并发操作把内存写的不符合预期。我们来看这样一个例子：
假设现在是一个银行账户系统，用户A想要自己从自己的账户中转1万元到用户B的账户上，直到转帐成功完成一个事务，主要做这两件事：

- 从A的账户中减去1万元，如果A的账户原来就有2万元，现在就变成了1万元
- 给B的账户添加1万元，如果B的账户原来有2万元，那么现在就变成了3万元

假设在操作一的时候，系统发生了故障，导致给B账户添加款项失败了，那么就要进行回滚。回滚就是回到事务之前的状态，我们把这种要么一起成功的操作叫做原子操作，而原子性就是要么完整的被执行、要么完全不执行。


## 如何保证原子性

- 锁机制

在处理器层面，可以采用总线加锁或者对缓存加锁的方式来实现多处理器之间的原子操作。通过加锁保证从系统内存中读取或写入一个字节是原子的，也就是当一个处理器读取一个字节时，其他处理器不能访问这个字节的内存地址。

总线锁：处理器提供一个`Lock#`信号，当一个处理器上在总线上输出此信号时，其他处理器的请求将被阻塞住，那么该处理器可以独占共享内存。总线锁会把`CPU`和内存之间的通信锁住了，在锁定期间，其他处理就不能操作其他内存地址的数据，所以总线锁定的开销比较大，所以处理会在某些场合使用缓存锁进行优化。
缓存锁：内存区域如果被缓存在处理器上的缓存行中，并且在`Lock#`操作期间，那么当它执行操作回写到内存时，处理不在总线上声言`Lock#`信号，而是修改内部的内存地址，并允许它的缓存一致机制来保证操作的原子性，因为缓存一致性机制会阻止同时修改由两个以上处理器缓存的内存区域的数据，其他处理器回写已被锁定的缓存行的数据时，就会使缓存无效。

锁机制虽然可以保证原子性，但是锁机制会存在以下问题：

- 多线程竞争的情况下，频繁的加锁、释放锁会导致较多的上下文切换和调度延时，性能会很差
- 当一个线程占用时间比较长时，就多导致其他需要此锁的线程挂起.

上面我们说的都是悲观锁，要解决这种低效的问题，我们可以采用乐观锁，每次不加锁，而是假设没有冲突去完成某项操作，如果因为冲突失败就重试，直到成功为止。也就是我们接下来要说的CAS(compare and swap).

- CAS(compare and swap)

CAS的全称为`Compare And Swap`，直译就是比较交换。是一条CPU的原子指令，其作用是让`CPU`先进行比较两个值是否相等，然后原子地更新某个位置的值，其实现方式是给予硬件平台的汇编指令，在`intel`的`CPU`中，使用的`cmpxchg`指令，就是说`CAS`是靠硬件实现的，从而在硬件层面提升效率。简述过程是这样：
> 假设包含3个参数内存位置(V)、预期原值(A)和新值(B)。`V`表示要更新变量的值，`E`表示预期值，`N`表示新值。仅当`V`值等于`E`值时，才会将`V`的值设为`N`，如果`V`值和`E`值不同，则说明已经有其他线程在做更新，则当前线程什么都不做，最后`CAS`返回当前`V`的真实值。CAS操作时抱着乐观的态度进行的，它总是认为自己可以成功完成操作。基于这样的原理，CAS操作即使没有锁，也可以发现其他线程对于当前线程的干扰。


伪代码可以这样写：
```go
func CompareAndSwap(int *addr,int oldValue,int newValue) bool{
    if *addr == nil{
        return false
    }
    if *addr == oldValue {
        *addr = newValue
        return true
    }
    return false
}
```
不过上面的代码可能会发生一个问题，也就是`ABA`问题，因为CAS需要在操作值的时候检查下值有没有发生变化，如果没有发生变化则更新，但是如果一个值原来是A，变成了B，又变成了A，那么使用CAS进行检查时会发现它的值没有发生变化，但是实际上却变化了。ABA问题的解决思路就是使用版本号。在变量前面追加上版本号，每次变量更新的时候把版本号加一，那么A－B－A 就会变成1A-2B－3A。


## go语言中如何进行原子操作

在`Go`语言标准库中，`sync/atomic`包将底层硬件提供的原子操作封装成了`Go`的函数，主要分为5个系列的函数，分别是：
- `func SwapXXXX(addr *int32, new int32) (old int32)`系列：其实就是原子性的将`new`值保存到`*addr`并返回旧值。代码表示：
```go
old = *addr
*addr = new
return old
```
- `func CompareAndSwapXXXX((addr *int64, old, new int64) (swapped bool)`系列：其就是原子性的比较`*addr`和old的值，如果相同则将`new`赋值给`*addr`并返回真，代码表示：

```go
if *addr == old{
    *addr = new
    return ture
}
return false
```
- `func AddXXXX(addr *int64, delta int64) (new int64)`系列：原子性的将`val`的值添加到`*addr`并返回新值。代码表示：
```go
*addr += delta
return *addr
```

- `func LoadXXXX(addr *uint32) (val uint32)`系列：原子性的获取`*addr`的值
- `func StoreXXXX(addr *int32, val int32)`原子性的将val值保存到`*addr`

`Go`语言在`1.4`版本时添加一个新的类型`Value`，此类型的值就相当于一个容器，可以被用来"原子地"存储(store)和加载(Load)任意类型的值。这些使用起来都还比较简单，就不写例子了，接下来我们一起看一看这些方法是如何实现的。


## 源码解析
由于系列比较多。底层实现的方法也大同小异样，这里就主要分析一下`Value`的实现方法吧。为什么不分析其他系列的呢？因为原子操作由底层硬件支持，所以看其他系列实现都要看汇编，Go的汇编是基于`Plan9`的，这个汇编语言真的资料甚少，我也是真的不懂，水平不够，也不自讨苦吃了，等后面真的能看懂这些汇编了，再来分析吧。这个网站有一些关于`plan9`汇编的知识，有兴趣可以看一看：http://doc.cat-v.org/plan_9/4th_edition/papers/asm。

### `Value`结构
我们先来看一下`Value`的结构：
```go
type Value struct {
	v interface{}
}
```
`Value`结构里就只有一个字段，是interface类型，虽然这里是`interface`类型，但是这里要注意，第一次`Store`写入的类型就确定了之后写入的类型，否则会发生`panic`。因为这里是`interface`类型，所以为了之后写入与读取操作方便，又在这个包里定义了一个`ifaceWords`结构，其实他就是一个空`interface`，他的作用就是将`interface`
分解成类型和数值。结构如下：
```go
// ifaceWords is interface{} internal representation.
type ifaceWords struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}
```

### `Value`的写入操作
我们一起来看一看他是如何实现写入操作的：
```go
// Store sets the value of the Value to x.
// All calls to Store for a given Value must use values of the same concrete type.
// Store of an inconsistent type panics, as does Store(nil).
func (v *Value) Store(x interface{}) {
	if x == nil {
		panic("sync/atomic: store of nil value into Value")
	}
	vp := (*ifaceWords)(unsafe.Pointer(v))
	xp := (*ifaceWords)(unsafe.Pointer(&x))
	for {
		typ := LoadPointer(&vp.typ)
		if typ == nil {
			// Attempt to start first store.
			// Disable preemption so that other goroutines can use
			// active spin wait to wait for completion; and so that
			// GC does not see the fake type accidentally.
			runtime_procPin()
			if !CompareAndSwapPointer(&vp.typ, nil, unsafe.Pointer(^uintptr(0))) {
				runtime_procUnpin()
				continue
			}
			// Complete first store.
			StorePointer(&vp.data, xp.data)
			StorePointer(&vp.typ, xp.typ)
			runtime_procUnpin()
			return
		}
		if uintptr(typ) == ^uintptr(0) {
			// First store in progress. Wait.
			// Since we disable preemption around the first store,
			// we can wait with active spinning.
			continue
		}
		// First store completed. Check type and overwrite data.
		if typ != xp.typ {
			panic("sync/atomic: store of inconsistently typed value into Value")
		}
		StorePointer(&vp.data, xp.data)
		return
	}
}

// Disable/enable preemption, implemented in runtime.
func runtime_procPin()
func runtime_procUnpin()
```
这段代码中的注释集已经告诉了我们，调用`Store`方法写入的类型必须与愿类型相同，不一致便会发生panic。接下来分析代码实现：
1. 首先判断条件写入参数不能为`nil`，否则触发`panic`
2. 通过使用`unsafe.Pointer`将`oldValue`和`newValue`转换成`ifaceWords`类型。方便我们获取他的原始类型(typ)和值(data).
3. 为了保证原子性，所以这里使用一个`for`换来处理，当已经有`Store`正在进行写入时，会进行等待.
4. 如果还没写入过数据，那么获取不到原始类型，就会开始第一次写入操作，这里会把先调用`runtime_procPin()`方法禁止调度器对当前 goroutine 的抢占（preemption），这样也可以防止`GC`线程看到一假类型。
5. 调用`CAS`方法来判断当前地址是否有被抢占，这里大家可能对`unsafe.Pointer(^uintptr(0))`这一句话有点不明白，因为是第一个写入数据，之前是没有数据的，所以通过这样一个中间值来做判断，如果失败就会解除抢占锁，解除禁止调度器，继续循环等待.
6. 设置中间值成功后，我们接下来就可以安全的把`v`设为传入的新值了，这里会先写入值，在写入类型(typ)，因为我们会根据ty来做完成判断。
7. 第一次写入没完成，我们还会通过`uintptr(typ) == ^uintptr(0)`来进行判断，因为还是第一次放入的中间类型，他依然会继续等待第一次完成。
8. 如果第一次写入完成，会检查上一次写入的类型与这次写入的类型是否一致，不一致则会抛出`panic`.

这里代码量没有多少，相信大家一定看懂了吧～。



### `Value`的读操作
先看一下代码：
```go
// Load returns the value set by the most recent Store.
// It returns nil if there has been no call to Store for this Value.
func (v *Value) Load() (x interface{}) {
	vp := (*ifaceWords)(unsafe.Pointer(v))
	typ := LoadPointer(&vp.typ)
	if typ == nil || uintptr(typ) == ^uintptr(0) {
		// First store not yet completed.
		return nil
	}
	data := LoadPointer(&vp.data)
	xp := (*ifaceWords)(unsafe.Pointer(&x))
	xp.typ = typ
	xp.data = data
	return
}
```
读取操作的代码就很简单了：
1.第一步使用`unsafe.Pointer`将`oldValue`转换成`ifaceWords`类型，然后获取他的类型，如果没有类型或者类型出去中间值，那么说明现在还没数据或者第一次写入还没有完成。
2. 通过检查后，调用`LoadPointer`方法可以获取他的值，然后构造一个新`interface`的`typ`和`data`返回。


### 小彩蛋

前面我们在说CAS时，说到了`ABA`问题，所以我就写了`demo`试一试`Go`标准库`atomic.CompareAndSwapXXX`方法是否有解决这个问题，看运行结果是没有，所以这里大家使用的时候要注意一下(虽然我也没想到什么现在什么业务场景会出现这个问题，但是还是要注意一下，需要自己评估)。

```go
func main()  {
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
```

## 总结

原子操作是并发编程的一个基础，也是为我学习`sync.once`打基础，好啦，现在你们应该知道下篇文章的内容是什么啦，敬请期待～。

**好啦，这篇文章就到这里啦，素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：加我vx拉你入群，或者公众号获取入群二维码**

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**翻译了一份Machinery中文文档，会定期进行维护，有需要的小伙伴们后台回复[machinery]即可获取。**

**我是asong，一名普普通通的程序猿，让我们一起慢慢变强吧。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [machinery-go异步任务队列](https://mp.weixin.qq.com/s/4QG69Qh1q7_i0lJdxKXWyg)
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

