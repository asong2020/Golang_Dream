`欢迎大家点击上方文字「Golang梦工厂」关注公众号，设为星标，第一时间接收推送文章。`

## 前言

> 哈喽，大家好，我是`asong`。`Go`语言中，`goroutine`的创建成本很低，调度效率很高，人称可以开几百几千万个`goroutine`，但是真正开几百几千万个`goroutine`就不会有任何影响吗？本文我们就一起来看一看`goroutine`是否有数量限制并介绍几种正确使用`goroutine`的姿势～。



## 现状

在`Go`语言中，`goroutine`的创建成本很低，调度效率高，`Go`语言在设计时就是按以数万个`goroutine`为规范进行设计的，数十万个并不意外，但是`goroutine`在内存占用方面确实具有有限的成本，你不能创造无限数量的它们，比如这个例子：

```go
ch := generate() 
go func() { 
        for range ch { } 
}()
```

这段代码通过`generate()`方法获得一个`channel`，然后启动一个`goroutine`一直去处理这个`channel`的数据，这个`goroutine`什么时候会退出？答案是不确定，`ch`是由函数`generate()`来决定的，所以有可能这个`goroutine`永远都不会退出，这就有可能会引发内存泄漏。

`goroutine`就是`G-P-M`调度模型中的`G`，我们可以把`goroutine`看成是一种协程，创建`goroutine`也是有开销的，但是开销很小，初始只需要`2-4k`的栈空间，当`goroutine`数量越来越大时，同时存在的`goroutine`也越来越多时，程序就隐藏内存泄漏的问题。看一个例子：

```go
func main()  {
	for i := 0; i < math.MaxInt64; i++ {
		go func(i int) {
			time.Sleep(5 * time.Second)
		}(i)
	}
}
```

大家可以在自己的电脑上运行一下这个程序，观察一下`CPU`和内存占用情况，我说下我运行后的现象：

![截屏2021-08-22 下午12.40.49](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-08-22%20%E4%B8%8B%E5%8D%8812.40.49.png)![截屏2021-08-22 下午12.41.05](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-08-22%20%E4%B8%8B%E5%8D%8812.41.05.png)

- CPU使用率疯狂上涨
- 内存占用率也不断上涨
- 运行一段时间后主进程崩溃了。。。

因此每次在编写`GO`程序时，都应该仔细考虑一个问题：

> 您将要启动的`goroutine`将如何以及在什么条件下结束？

接下来我们就来介绍几种方式可以控制`goroutine`和`goroutine`的数量。

## 控制`goroutine`的方法

### `Context`包

Go 语言中的每一个请求的都是通过一个单独的 `goroutine` 进行处理的，`HTTP/RPC` 请求的处理器往往都会启动新的` Goroutine` 访问数据库和 `RPC` 服务，我们可能会创建多个` goroutine` 来处理一次请求，而 `Context` 的主要作用就是在不同的 `goroutine` 之间同步请求特定的数据、取消信号以及处理请求的截止日期。

`Context`包主要衍生了四个函数：

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
func WithValue(parent Context, key, val interface{}) Context
```

使用这四个函数我们对`goroutine`进行控制，具体展开就不再本文说了，我们以`WithCancel`方法写一个例子：

```go
func main()  {
	ctx,cancel := context.WithCancel(context.Background())
	go Speak(ctx)
	time.Sleep(10*time.Second)
	cancel()
	time.Sleep(2 * time.Second)
	fmt.Println("bye bye!")
}

func Speak(ctx context.Context)  {
	for range time.Tick(time.Second){
		select {
		case <- ctx.Done():
			fmt.Println("asong哥，我收到信号了，要走了，拜拜！")
			return
		default:
			fmt.Println("asong哥，你好帅呀～balabalabalabala")
		}
	}
}
```

运行结果：

```go
asong哥，你好帅呀～balabalabalabala
# ....... 省略部分
asong哥，我收到信号了，要走了，拜拜！
bye bye!
```

这里我们使用`withCancel`创建了一个基于`Background`的`ctx`，然后启动了一个`goroutine`每隔`1s`夸我一句，`10s`后在主`goroutine`中发送取消新信号，那么启动的`goroutine`在检测到信号后就会取消退出。



### `channel` 

我们知道`channel`是用于`goroutine`的数据通信，在`Go`中通过`goroutine+channel`的方式，可以简单、高效地解决并发问题。上面我们介绍了使用`context`来达到对`goroutine`的控制，实际上`context`的内部实现也是使用的`channel`，所以有时候为了实现方便，我们可以直接通过`channel+select`或者`channel+close`的方式来控制`goroutine`的退出，我们分别来一写一个例子：

- `channel+select`

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

上面的例子是计算斐波那契数列的结果，我们使用两个`channel`，一个`channel`用来传输数据，另外一个`channel`用来做结束信号，这里我们使用的是`select`的阻塞式的收发操作，直到有一个`channel`发生状态改变，我们也可以在`select`中使用`default`语句，那么`select`语句在执行时会遇到这两种情况：

- 当存在可以收发的`Channel`时，直接处理该`Channel` 对应的 `case`；
- 当不存在可以收发的`Channel` 时，执行 `default` 中的语句；

建议大家使用带`default`的方式，因为在一个`nil channel`上的操作会一直被阻塞，如果没有`default case`,只有`nil channel`的`select`会一直被阻塞。



- `channel+close`

`channel`可以单个出队，也可以循环出队，因为我们可以使用`for-range`循环处理`channel`，`range ch`会一直迭代到`channel`被关闭，根据这个特性，我们也可做到对`goroutine`的控制：

```go
func main()  {
	ch := make(chan int, 10)
	go func() {
		for i:=0; i<10;i++{
			ch <- i
		}
		close(ch)
	}()
	go func() {
		for val := range ch{
			fmt.Println(val)
		}
		fmt.Println("receive data over")
	}()
	time.Sleep(5* time.Second)
	fmt.Println("program over")
}

```



如果对`channel`不熟悉的朋友可以看一下我之前的文章：[学习channel设计：从入门到放弃](https://mp.weixin.qq.com/s/E2XwSIXw1Si1EVSO1tMW7Q)



## 控制`goroutine`的数量

我们可以通过以下方式达到控制`goroutine`数量的目的，不过本身`Go`的`goroutine`就已经很轻量了，所以控制`goroutine`的数量还是要根据具体场景分析，并不是所有场景都需要控制`goroutine`的数量的，一般在并发场景我们会考虑控制`goroutine`的数量，接下来我们来看一看如下几种方式达到控制`goroutine`数量的目的。

### 协程池

写 `go` 并发程序的时候如果程序会启动大量的` goroutine` ，势必会消耗大量的系统资源（内存，CPU），所以可以考虑使用`goroutine`池达到复用`goroutine`，节省资源，提升性能。也有一些开源的协程池库，例如：`ants`、`go-playground/pool`、`jeffail/tunny`等，这里我们看`ants`的一个官方例子：

```go
var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

func demoFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

func main() {
	defer ants.Release()

	runTimes := 1000

	// Use the common pool.
	var wg sync.WaitGroup
	syncCalculateSum := func() {
		demoFunc()
		wg.Done()
	}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(syncCalculateSum)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks.\n")

	// Use the pool with a function,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc(i)
		wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(int32(i))
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
}
```

这个例子其实就是计算大量整数和的程序，这里通过`ants.NewPoolWithFunc()`创建了一个 `goroutine` 池。第一个参数是池容量，即池中最多有 `10` 个` goroutine`。第二个参数为每次执行任务的函数。当我们调用`p.Invoke(data)`的时候，`ants`池会在其管理的 `goroutine` 中找出一个空闲的，让它执行函数`taskFunc`，并将`data`作为参数。

具体这个库的设计就不详细展开了，后面会专门写一篇文章来介绍如何设计一个协程池。



### 信号量`Semaphore`

`Go`语言的官方扩展包为我们提供了一个基于权重的信号量[Semaphore](https://github.com/golang/sync/blob/master/semaphore/semaphore.go)，我可以根据信号量来控制一定数量的 `goroutine` 并发工作，官方也给提供了一个例子：[workerPool](https://pkg.go.dev/golang.org/x/sync/semaphore#example-package-WorkerPool)，代码有点长就不在这里贴了，我们来自己写一个稍微简单点的例子：

```go
const (
	Limit = 3  // 同时运行的goroutine上限
	Weight = 1 // 信号量的权重
)
func main() {
	names := []string{
		"asong1",
		"asong2",
		"asong3",
		"asong4",
		"asong5",
		"asong6",
		"asong7",
	}

	sem := semaphore.NewWeighted(Limit)
	var w sync.WaitGroup
	for _, name := range names {
		w.Add(1)
		go func(name string) {
			sem.Acquire(context.Background(), Weight)
			fmt.Println(name)
			time.Sleep(2 * time.Second) // 延时能更好的体现出来控制
			sem.Release(Weight)
			w.Done()
		}(name)
	}
	w.Wait()

	fmt.Println("over--------")
}
```

上面的例子我们使用 `NewWeighted()` 函数创建一个并发访问的最大资源数，也就是同时运行的`goroutine`上限为`3`，使用`Acquire`函数来获取指定个数的资源，如果当前没有空闲资源可用，则当前`goroutine`将陷入休眠状态，最后使用`release`函数释放已使用资源数量（计数器）进行更新减少，并通知其它 `waiters`。



### `channel+waitgroup`实现

这个方法我是在煎鱼大佬的一篇文章学到的：[来，控制一下Goroutine的并发数量](https://segmentfault.com/a/1190000017956396)

主要实现原理是利用`waitGroup`做并发控制，利用`channel`可以在`goroutine`之间进行数据通信，通过限制`channel`的队列长度来控制同时运行的`goroutine`数量，例子如下：

```go
func main()  {
	count := 9 // 要运行的goroutine数量
	limit := 3 // 同时运行的goroutine为3个
	ch := make(chan bool, limit)
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i:=0; i < count; i++{
		go func(num int) {
			defer wg.Done()
			ch <- true // 发送信号
			fmt.Printf("%d 我在干活 at time %d\n",num,time.Now().Unix())
			time.Sleep(2 * time.Second)
			<- ch // 接收数据代表退出了
		}(i)
	}
	wg.Wait()
}
```

这种实现方式真的妙，与信号量的实现方式基本相似，某些场景大家也可以考虑使用这种方式来达到控制`goroutine`的目的，不过最好封装一下，要不有点丑陋，感兴趣的可以看一下煎鱼大佬是怎么封装的：https://github.com/eddycjy/gsema/blob/master/sema.go



## 总结

本文主要目的是介绍控制`goroutine`的几种方式、控制`goroutine`数量的几种方式，`goroutine`的创建成本低、效率高带来了很大优势，同时也会有一些弊端，这就需要我们在实际开发中根据具体场景选择正确的方式使用`goroutine`，本文介绍的技术方案也可能是片面的，如果你有更好的方式可以在评论区中分享出来，我们大家一起学习学习～。

文中代码已经上传github，欢迎star：https://github.com/asong2020/Golang_Dream/tree/master/code_demo/goroutine_demo

**素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！我是`asong`，我们下期见。**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E7%99%BD%E8%89%B2%E7%89%88-20210717170231906-20210801174715998.png)

推荐往期文章：

- [学习channel设计：从入门到放弃](https://mp.weixin.qq.com/s/E2XwSIXw1Si1EVSO1tMW7Q)
- [详解内存对齐](https://mp.weixin.qq.com/s/ig8LDNdpflEBWlypU1NRhw)
- [Go语言中new和make你使用哪个来分配内存？](https://mp.weixin.qq.com/s/xNdnVXxC5Ji2ApgbfpRaXQ)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/yJ05a6pNxr_G72eiWTJ-rw)
- [面试官：小松子来聊一聊内存逃逸](https://mp.weixin.qq.com/s/MepbrrSlGVhNrEkTQhfhhQ)
- [面试官：你能聊聊string和[]byte的转换吗？](https://mp.weixin.qq.com/s/jztwFH6thFdcySzowXOH_Q)
- [面试官：两个nil比较结果是什么？](https://mp.weixin.qq.com/s/CNOLLLRzHomjBnbZMnw0Gg)
- [并发编程包之 errgroup](https://mp.weixin.qq.com/s/NcrENqRyK9dYrOBBI0SGkA)

