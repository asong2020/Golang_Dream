## 前言

> 哈喽，大家好，我是`asong`。终于回归了，停更了两周了，这两周一直在搞留言号的事，经过漫长的等待，终于搞定了。兄弟们，以后就可以在留言区尽情开喷了，只要你敢喷，我就敢精选🐶。(因为发生了账号迁移，需点击右上角重新添加星标，优质文章第一时间获取！)
>
> 今天给大家带来的是`Go`语言中的`channel`。`Go`语言从出世以来就以高并发著称，得益于其`Goroutine`的设计，`Goroutine`也就是一个可执行的轻量级协程，有了`Goroutine`我们可以轻松的运行协程，但这并不能满足我们的需求，我们往往还希望多个线程/协程是能够通信的，`Go`语言为了支持多个`Goroutine`通信，设计了`channel`，本文我们就一起从`GO1.15`的源码出发，看看`channel`到底是如何设计的。
>
> 好啦，开往幼儿园的列车就要开了，朋友们系好安全带，我要开车啦🐶



## 什么是channel

通过开头的介绍我们可以知道`channel`是用于`goroutine`的数据通信，在`Go`中通过`goroutine+channel`的方式，可以简单、高效地解决并发问题。我们先来看一下简单的示例：

```go
func GoroutineOne(ch chan <-string)  {
	fmt.Println("GoroutineOne running")
	ch <- "asong真帅"
	fmt.Println("GoroutineOne end of the run")
}

func GoroutineTwo(ch <- chan string)  {
	fmt.Println("GoroutineTwo running")
	fmt.Printf("女朋友说：%s\n",<-ch)
	fmt.Println("GoroutineTwo end of the run")
}


func main()  {
	ch := make(chan string)
	go GoroutineOne(ch)
	go GoroutineTwo(ch)
	time.Sleep(3 * time.Second)
}
// 运行结果
// GoroutineOne running
// GoroutineTwo running
// 女朋友说：asong真帅
// GoroutineTwo end of the run
// GoroutineOne end of the run
```

这里我们运行了两个`Goroutine`，在`GoroutineOne`中我们向`channel`中写入数据，在`GoroutineTwo`中我们监听`channel`，直到读取到"asong真帅"。我们可以画一个简单的图来表明一下这个顺序：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-06-30%20%E4%B8%8A%E5%8D%889.27.21.png)

上面的例子是对无缓冲`channel`的一个简单应用，其实`channel`的使用语法还是挺多的，下面且听我慢慢道来，毕竟是从入门到放弃嘛，那就先从入门开始。



## 入门`channel`

### `channel`类型

`channel`有三种类型的定义，分别是：`chan`、`chan <-`、`<- chan`，可选的`<-`代表`channel`的方向，如果我们没有指定方向，那么`channel`就是双向的，既可以接收数据，也可以发送数据。

```go
chan T // 接收和发送类型为T的数据
chan<- T // 只可以用来发送 T 类型的数据
<- chan T // 只可以用来接收 T 类型的数据
```

### 创建`channel`

我们可以使用`make`初始化`channel`，可以创建两种两种类型的`channel`：无缓冲的`channel`和有缓冲的`channel`。

示例：

```go
ch_no_buffer := make(chan int)
ch_no_buffer := make(chan int, 0)
ch_buffer := make(chan int, 100)
```

没有设置容量或者容量设置为`0`，则说明`channel`没有缓存，此时只有发送方和接收方都准备好后他们才可以进行通讯，否则就是一直阻塞。如果容量设置大于`0`，那就是一个带缓冲的`channel`，发送方只有`buffer`满了之后才会阻塞，接收方只有缓存空了才会阻塞。

**注意：未初始化（为nil）的`channel`是不可以通信的**

```go
func main()  {
	var ch chan string
	ch <- "asong真帅"
	fmt.Println(<- ch)
}
// 运行报错
fatal error: all goroutines are asleep - deadlock!
goroutine 1 [chan send (nil chan)]:
```


### `channel`入队

`channel`的入队定义如下：

```go
"channel" <- "要入队的值（可以是表达式）"
```

在无缓冲的`channel`中，只有在出队方准备好后，`channel`才会入队，否则一直阻塞着，所以说无缓冲`channel`是同步的。

在有缓冲的`channel`中，缓存未满时，就会执行入队操作。

**向`nil`的`channel`中入队会一直阻塞，导致死锁。**

### `channel`单个出队

`channel`的单个出队定义如下：

```go
<- "channel"
```

无论是有无缓冲的`channel`在接收不到数据时都会阻塞，直到有数据可以接收。

**从`nil`的`channel`中接收数据会一直阻塞。**

`channel`的出队还有一种非阻塞写法，定义如下：

```go
val, ok := <-ch
```

这么写可以判断当前`channel`是否关闭，如果这个`channel`被关闭了，`ok`会被设置为`false`，`val`就是零值。



### `channel`循环出队

我们可以使用`for-range`循环处理`channel`。

```go
func main()  {
	ch := make(chan int,10)
	go func() {
		for i:=0;i<10;i++{
			ch <- i
		}
		close(ch)
	}()
	for val := range ch{
		fmt.Println(val)
	}
	fmt.Println("over")
}
```

`range ch`会一直迭代到`channel`被关闭。在使用有缓冲`channel`时，配合`for-range`是一个不错的选择。



### 配合`select`使用

`Go`语言中的`select`能够让`Goroutine`同时等待多个`channel`读或者写，在`channel`状态未改变之前，`select`会一直阻塞当前线程或`Goroutine`。先看一个例子：

```go
func fibonacci(ch chan int, done chan struct{}) {
	x, y := 0, 1
	for {
		select {
		case ch <- x:
			x, y = y, x+y
		case <-done:
			fmt.Println("over")
			return
		}
	}
}
func main() {
	ch := make(chan int)
	done := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-ch)
		}
		done <- struct{}{}
	}()
	fibonacci(ch, done)
}
```

`select`与`switch`具有相似的控制结构，与`switch`不同的是，`select`中的`case`中的表达式必须是`channel`的收发操作，当`select`中的两个`case`同时被触发时，会随机执行其中的一个。为什么是随机执行的呢？随机的引入就是为了避免饥饿问题的发生，如果我们每次都是按照顺序依次执行的，若两个`case`一直都是满足条件的，那么后面的`case`永远都不会执行。

上面例子中的`select`用法是阻塞式的收发操作，直到有一个`channel`发生状态改变。我们也可以在`select`中使用`default`语句，那么`select`语句在执行时会遇到这两种情况：

- 当存在可以收发的` Channel `时，直接处理该` Channel` 对应的 `case`；
- 当不存在可以收发的` Channel` 时，执行 `default` 中的语句；

**注意：`nil channel`上的操作会一直被阻塞，如果没有`default case`,只有`nil channel`的`select`会一直被阻塞。**

### 关闭`channel`

内建的`close`方法可以用来关闭`channel`。如果`channel`已经关闭，不可以继续发送数据了，否则会发生`panic`，但是从这个关闭的`channel`中不但可以读取出已发送的数据，还可以不断的读取零值。

```go
func main()  {
	ch := make(chan int, 10)
	ch <- 10
	ch <- 20
	close(ch)
	fmt.Println(<-ch) //1
	fmt.Println(<-ch) //2
	fmt.Println(<-ch) //0
	fmt.Println(<-ch) //0
}
```



## `channel`基本设计思想

`channel`设计的基本思想是：**不要通过共享内存来通信，而是通过通信来实现共享内存（Do not communicate by sharing memory; instead, share memory by communicating）**。

这个思想大家是否理解呢？我在这里分享一下我的理解(查找资料+个人理解)，有什么不对的，留言区指正或开喷！

> 什么是使用共享内存来通信？其实就是多个线程/协程使用同一块内存，通过加锁的方式来宣布使用某块内存，通过解锁来宣布不再使用某块内存。
>
> 什么是通过通信来实现共享内存？其实就是把一份内存的开销变成两份内存开销而已，再说的通俗一点就是，我们使用发送消息的方式来同步信息。
>
> 为什么鼓励使用通过通信来实现共享内存？使用发送消息来同步信息相比于直接使用共享内存和互斥锁是一种更高级的抽象，使用更高级的抽象能够为我们在程序设计上提供更好的封装，让程序的逻辑更加清晰；其次，消息发送在解耦方面与共享内存相比也有一定优势，我们可以将线程的职责分成生产者和消费者，并通过消息传递的方式将它们解耦，不需要再依赖共享内存。
>
> 对于这个理解更深的文章，建议读一下这篇文章：[为什么使用通信来共享内存](https://draveness.me/whys-the-design-communication-shared-memory/)

`channel`在设计上本质就是一个有锁的环形队列，包括发送方队列、接收方队列、互斥锁等结构，下面我就一起从源码出发，剖析这个有锁的环形队列是怎么设计的！



## 源码剖析

### 数据结构

在`src/runtime/chan.go`中我们可以看到`hchan`的结构如下：

```go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters
	lock mutex
}
```

我们来解释一下`hchan`中每个字段都是什么意思：

- `qcount`：循环数组中的元素数量
- `dataqsiz`：循环数组的长度
- `buf`：只针对有缓冲的`channel`，指向底层循环数组的指针
- `elemsize`：能够接收和发送的元素大小
- `closed`：`channel`是否关闭标志
- `elemtype`：记录`channel`中元素的类型
- `sendx`：已发送元素在循环数组中的索引
- `sendx`：已接收元素在循环数组中的索引
- `recvq`：等待接收的`goroutine`队列
- `senq`：等待发送的`goroutine`队列
- `lock`：互斥锁，保护`hchan`中的字段，保证读写`channel`的操作都是原子的。

这个结构结合上面那个图理解就更清晰了：

- `buf`是指向底层的循环数组，`dataqsiz`就是这个循环数组的长度，`qcount`就是当前循环数组中的元素数量，缓冲的`channel`才有效。
- `elemsize`和`elemtype`就是我们创建`channel`时设置的容量大小和元素类型。
- `sendq`、`recvq`是一个双向链表结构，分别表示被阻塞的`goroutine`链表，这些 goroutine 由于尝试读取 `channel` 或向 `channel` 发送数据而被阻塞。

对于上面的描述，我们可以画出来这样的一个理解图：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-06-26%20%E4%B8%8B%E5%8D%884.45.09.png)

### `channel`的创建

前面介绍`channel`入门的时候我们就说到了，我们使用`make`进行创建，`make`在经过编译器编译后对应的`runtime.makechan`或`runtime.makechan64`。为什么会有这个区别呢？先看一下代码：

```go
// go 1.15.7
func makechan64(t *chantype, size int64) *hchan {
	if int64(int(size)) != size {
		panic(plainError("makechan: size out of range"))
	}

	return makechan(t, int(size))
}
```

`runtime.makechan64`本质也是调用的`makechan`方法，只不过多了一个数值溢出的校验。`runtime.makechan64`是用于处理缓冲区大于2的32方，所以这两个方法会根据传入的参数类型和缓冲区大小进行选择。大多数情况都是使用`makechan`。我们只需要分析`makechan`函数就可以了。

```go
func makechan(t *chantype, size int) *hchan {
	elem := t.elem
	// 对发送元素进行限制 1<<16 = 65536
	if elem.size >= 1<<16 {
		throw("makechan: invalid channel element type")
	}
  // 检查是否对齐
	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
		throw("makechan: bad alignment")
	}
  // 判断是否会发生内存溢出
	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
	if overflow || mem > maxAlloc-hchanSize || size < 0 {
		panic(plainError("makechan: size out of range"))
	}
  // 构造hchan对象
	var c *hchan
	switch {
  // 说明是无缓冲的channel
	case mem == 0:
		// Queue or element size is zero.
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		// Race detector uses this location for synchronization.
		c.buf = c.raceaddr()
  // 元素类型不包含指针，只进行一次内存分配
	// 如果hchan结构体中不含指针，gc就不会扫描chan中的元素，所以我们只需要分配
  // "hchan 结构体大小 + 元素大小*个数" 的内存
	case elem.ptrdata == 0:
		// Allocate hchan and buf in one call.
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
  // 元素包含指针，进行两次内存分配操作
	default:
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}
	// 初始化hchan中的对象
	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)
	lockInit(&c.lock, lockRankHchan)

	if debugChan {
		print("makechan: chan=", c, "; elemsize=", elem.size, "; dataqsiz=", size, "\n")
	}
	return c
}
```

注释我都添加上了，应该很容易懂吧，这里在特殊说一下分配内存这块的内容，其实归一下类，就只有两块：

- 分配一次内存：若创建的`channel`是无缓冲的，或者创建的有缓冲的`channel`中存储的类型不存在指针引用，就会调用一次`mallocgc`分配一段连续的内存空间。
- 分配两次内存：若创建的有缓冲`channel`存储的类型存在指针引用，就会连同`hchan`和底层数组同时分配一段连续的内存空间。

因为都是调用`mallocgc`方法进行内存分配，所以`channel`都是在堆上创建的，会进行垃圾回收，不关闭`close`方法也是没有问题的（但是想写出漂亮的代码就不建议你这么做了）。



### `channel`入队

`channel`发送数据部分的代码经过编译器编译后对应的是`runtime.chansend1`，其调用的也是`runtime.chansend`方法：

```go
func chansend1(c *hchan, elem unsafe.Pointer) {
	chansend(c, elem, true, getcallerpc())
}
```

我们主要分析一下`chansend`方法，代码有点长，我们分几个步骤来看这段代码：

- 前置检查
- 加锁/异常检查
- `channel`直接发送数据
- `channel`发送数据缓冲区有可用空间
- `channel`发送数据缓冲区无可用空间

#### 前置检查

```go
	if c == nil {
		if !block {
			return false
		}
		gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	if debugChan {
		print("chansend: chan=", c, "\n")
	}

	if raceenabled {
		racereadpc(c.raceaddr(), callerpc, funcPC(chansend))
	}
	if !block && c.closed == 0 && full(c) {
		return false
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}
```

这里最主要的检查就是判断当前`channel`是否为`nil`，往一个`nil`的`channel`中发送数据时，会调用`gopark`函数将当前执行的`goroutine`从`running`状态转入`waiting`状态，这让就会导致进程出现死锁，表象出`panic`事件。

紧接着会对非阻塞的`channel`进行一个上限判断，看看是否快速失败，这里相对于之前的版本做了改进，使用`full`方法来对`hchan`结构进行校验。

```go
func full(c *hchan) bool {
	if c.dataqsiz == 0 {
		return c.recvq.first == nil
	}
	return c.qcount == c.dataqsiz
}
```

这里快速失败校验逻辑如下：

- 若是 `qcount` 与 `dataqsiz` 大小相同（缓冲区已满）时，则会返回失败。
- 非阻塞且未关闭，同时底层数据 `dataqsiz` 大小为` 0`（无缓冲`channel`），如果接收方没准备好则直接返回失败。



#### 加锁/异常检查

```go
lock(&c.lock)

if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
}
```

前置校验通过后，在发送数据的逻辑执行之前会先为当前的`channel`加锁，防止多个协程并发修改数据。如果` Channel` 已经关闭，那么向该 `Channel `发送数据时会报` “send on closed channel” `错误并中止程序。

#### `channel`直接发送数据

直接发送数据是指 如果已经有阻塞的接收`goroutines`（即`recvq`中指向非空），那么数据将被直接发送给接收`goroutine`。

```go
if sg := c.recvq.dequeue(); sg != nil {
		//找到一个等待的接收器。我们将想要发送的值直接传递给接收者，绕过通道缓冲区(如果有的话)。
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
}
```

这里主要是调用`Send`方法，我们来看一下这个函数：

```go
func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
 // 静态竞争省略掉
  // elem是指接收到的值存放的位置
	if sg.elem != nil {
    // 调用sendDirect方法直接进行内存拷贝
    // 从发送者拷贝到接收者
		sendDirect(c.elemtype, sg, ep)
		sg.elem = nil
	}
  // 绑定goroutine
	gp := sg.g
  // 解锁
	unlockf()
	gp.param = unsafe.Pointer(sg)
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
  // 唤醒接收的 goroutine
	goready(gp, skip+1)
}
```

我们再来看一下`SendDirect`方法：

```go
func sendDirect(t *_type, sg *sudog, src unsafe.Pointer) {
	dst := sg.elem
	typeBitsBulkBarrier(t, uintptr(dst), uintptr(src), t.size)
	memmove(dst, src, t.size)
}
```

这里调用了`memove`方法进行内存拷贝，这里是从一个 `goroutine` 直接写另一个 `goroutine` 栈的操作，这样做的好处是减少了一次内存 `copy`：不用先拷贝到 `channel` 的` buf`，直接由发送者到接收者，没有中间商赚差价，效率得以提高，完美。



#### `channel`发送数据缓冲区有可用空间

接着往下看代码，判断`channel`缓冲区是否还有可用空间：

```go
// 判断通道缓冲区是否还有可用空间
if c.qcount < c.dataqsiz {
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			raceacquire(qp)
			racerelease(qp)
		}
		typedmemmove(c.elemtype, qp, ep)
  	// 指向下一个待发送元素在循环数组中的位置
		c.sendx++
   // 因为存储数据元素的结构是循环队列，所以当当前索引号已经到队末时，将索引号调整到队头
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
  	// 当前循环队列中存储元素数+1
		c.qcount++
   // 释放锁，发送数据完毕
		unlock(&c.lock)
		return true
}
```

这里的几个步骤还是挺好理解的，注释已经添加到代码中了，我们再来详细解析一下：

- 如果当前缓冲区还有可用空间，则调用`chanbuf`方法获取底层缓冲数组中`sendx`索引的元素指针值
- 调用`typedmemmove`方法将发送的值拷贝到缓冲区中
- 数据拷贝成功，`sendx`进行+1操作，指向下一个待发送元素在循环数组中的位置。如果下一个索引位置正好是循环队列的长度，那么就需要把所谓位置归0，因为这是一个循环环形队列。
- 发送数据成功后，队列元素长度自增，至此发送数据完毕，释放锁，返回结果即可。



#### `channel`发送数据缓冲区无可用空间

缓冲区空间也会有满了的时候，这是有两种方式可以选择，一种是直接返回，另外一种是阻塞等待。

直接返回的代码就很简单了，做一个简单的是否阻塞判断，不阻塞的话，直接释放锁，返回即可。

```go
if !block {
		unlock(&c.lock)
		return false
}
```

阻塞的话代码稍微长一点，我们来分析一下：

```go
  gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	mysg.elem = ep
	mysg.waitlink = nil
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.waiting = mysg
	gp.param = nil
	c.sendq.enqueue(mysg)
	atomic.Store8(&gp.parkingOnChan, 1)
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceEvGoBlockSend, 2)
  KeepAlive(ep)
```

首先通过调用`gettg`获取当前执行的`goroutine`，然后调用`acquireSudog`方法构造`sudog`结构体，然后设置待发送信息和`goroutine`等信息（`sudog` 通过 `g` 字段绑定 `goroutine`，而` goroutine` 通过` waiting `绑定 `sudog`，`sudog` 还通过 `elem` 字段绑定待发送元素的地址）；构造完毕后调用`c.sendq.enqueue`将其放入待发送的等待队列，最后调用`gopark`方法挂起当前的`goroutine`进入`wait`状态。

这里在最后调用了`KeepAlive`方法，很多人对这个比较懵逼，我来解释一下。这个方法就是为了保证待发送的数据处于活跃状态，也就是分配在堆上避免被GC。这里我在画一个图解释一下上面的绑定过程，更加深理解。

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-06-26%20%E4%B8%8B%E5%8D%885.20.19.png)

现在`goroutine`处于`wait`状态了，等待被唤醒，唤醒代码如下：

```go
 if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	if gp.param == nil {
		if c.closed == 0 {
			throw("chansend: spurious wakeup")
		}
    // 唤醒后channel被关闭了，直接panic
		panic(plainError("send on closed channel"))
	}
	gp.param = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
 // 去掉mysg上绑定的channel
	mysg.c = nil
  // 释放sudog
	releaseSudog(mysg)
	return true
```

唤醒的逻辑比较简单，首先判断`goroutine`是否还存在，不存在则抛出异常。唤醒后还有一个检查是判断当前`channel`是否被关闭了，关闭了则触发`panic`。最后我们开始取消`mysg`上的`channel`绑定和`sudog`的释放。

这里大家肯定好奇，怎么没有看到唤醒后执行发送数据动作？之所以有这个想法，就是我们理解错了。在上面我们已经使`goroutine`进入了`wait`状态，那么调度器在停止`g` 时会记录运行线程和方法内执行的位置，也就是这个`ch <- "asong"`位置，唤醒后会在这个位置开始执行，代码又开始重新执行了，但是我们之前进入`wait`状态的绑定是要解绑与释放的，否则下次进来就会出现问题喽。



### 接收数据

之前我们介绍过`channel`接收数据有两种方式，如下：

```go
val := <- ch
val, ok := <- ch
```

它们在经过编译器编译后分别对应的是`runtime.chanrecv1` 和 `runtime.chanrecv2`：

```go
//go:nosplit
func chanrecv1(c *hchan, elem unsafe.Pointer) {
	chanrecv(c, elem, true)
}

//go:nosplit
func chanrecv2(c *hchan, elem unsafe.Pointer) (received bool) {
	_, received = chanrecv(c, elem, true)
	return
}
```

其实都是调用`chanrecv`方法，所以我们只需要解析这个方法就可以了。接收部分的代码和接收部分的代码是相对应的，所以我们也可以分几个步骤来看这部分代码：

- 前置检查
- 加锁和提前返回
- `channel`直接接收数据
- `channel`缓冲区有数据
- `channel`缓冲区无数据



#### 前置检查

```go
if c == nil {
		if !block {
			return
}
		gopark(nil, nil, waitReasonChanReceiveNilChan, traceEvGoStop, 2)
		throw("unreachable")
}
if atomic.Load(&c.closed) == 0 {
			return
}
if empty(c) {
		if raceenabled {
				raceacquire(c.raceaddr())
		}
		if ep != nil {
				typedmemclr(c.elemtype, ep)
		}
			return true, false
  }
}

var t0 int64
if blockprofilerate > 0 {
	t0 = cputicks()
}
```

首先也会判断当前`channel`是否为`nil channel`，如果是`nil channel`且为非阻塞接收，则直接返回即可。如果是`nil channel`且为阻塞接收，则直接调用`gopark`方法挂起当前`goroutine`。

然后也会进行快速失败检查，这里只会对非阻塞接收的`channel`进行快速失败检查，检查规则如下：

```go
func empty(c *hchan) bool {
	// c.dataqsiz is immutable.
	if c.dataqsiz == 0 {
		return atomic.Loadp(unsafe.Pointer(&c.sendq.first)) == nil
	}
	return atomic.Loaduint(&c.qcount) == 0
}
```

当循环队列为 `0`且等待队列 `sendq `内没有 `goroutine` 正在等待或者缓冲区数组为空时，如果`channel`还未关闭，这说明没有要接收的数据，直接返回即可。如果`channel`已经关闭了且缓存区没有数据了，则会清理`ep`指针中的数据并返回。这里为什么清理`ep`指针呢？`ep`指针是什么？这个`ep`就是我们要接收的值存放的地址（`val := <-ch val`就是`ep`  ），即使`channel`关闭了，我们也可以接收零值。

#### 加锁和提前返回

```go
	lock(&c.lock)

	if c.closed != 0 && c.qcount == 0 {
		if raceenabled {
			raceacquire(c.raceaddr())
		}
		unlock(&c.lock)
		if ep != nil {
			typedmemclr(c.elemtype, ep)
		}
		return true, false
	}
```

前置校验通过后，在执行接收数据的逻辑之前会先为当前的`channel`加锁，防止多个协程并发接收数据。同样也会判断当前`channel`是否被关闭，如果`channel`被关闭了，并且缓存区没有数据了，则直接释放锁和清理`ep`中的指针数据，不需要再走接下来的流程。



#### `channel`直接接收数据

这一步与`channel`直接发送数据是对应的，当发现`channel`上有正在阻塞等待的发送方时，则直接进行接收。

```go
if sg := c.sendq.dequeue(); sg != nil {
		recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true, true
	}
```

等待发送队列里有`goroutine`存在，有两种可能：

- 非缓冲的`channel`
- 缓冲的`channel`，但是缓冲区满了

针对这两种情况，在`recv`方法中的执行逻辑是不同的：

```go
func recv(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
  // 非缓冲channel
	if c.dataqsiz == 0 {
    // 未忽略接收值
		if ep != nil {
			// 直接从发送方拷贝数据到接收方
			recvDirect(c.elemtype, sg, ep)
		}
	} else { // 有缓冲channel，但是缓冲区满了
    // 缓冲区满时，接收方和发送方游标重合了
    // 因为是循环队列，都是游标0的位置
    // 获取当前接收方游标位置下的值
		qp := chanbuf(c, c.recvx)
		// 未忽略值的情况下直接从发送方拷贝数据到接收方
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		// 将发送者数据拷贝到缓冲区中
		typedmemmove(c.elemtype, qp, sg.elem)
    // 自增到下一个待接收位置
		c.recvx++
    // 如果下一个待接收位置等于队列长度了，则下一个待接收位置为队头，因为是循环队列
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
    // 上面已经将发送者数据拷贝到缓冲区中了，所以缓冲区还是满的，所以发送方位置仍然等于接收方位置。
		c.sendx = c.recvx // c.sendx = (c.sendx+1) % c.dataqsiz
	}
	sg.elem = nil
  // 绑定发送方goroutine
	gp := sg.g
	unlockf()
	gp.param = unsafe.Pointer(sg)
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
  // 唤醒发送方的goroutine
	goready(gp, skip+1)
}
```

代码中的注释已经很清楚了，但还是想在解释一遍，这里主要就是分为两种情况：

- 非缓冲区`channel`：未忽略接收值时直接调用`recvDirect`方法直接从发送方的`goroutine`调用栈中将数据拷贝到接收方的`goroutine`。
- 带缓冲区的`channel`：首先调用`chanbuf`方法根据`recv`索引的位置读取缓冲区元素，并将其拷贝到接收方的内存地址；拷贝完毕后调整`sendx`和`recvx`索引位置。

最后别忘了还有一个操作就是调用`goready`方法唤醒发送方的`goroutine`可以继续发送数据了。



#### `channel`缓冲区有数据

我们接着往下看代码，若当前`channel`的缓冲区有数据时，代码逻辑如下：

```go
  // 缓冲channel，buf里有可用元素，发送方也可以正常发送
   if c.qcount > 0 {
     // 直接从循环队列中找到要接收的元素
		qp := chanbuf(c, c.recvx)
    // 未忽略接收值，直接把缓冲区的值拷贝到接收方中
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
     // 清理掉循环数组里相应位置的值
		typedmemclr(c.elemtype, qp)
     // 接收游标向前移动
		c.recvx++
     // 超过循环队列的长度时，接收游标归0（循环队列）
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
     // 循环队列中的数据数量减1
		c.qcount--
    // 接收数据完毕，释放锁
		unlock(&c.lock)
		return true, true
	}

	if !block {
		unlock(&c.lock)
		return false, false
	}
```

这段代码没什么难度，就不再解释一遍了。



#### `channel`缓冲区无数据

经过上面的步骤，现在可以确定目前这个`channel`既没有待发送的`goroutine`，并且缓冲区也没有数据。接下来就看我们是否阻塞等待接收数据了，也就有了如下判断：

```go
	if !block {
		unlock(&c.lock)
		return false, false
	}
```

非阻塞接收数据的话，直接返回即可；否则则进入阻塞接收模式：

```go
  gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	mysg.elem = ep
	mysg.waitlink = nil
	gp.waiting = mysg
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.param = nil
	c.recvq.enqueue(mysg)
	atomic.Store8(&gp.parkingOnChan, 1)
	gopark(chanparkcommit,  unsafe.Pointer(&c.lock), waitReasonChanReceive, traceEvGoBlockRecv, 2)
```

这一部分的逻辑基本与发送阻塞部分一模一样，大概逻辑就是获取当前的`goroutine`，然后构建`sudog`结构保存待接收数据的地址信息和`goroutine`信息，并将`sudog`加入等待接收队列，最后挂起当前`goroutine`，等待唤醒。

接下来的环境逻辑也没有特别要说的，与发送方唤醒部分一模一样，不懂的可以看前面。唤醒后的主要工作就是恢复现场，释放绑定信息。



## 关闭`channel`

使用`close`可以关闭`channel`，其经过编译器编译后对应的是`runtime.closechan`方法，详细逻辑我们通过注释到代码中：

```go
func closechan(c *hchan) {
  // 对一个nil的channel进行关闭会引发panic
	if c == nil {
		panic(plainError("close of nil channel"))
	}
  // 加锁
	lock(&c.lock)
  // 关闭一个已经关闭的channel也会引发channel
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("close of closed channel"))
	}
	// 关闭channnel标志
	c.closed = 1
 // Goroutine集合
	var glist gList

	// 接受者的 sudog 等待队列（recvq）加入到待清除队列 glist 中
	for {
		sg := c.recvq.dequeue()
		if sg == nil {
			break
		}
		if sg.elem != nil {
			typedmemclr(c.elemtype, sg.elem)
			sg.elem = nil
		}
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		gp := sg.g
		gp.param = nil
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		glist.push(gp)
	}

	// 发送方的sudog也加入到到待清除队列 glist 中
	for {
		sg := c.sendq.dequeue()
		if sg == nil {
			break
		}
    // 要关闭的goroutine，发送的值设为nil
		sg.elem = nil
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		gp := sg.g
		gp.param = nil
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		glist.push(gp)
	}
  // 释放了发送方和接收方后，释放锁就可以了。
	unlock(&c.lock)

	// 将所有 glist 中的 goroutine 状态从 _Gwaiting 设置为 _Grunnable 状态，等待调度器的调度。
  // 我们既然是从sendq和recvq中获取的goroutine，状态都是挂起状态，所以需要唤醒他们，走后面的流程。
	for !glist.empty() {
		gp := glist.pop()
		gp.schedlink = 0
		goready(gp, 3)
	}
}
```

这里逻辑还是比较简单，归纳总结一下：

- 一个为`nil`的`channel`不允许进行关闭
- 不可以重复关闭`channel`
- 获取当前正在阻塞的发送或者接收的`goroutine`，他们都处于挂起状态，然后进行唤醒。这是发送方不允许在向`channel`发送数据了，但是不影响接收方继续接收元素，如果没有元素，获取到的元素是零值。使用`val,ok := <-ch`可以判断当前`channel`是否被关闭。



## 总结

哇塞，开往幼儿园的车终于停了，小松子唠唠叨叨一路了，你们学会了吗？

我们从入门开始到最后的源码剖析，其实`channel`的设计一点也不复杂，源码也是很容易看懂的，本质就是维护了一个循环队列嘛，发送数据遵循FIFO（First In First Out）原语，数据传递依赖于内存拷贝。不懂的可以再看一遍，很容易理解的哦～。

最后我想说的是：`channel`内部也是使用互斥锁，那么`channel`和互斥锁谁更轻量呢？（**评论区我们一起探讨一下**）。

**素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！我是`asong`，我们下期见。**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![扫码_搜索联合传播样式-白色版](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E7%99%BD%E8%89%B2%E7%89%88.png)

推荐往期文章：

- [Go语言如何实现可重入锁？](https://mp.weixin.qq.com/s/S_EzyWZmFzzbBbxoSNe6Hw)
- [Go语言中new和make你使用哪个来分配内存？](https://mp.weixin.qq.com/s/xNdnVXxC5Ji2ApgbfpRaXQ)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/yJ05a6pNxr_G72eiWTJ-rw)
- [空结构体引发的大型打脸现场](https://mp.weixin.qq.com/s/aHwGWWmnDFkcw2cyw5jmgw)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/UJKBHm58TXi37v53iZP8xA)
- [面试官：两个nil比较结果是什么？](https://mp.weixin.qq.com/s/CNOLLLRzHomjBnbZMnw0Gg)
- [面试官：你能用Go写段代码判断当前系统的存储方式吗?](https://mp.weixin.qq.com/s/DWMqzOi7wf79DoUUAJnr1w)
- [面试中如果这样写二分查找](https://mp.weixin.qq.com/s/z7NIzrcVRhpoLUQdFAa8JQ)

