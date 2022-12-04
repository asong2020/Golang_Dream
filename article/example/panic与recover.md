## 前言

> 哈喽，大家好，我是`asong`，今天与大家来聊一聊`go`语言中的"throw、try.....catch{}"。如果你之前是一名`java`程序员，我相信你一定吐槽过`go`语言错误处理方式，但是这篇文章不是来讨论好坏的，我们本文的重点是带着大家看一看`panic`与`recover`是如何实现的。上一文我们讲解了[`defer`是如何实现的](https://mp.weixin.qq.com/s/FUmoBB8OHNSfy7STR0GsWw)，但是没有讲解与`defoer`紧密相连的`recover`，想搞懂`panic`与`recover`的实现也没那么简单，就放到这一篇来讲解了。废话不多说，直接开整。



## 什么是`panic`、`recover`

Go 语言中 `panic` 关键字主要用于主动抛出异常，类似 `java` 等语言中的 `throw` 关键字。`panic` 能够改变程序的控制流，调用 `panic` 后会立刻停止执行当前函数的剩余代码，并在当前 Goroutine 中递归执行调用方的 `defer`；

Go 语言中 `recover` 关键字主要用于捕获异常，让程序回到正常状态，类似 `java` 等语言中的 `try ... catch` 。`recover` 可以中止 `panic` 造成的程序崩溃。它是一个只能在 `defer` 中发挥作用的函数，在其他作用域中调用不会发挥作用；

`recover`只能在`defer`中使用这个在标准库的注释中已经写明白了，我们可以看一下：

```go
// The recover built-in function allows a program to manage behavior of a
// panicking goroutine. Executing a call to recover inside a deferred
// function (but not any function called by it) stops the panicking sequence
// by restoring normal execution and retrieves the error value passed to the
// call of panic. If recover is called outside the deferred function it will
// not stop a panicking sequence. In this case, or when the goroutine is not
// panicking, or if the argument supplied to panic was nil, recover returns
// nil. Thus the return value from recover reports whether the goroutine is
// panicking.
func recover() interface{}
```

这里有一个要注意的点就是`recover`必须要要在`defer`函数中使用，否则无法阻止`panic`。最好的验证方法是先写两个例子：

```go
func main()  {
	example1()
	example2()
}

func example1()  {
	defer func() {
		if err := recover(); err !=nil{
			fmt.Println(string(Stack()))
		}
	}()
	panic("unknown")
}

func example2()  {
	defer recover()
	panic("unknown")
}

func Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
```

运行我们会发现`example2()`方法的`panic`是没有被`recover`住的，导致整个程序直接`crash`了。这里大家肯定会有疑问，为什么直接写`recover()`就不能阻止`panic`了呢。我们在[详解defer实现机制(附上三道面试题，我不信你们都能做对)](https://mp.weixin.qq.com/s/FUmoBB8OHNSfy7STR0GsWw)讲解了`defer`实现原理，一个重要的知识点**`defer`将语句放入到栈中时，也会将相关的值拷贝同时入栈。**所以`defer recover()`这种写法在放入`defer`栈中时就已经被执行过了，`panic`是发生在之后，所以根本无法阻止住`panic`。




## 特性

上面我们简单的介绍了一下什么是`panic`与`recover`，下面我一起来看看他们有什么特性，避免我们踩坑。

- `recover`只有在`defer`函数中使用才有效，上面已经举例说明了，这里就不在赘述了。
- `panic`允许在`defer`中嵌套多次调用.程序多次调用 `panic` 也不会影响 `defer` 函数的正常执行，所以使用 `defer` 进行收尾工作一般来说都是安全的。写个例子验证一下：

```go
func example3()  {
	defer fmt.Println("this is a example3 for defer use panic")
	defer func() {
		defer func() {
			panic("panic defer 2")
		}()
		panic("panic defer 1")
	}()
	panic("panic example3")
}
// 运行结果
this is a example3 for defer use panic
panic: panic example3
        panic: panic defer 1
        panic: panic defer 2
.......... 省略
```

通过运行结果可以看出`panic`不会影响`defer`函数的使用，所以他是安全的。

- `panic`只会对当前`Goroutine`的`defer`有效，还记得我们上一文分析的`deferproc`函数吗？在`newdefer`中分配`_defer`结构体对象的时，会把分配到的对象链入当前 `goroutine`的`_defer` 链表的表头，也就是把延迟调用函数与调用方所在的`Goroutine`进行关联。因此当程序发生`panic`时只会调用当前 Goroutine 的延迟调用函数是没有问题的。写个例子验证一下：

```go
func main()  {
	go example4()
	go example5()
	time.Sleep(10 * time.Second)
}

func example4()  {
	fmt.Println("goroutine example4")
	defer func() {
		fmt.Println("test defer")
	}()
	panic("unknown")
}

func example5()  {

	defer fmt.Println("goroutine example5")
	time.Sleep(5 * time.Second)
}
// 运行结果
goroutine example4
test defer
panic: unknown
............. 省略部分代码
```

这里我开了两个协程，一个协程会发生`panic`，导致程序崩溃，但是只会执行自己所在`Goroutine`的延迟函数，所以正好验证了多个 `Goroutine` 之间没有太多的关联，一个 `Goroutine` 在 `panic` 时也不应该执行其他 `Goroutine` 的延迟函数。



## 典型应用

其实我们在实际项目开发中，经常会遇到`panic`问题， Go 的 `runtime` 代码中很多地方都调用了 `panic` 函数，对于不了解 Go 底层实现的新人来说，这无疑是挖了一堆深坑。我们在实际生产环境中总会出现`panic`，但是我们的程序仍能正常运行，这是因为我们的框架已经做了`recover`，他已经为我们兜住底，比如`gin`，我们看一看他是怎么做的。

先看代码部分吧：

```go
func Default() *Engine {
	debugPrintWARNINGDefault()
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}
// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc {
	return RecoveryWithWriter(DefaultErrorWriter)
}

// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func RecoveryWithWriter(out io.Writer) HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
			...................// 省略
			}
		}()
		c.Next()
	}
}
```
我们在使用`gin`时，第一步会初始化一个`Engine`实例，调用`Default`方法会把`recovery middleware`附上，`recovery`中使用了`defer`函数，通过`recover`来阻止`panic`，当发生`panic`时，会返回500错误码。这里有一个需要注意的点是只有主程序中的`panic`是会被自动`recover`的,协程中出现`panic`会导致整个程序`crash`。还记得我们上面讲的第三个特性嘛，**一个协程会发生`panic`，导致程序崩溃，但是只会执行自己所在`Goroutine`的延迟函数，所以正好验证了多个 `Goroutine` 之间没有太多的关联，一个 `Goroutine` 在 `panic` 时也不应该执行其他 `Goroutine` 的延迟函数。** 这就能解释通了吧， 所以为了程序健壮性，我们应该自己主动检查我们的协程程序，在我们的协程函数中添加`recover`是很有必要的，比如这样：

```go
func main()  {
		r := gin.Default()
		r.GET("/asong/test/go-panic", func(ctx *gin.Context) {
			go func() {
				defer func() {
					if err := recover();err != nil{
						fmt.Println(err)
					}
				}()
				panic("panic")
			}()
		})
		r.Run()
}
```
如果使用的`Gin`框架，切记要检查协程中是否会出现`panic`，否则线上将付出沉重的代价。**非常危险！！！**


## 源码解析

go-version: 1.15.3

我们先来写个简单的代码，看看他的汇编调用：

```go
func main()  {
	defer func() {
		if err:= recover();err != nil{
			fmt.Println(err)
		}
	}()
	panic("unknown")
}
```

执行`go tool compile -N -l -S main.go`就可以看到对应的汇编码了，我们截取部分片段分析：

<img src="https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/panic.png" style="zoom: 67%;" />

上面重点部分就是画红线的三处，第一步调用`runtime.deferprocStack`创建`defer`对象，这一步大家可能会有疑惑，我上一文忘记讲个这个了，这里先简单概括一下，`defer`总共有三种模型，编译一个函数里只会有一种`defer`模式

- 第一种，堆上分配(deferproc)，基本是依赖运行时来分配"_defer"对象并加入延迟参数。在函数的尾部插入`deferreturn`方法来消费`defer`link。
- 第二种，栈上分配(deferprocStack)，基本上跟堆差不多，只是分配方式改为在栈上分配，压入的函数调用栈存有`_defer`记录，编译器在`ssa`过程中会预留`defer`空间。
- 第三种，开放编码模式(open coded)，不过是有条件的，默认open-coded最多支持8个defer，超过则取消。在构建ssa时如发现gcflags有N禁止优化的参数 或者 return数量 * defer数量超过了 15不适用open-coded模式。并不能处于循环中。

按理说我们的版本是`1.15+`，应该使用开放编码模式呀，但是这里怎么还会在栈上分配？注意看呀，伙计们，我在汇编处理时禁止了编译优化，那肯定不会走开放编码模式呀，这个不是重点，我们接着分析上面的汇编。

第二个红线在程序发生`panic`时会调用`runtime.gopanic`，现在程序处于`panic`状态，在函数返回时调用`runtime.deferreturn`，也就是调用延迟函数处理。上面这一步是主程序执行部分，下面我们在看一下延迟函数中的执行：

<img src="https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/recover.png" style="zoom:67%;" />

这里最重点的就只有一个，调用`runtime.gorecover`，也就是在这一步，对主程序中的`panic`进行了恢复了，这就是`panic`与`recover`的执行过程，接下来我们就仔细分析一下`runtime.gopanic`、`runtime.gorecover`这两个方法是如何实现的！



### _panic结构

在讲[`defer`实现机制](https://mp.weixin.qq.com/s/FUmoBB8OHNSfy7STR0GsWw)时，我们一起看过`defer`的结构，其中有一个字段就是`_panic`，是触发`defer`的作用，我们来看看的`panic`的结构：

```go
type _panic struct {
	argp      unsafe.Pointer // pointer to arguments of deferred call run during panic; cannot move - known to liblink
	arg       interface{}    // argument to panic
	link      *_panic        // link to earlier panic
	pc        uintptr        // where to return to in runtime if this panic is bypassed
	sp        unsafe.Pointer // where to return to in runtime if this panic is bypassed
	recovered bool           // whether this panic is over
	aborted   bool           // the panic was aborted
	goexit    bool
}
```

简单介绍一下上面的字段：

- `argp`是指向`defer`调用时参数的指针。
- `arg`是我们调用`panic`时传入的参数
- `link`指向的是更早调用`runtime._panic`结构，也就是说`painc`可以被连续调用，他们之间形成链表
- `recovered` 表示当前`runtime._panic`是否被`recover`恢复
- `aborted`表示当前的`panic`是否被强行终止

上面的`pc`、`sp`、`goexit`我们单独讲一下，`runtime`包中有一个`Goexit`方法，`Goext`能够终止调用它的`goroutine`，其他的`goroutine`是不受影响的，`goexit`也会在终止`goroutine`之前运行所有延迟调用函数，`Goexit`不是一个`panic`，所以这些延迟函数中的任何`recover`调用都将返回`nil`。如果我们在主函数中调用了`Goexit`会终止该`goroutine`但不会返回`func main`。由于`func main`没有返回，因此程序将继续执行其他`gorountine`，直到所有其他`goroutine`退出，程序才会`crash`。写个简单的例子：

```go
func main()  {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		runtime.Goexit()
	}()
	go func() {
		for true {
			fmt.Println("test")
		}
	}()
	runtime.Goexit()
	fmt.Println("main")
	select {

	}
}
```

运行上面的例子你就会发现，即使在主`goroutine`中调用了`runtime.Goexit`，其他`goroutine `是没有任何影响的。所以结构中的`pc`、`sp`、`goexit`三个字段都是为了修复`runtime.Goexit`，这三个字段就是为了保证该函数的一定会生效，因为如果在`defer`中发生`panic`，那么`goexit`函数就会被取消，所以才有了这三个字段做保护。看这个例子：

```go
func main()  {
	maybeGoexit()
}
func maybeGoexit() {
	defer func() {
		fmt.Println(recover())
	}()
	defer panic("cancelled Goexit!")
	runtime.Goexit()
}
```

英语好的可以看一看这个：https://github.com/golang/go/issues/29226，这就是上面的一个例子，这里就不过多解释了，了解就好。

下面就开始我们的重点吧～。



### gopanic

gopanic的代码有点长，我们一点一点来分析：

- 第一部分，判断`panic`类型：

```go
gp := getg()
	if gp.m.curg != gp {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic on system stack")
	}

	if gp.m.mallocing != 0 {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic during malloc")
	}
	if gp.m.preemptoff != "" {
		print("panic: ")
		printany(e)
		print("\n")
		print("preempt off reason: ")
		print(gp.m.preemptoff)
		print("\n")
		throw("panic during preemptoff")
	}
	if gp.m.locks != 0 {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic holding locks")
	}
```

根据不同的类型判断当前发生`panic`错误，这里没什么多说的，接着往下看。

- 第二部分，确保每个`recover`都试图恢复当前协程中最新产生的且尚未恢复的`panic`

```go
var p _panic // 声明一个panic结构
	p.arg = e // 把panic传入的值赋给`arg`
	p.link = gp._panic // 指向runtime.panic结构
	gp._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

	atomic.Xadd(&runningPanicDefers, 1)

	// By calculating getcallerpc/getcallersp here, we avoid scanning the
	// gopanic frame (stack scanning is slow...)
	addOneOpenDeferFrame(gp, getcallerpc(), unsafe.Pointer(getcallersp()))

	for {
		d := gp._defer // 获取当前gorourine的 defer
		if d == nil {
			break // 如果没有defer直接退出了
		}

		// If defer was started by earlier panic or Goexit (and, since we're back here, that triggered a new panic),
		// take defer off list. An earlier panic will not continue running, but we will make sure below that an
		// earlier Goexit does continue running.
		if d.started {
			if d._panic != nil {
				d._panic.aborted = true
			}
			d._panic = nil
			if !d.openDefer {
				// For open-coded defers, we need to process the
				// defer again, in case there are any other defers
				// to call in the frame (not including the defer
				// call that caused the panic).
				d.fn = nil
				gp._defer = d.link
				freedefer(d)
				continue
			}
		}

		// Mark defer as started, but keep on list, so that traceback
		// can find and update the defer's argument frame if stack growth
		// or a garbage collection happens before reflectcall starts executing d.fn.
		d.started = true
    // Record the panic that is running the defer.
		// If there is a new panic during the deferred call, that panic
		// will find d in the list and will mark d._panic (this panic) aborted.
		d._panic = (*_panic)(noescape(unsafe.Pointer(&p)))
```

上面的代码不太好说的部分，我添加了注释，就不在这解释一遍了，直接看 `d.Started`部分，这里的意思是如果`defer`是由先前的`panic`或`Goexit`启动的(循环处理回到这里，这触发了新的`panic`)，将`defer`从列表中删除。早期的`panic`将不会继续运行，但我们将确保早期的Goexit会继续运行，代码中的`if d._panic != nil{d._panic.aborted =true}`就是确保将先前的`panic`终止掉，将`aborted`设置为`true`，在下面执行`recover`时保证`goexit`不会被取消。

- 第三部分，`defer`内联优化调用性能

```go
	if !d.openDefer {
				// For open-coded defers, we need to process the
				// defer again, in case there are any other defers
				// to call in the frame (not including the defer
				// call that caused the panic).
				d.fn = nil
				gp._defer = d.link
				freedefer(d)
				continue
			}

		done := true
		if d.openDefer {
			done = runOpenDeferFrame(gp, d)
			if done && !d._panic.recovered {
				addOneOpenDeferFrame(gp, 0, nil)
			}
		} else {
			p.argp = unsafe.Pointer(getargp(0))
			reflectcall(nil, unsafe.Pointer(d.fn), deferArgs(d), uint32(d.siz), uint32(d.siz))
		}
```

上面的代码都是截图片段，这些部分都是为了判断当前`defer`是否可以使用开发编码模式，具体怎么操作的就不展开了。



- 第四部分，`gopanic`中执行程序恢复

在第三部分进行`defer`内联优化选择时会执行调用延迟函数(reflectcall就是这个作用)，也就是会调用`runtime.gorecover`把`recoverd = true`，具体这个函数的操作留在下面讲，因为`runtime.gorecover`函数并不包含恢复程序的逻辑，程序的恢复是在`gopanic`中执行的。先看一下代码：

```go
		if p.recovered { // 在runtime.gorecover中设置为true
			gp._panic = p.link 
			if gp._panic != nil && gp._panic.goexit && gp._panic.aborted { 
				// A normal recover would bypass/abort the Goexit.  Instead,
				// we return to the processing loop of the Goexit.
				gp.sigcode0 = uintptr(gp._panic.sp)
				gp.sigcode1 = uintptr(gp._panic.pc)
				mcall(recovery)
				throw("bypassed recovery failed") // mcall should not return
			}
			atomic.Xadd(&runningPanicDefers, -1)

			if done {
				// Remove any remaining non-started, open-coded
				// defer entries after a recover, since the
				// corresponding defers will be executed normally
				// (inline). Any such entry will become stale once
				// we run the corresponding defers inline and exit
				// the associated stack frame.
				d := gp._defer
				var prev *_defer
				for d != nil {
					if d.openDefer {
						if d.started {
							// This defer is started but we
							// are in the middle of a
							// defer-panic-recover inside of
							// it, so don't remove it or any
							// further defer entries
							break
						}
						if prev == nil {
							gp._defer = d.link
						} else {
							prev.link = d.link
						}
						newd := d.link
						freedefer(d)
						d = newd
					} else {
						prev = d
						d = d.link
					}
				}
			}

			gp._panic = p.link
			// Aborted panics are marked but remain on the g.panic list.
			// Remove them from the list.
			for gp._panic != nil && gp._panic.aborted {
				gp._panic = gp._panic.link
			}
			if gp._panic == nil { // must be done with signal
				gp.sig = 0
			}
			// Pass information about recovering frame to recovery.
			gp.sigcode0 = uintptr(sp)
			gp.sigcode1 = pc
			mcall(recovery)
			throw("recovery failed") // mcall should not return
		}
```

这段代码有点长，主要就是分为两部分：

第一部分主要是这个判断`if gp._panic != nil && gp._panic.goexit && gp._panic.aborted { ... }`，正常recover是会绕过`Goexit`的，所以为了解决这个，添加了这个判断，这样就可以保证`Goexit`也会被`recover`住，这里是通过从`runtime._panic`中取出了程序计数器`pc`和栈指针`sp`并且调用`runtime.recovery`函数触发`goroutine`的调度，调度之前会准备好 `sp`、`pc` 以及函数的返回值。

第二部分主要是做`panic`的`recover`，这也与上面的流程基本差不多，他是从`runtime._defer`中取出了程序计数器`pc`和`栈指针sp`并调用`recovery`函数触发`Goroutine`，跳转到`recovery`函数是通过`runtime.call`进行的，我们看一下其源码(src/runtime/asm_amd64.s 289行)：

```go
// func mcall(fn func(*g))
// Switch to m->g0's stack, call fn(g).
// Fn must never return. It should gogo(&g->sched)
// to keep running g.
TEXT runtime·mcall(SB), NOSPLIT, $0-8
	MOVQ	fn+0(FP), DI

	get_tls(CX)
	MOVQ	g(CX), AX	// save state in g->sched
	MOVQ	0(SP), BX	// caller's PC
	MOVQ	BX, (g_sched+gobuf_pc)(AX)
	LEAQ	fn+0(FP), BX	// caller's SP
	MOVQ	BX, (g_sched+gobuf_sp)(AX)
	MOVQ	AX, (g_sched+gobuf_g)(AX)
	MOVQ	BP, (g_sched+gobuf_bp)(AX)

	// switch to m->g0 & its stack, call fn
	MOVQ	g(CX), BX
	MOVQ	g_m(BX), BX
	MOVQ	m_g0(BX), SI
	CMPQ	SI, AX	// if g == m->g0 call badmcall
	JNE	3(PC)
	MOVQ	$runtime·badmcall(SB), AX
	JMP	AX
	MOVQ	SI, g(CX)	// g = m->g0
	MOVQ	(g_sched+gobuf_sp)(SI), SP	// sp = m->g0->sched.sp
	PUSHQ	AX
	MOVQ	DI, DX
	MOVQ	0(DI), DI
	CALL	DI
	POPQ	AX
	MOVQ	$runtime·badmcall2(SB), AX
	JMP	AX
	RET
```

因为`go`语言中的`runtime`环境是有自己的堆栈和`goroutine`，`recovery`函数也是在`runtime`环境执行的，所以要调度到`m->g0`来执行`recovery`函数，我们在看一下`recovery`函数：

```go
// Unwind the stack after a deferred function calls recover
// after a panic. Then arrange to continue running as though
// the caller of the deferred function returned normally.
func recovery(gp *g) {
	// Info about defer passed in G struct.
	sp := gp.sigcode0
	pc := gp.sigcode1

	// d's arguments need to be in the stack.
	if sp != 0 && (sp < gp.stack.lo || gp.stack.hi < sp) {
		print("recover: ", hex(sp), " not in [", hex(gp.stack.lo), ", ", hex(gp.stack.hi), "]\n")
		throw("bad recovery")
	}

	// Make the deferproc for this d return again,
	// this time returning 1. The calling function will
	// jump to the standard return epilogue.
	gp.sched.sp = sp
	gp.sched.pc = pc
	gp.sched.lr = 0
	gp.sched.ret = 1
	gogo(&gp.sched)
}
```

在`recovery` 函数中，利用 `g` 中的两个状态码回溯栈指针 sp 并恢复程序计数器 pc 到调度器中，并调用 `gogo` 重新调度 `g` ，将 `g` 恢复到调用 `recover` 函数的位置， goroutine 继续执行，`recovery`在调度过程中会将函数的返回值设置为1。这个有什么作用呢？ 在`deferproc`函数中找到了答案：

```go
//go:nosplit
func deferproc(siz int32, fn *funcval) { // arguments of fn follow fn
  ............ 省略
// deferproc returns 0 normally.
	// a deferred func that stops a panic
	// makes the deferproc return 1.
	// the code the compiler generates always
	// checks the return value and jumps to the
	// end of the function if deferproc returns != 0.
	return0()
	// No code can go here - the C return register has
	// been set and must not be clobbered.
}
```

当延迟函数中`recover`了一个`panic`时，就会返回1，当 `runtime.deferproc` 函数的返回值是 1 时，编译器生成的代码会直接跳转到调用方函数返回之前并执行 `runtime.deferreturn`，跳转到`runtime.deferturn`函数之后，程序就已经从`panic`恢复了正常的逻辑。



- 第五部分，如果没有遇到`runtime.gorecover`就会依次遍历所有的`runtime._defer`，在最后调用`fatalpanic`中止程序，并打印`panic`参数返回错误码2。

```go
// fatalpanic implements an unrecoverable panic. It is like fatalthrow, except
// that if msgs != nil, fatalpanic also prints panic messages and decrements
// runningPanicDefers once main is blocked from exiting.
//
//go:nosplit
func fatalpanic(msgs *_panic) {
	pc := getcallerpc()
	sp := getcallersp()
	gp := getg()
	var docrash bool
	// Switch to the system stack to avoid any stack growth, which
	// may make things worse if the runtime is in a bad state.
	systemstack(func() {
		if startpanic_m() && msgs != nil {
			// There were panic messages and startpanic_m
			// says it's okay to try to print them.

			// startpanic_m set panicking, which will
			// block main from exiting, so now OK to
			// decrement runningPanicDefers.
			atomic.Xadd(&runningPanicDefers, -1)

			printpanics(msgs)
		}

		docrash = dopanic_m(gp, pc, sp)
	})

	if docrash {
		// By crashing outside the above systemstack call, debuggers
		// will not be confused when generating a backtrace.
		// Function crash is marked nosplit to avoid stack growth.
		crash()
	}

	systemstack(func() {
		exit(2)
	})

	*(*int)(nil) = 0 // not reached
}
```

在这里`runtime.fatalpanic`实现了无法被恢复的程序崩溃，它在中止程序之前会通过 `runtime.printpanics` 打印出全部的 `panic` 消息以及调用时传入的参数。

好啦，至此整个`gopanic`方法就全部看完了，接下来我们再来看一看`gorecover`方法。



### gorecover

这个函数就简单很多了，代码量比较少，先看一下代码吧：

```go
// The implementation of the predeclared function recover.
// Cannot split the stack because it needs to reliably
// find the stack segment of its caller.
//
// TODO(rsc): Once we commit to CopyStackAlways,
// this doesn't need to be nosplit.
//go:nosplit
func gorecover(argp uintptr) interface{} {
	// Must be in a function running as part of a deferred call during the panic.
	// Must be called from the topmost function of the call
	// (the function used in the defer statement).
	// p.argp is the argument pointer of that topmost deferred function call.
	// Compare against argp reported by caller.
	// If they match, the caller is the one who can recover.
	gp := getg()
	p := gp._panic
	if p != nil && !p.goexit && !p.recovered && argp == uintptr(p.argp) {
		p.recovered = true
		return p.arg
	}
	return nil
}
```

首先获取当前所在的`Goroutine`，如果当前`Goroutine`没有调用`panic`，那么该函数会直接返回`nil`，是否能`recover`住该`panic`的判断条件必须四个都吻合，`p.Goexit`判断当前是否是`goexit`触发的，如果是则无法`revocer`住，上面讲过会在`gopanic`中执行进行`recover`。`argp`是是最顶层延迟函数调用的实参指针，与调用者的`argp`进行比较，如果匹配说明调用者是可以`recover`，直接将`recovered`字段设置为`true`就可以了。这里主要的作用就是判断当前`panic`是否可以`recover`，具体的恢复逻辑还是由`gopanic`函数负责的。



## 流程总结

上面看了一篇源码，肯定也是一脸懵逼吧～。这正常，毕竟文字诉说，只能到这个程度了，还是要自己结合带去去看，这里只是起一个辅助作用，最后做一个流程总结吧。

- 在程序执行过程中如果遇到`panic`，那么会调用`runtime.gopanic`，然后取当前`Goroutine`的`defer`链表依次执行。
- 在调用`defer`函数是如果有`recover`就会调用`runtime.gorecover`，在`gorecover`中会把`runtime._panic`中的`recoved`标记为`true`，这里只是标记的作用，恢复逻辑仍在`runtime.panic`中。
- 在`gopanic`中会执行`defer`内联优化、程序恢复逻辑。在程序恢复逻辑中，会进行判断，如果是触发是`runtime.Goexit`，也会进行`recovery`。`panic`也会进行`recovery`，主要逻辑是`runtime.gopanic`会从`runtime._defer`结构体中取出程序计数器`pc`和栈指针`sp`并调用`runtime.recovery`函数恢复程序。`runtime.recvoery`函数中会根据传入的 `pc` 和 `sp` 在`gogo`中跳转回`runtime.deferproc`，如果返回值为1，就会调用`runtime.deferreturn`恢复正常流程。
- 在`gopanic`执行完所有的`_defer`并且也没有遇到`recover`，那么就会执行`runtime.fatalpanic`终止程序，并返回错误码2.

这就是这个逻辑流程，累死我了。。。。



## 小彩蛋

结尾给大家发一个小福利，哈哈，这个福利就是如果避免出现`panic`，要注意这些：

- 数组/切片下标越界，对于`go`这种静态语言来说，下标越界是致命问题。
- 不要访问未初始化的指针或`nil`指针
- 不要往已经`close`的`chan`里发送数据
- `map`不是线程安全的，不要并发读写`map`

这几个是比较典型的，还有很多会发生`panic`的地方，交给你们自行学习吧～。

## 总结

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

