## 前言

> 哈喽，大家好，我是`asong`。前几天一个读者问我如何使用`Go`语言实现可重入锁，突然想到`Go`语言中好像没有这个概念，平常在业务开发中也没有要用到可重入锁的概念，一时懵住了。之前在写`java`的时候，就会使用到可重入锁，然而写了这么久的`Go`，却没有使用过，这是怎么回事呢？这一篇文章就带你来解密～



## 什么是可重入锁

之前写过`java`的同学对这个概念应该了如指掌，可重入锁又称为递归锁，是指在同一个线程在外层方法获取锁的时候，在进入该线程的内层方法时会自动获取锁，不会因为之前已经获取过还没释放而阻塞。[美团技术团队](https://tech.meituan.com/2018/11/15/java-lock.html)的一篇关于锁的文章当中针对可重入锁进行了举例：

假设现在有多个村民在水井排队打水，有管理员正在看管这口水井，村民在打水时，管理员允许锁和同一个人的多个水桶绑定，这个人用多个水桶打水时，第一个水桶和锁绑定并打完水之后，第二个水桶也可以直接和锁绑定并开始打水，所有的水桶都打完水之后打水人才会将锁还给管理员。这个人的所有打水流程都能够成功执行，后续等待的人也能够打到水。这就是可重入锁。

下图摘自美团技术团队分享的文章：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-05-29%20%E4%B8%8B%E5%8D%884.30.49.png)

如果是非可重入锁，，此时管理员只允许锁和同一个人的一个水桶绑定。第一个水桶和锁绑定打完水之后并不会释放锁，导致第二个水桶不能和锁绑定也无法打水。当前线程出现死锁，整个等待队列中的所有线程都无法被唤醒。

下图依旧摘自美团技术团队分享的文章：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-05-29%20%E4%B8%8B%E5%8D%884.32.01.png)



## 用`Go`实现可重入锁

既然我们想自己实现一个可重入锁，那我们就要了解`java`中可重入锁是如何实现的，查看了`ReentrantLock`的源码，大致实现思路如下：

`ReentrantLock`继承了父类`AQS`，其父类`AQS`中维护了一个同步状态`status`来计数重入次数，`status`初始值为`0`，当线程尝试获取锁时，可重入锁先尝试获取并更新`status`值，如果`status == 0`表示没有其他线程在执行同步代码，则把`status`置为`1`，当前线程开始执行。如果`status != 0`，则判断当前线程是否是获取到这个锁的线程，如果是的话执行`status+1`，且当前线程可以再次获取锁。释放锁时，可重入锁同样先获取当前`status`的值，在当前线程是持有锁的线程的前提下。如果`status-1 == 0`，则表示当前线程所有重复获取锁的操作都已经执行完毕，然后该线程才会真正释放锁。

总结一下实现一个可重入锁需要这两点：

- 记住持有锁的线程
- 统计重入的次数

统计重入的次数很容易实现，接下来我们考虑一下怎么实现记住持有锁的线程？

我们都知道`Go`语言最大的特色就是从语言层面支持并发，`Goroutine`是`Go`中最基本的执行单元，每一个`Go`程序至少有一个`Goroutine`，主程序也是一个`Goroutine`，称为主`Goroutine`，当程序启动时，他会自动创建。每个`Goroutine`也是有自己唯一的编号，这个编号只有在`panic`场景下才会看到，`Go语言`却刻意没有提供获取该编号的接口，官方给出的原因是为了避免滥用。但是我们还是通过一些特殊手段来获取`Goroutine ID`的，可以使用`runtime.Stack`函数输出当前栈帧信息，然后解析字符串获取`Goroutine ID`，具体代码可以参考开源项目 - [goid](https://github.com/petermattis/goid/blob/master/goid.go)。

因为`go`语言中的`Goroutine`有`Goroutine ID`，那么我们就可以通过这个来记住当前的线程，通过这个来判断是否持有锁，就可以了，因此我们可以定义如下结构体：

```go
type ReentrantLock struct {
	lock *sync.Mutex
	cond *sync.Cond
	recursion int32
	host     int64
}
```

其实就是包装了`Mutex`锁，使用`host`字段记录当前持有锁的`goroutine id`，使用`recursion`字段记录当前`goroutine`的重入次数。这里有一个特别要说明的就是`sync.Cond`，使用`Cond`的目的是，当多个`Goroutine`使用相同的可重入锁时，通过`cond`可以对多个协程进行协调，如果有其他协程正在占用锁，则当前协程进行阻塞，直到其他协程调用释放锁。具体`sync.Cond`的使用大家可以参考我之前的一篇文章：[源码剖析sync.cond(条件变量的实现机制）](https://mp.weixin.qq.com/s/szSxatDakPQMUA8Vm9u3qQ)。

- 构造函数

```go

func NewReentrantLock()  sync.Locker{
	res := &ReentrantLock{
		lock: new(sync.Mutex),
		recursion: 0,
		host: 0,
	}
	res.cond = sync.NewCond(res.lock)
	return res
}
```

- `Lock`

```go
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
```
这里逻辑比较简单，大概解释一下：

首先我们获取当前`Goroutine`的`ID`，然后我们添加互斥锁锁住当前代码块，保证并发安全，如果当前`Goroutine`正在占用锁，则增加`resutsion`的值，记录当前线程加锁的数量，然后返回即可。如果当前`Goroutine`没有占用锁，则判断当前可重入锁是否被其他`Goroutine`占用，如果有其他`Goroutine`正在占用可重入锁，则调用`cond.wait`方法进行阻塞，直到其他协程释放锁。

- `Unlock`

```go
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
```

大概解释如下：

首先我们添加互斥锁锁住当前代码块，保证并发安全，释放可重入锁时，如果非持有锁的`Goroutine`释放锁则会导致程序出现`panic`，这个一般是由于用户用法错误导致的。如果当前`Goroutine`释放了锁，则调用`cond.Signal`唤醒其他协程。



测试例子就不在这里贴了，代码已上传`github`:https://github.com/asong2020/Golang_Dream/tree/master/code_demo/reentrantLock，欢迎star。



## 为什么`Go`语言中没有可重入锁

这问题的答案，我在：https://stackoverflow.com/questions/14670979/recursive-locking-in-go#14671462，这里找到了答案。`Go`语言的发明者认为，如果当你的代码需要重入锁时，那就说明你的代码有问题了，我们正常写代码时，从入口函数开始，执行的层次都是一层层往下的，如果有一个锁需要共享给几个函数，那么就在调用这几个函数的上面，直接加上互斥锁就好了，不需要在每一个函数里面都添加锁，再去释放锁。

举个例子，假设我们现在一段这样的代码：

```go
func F() {
	mu.Lock()
	//... do some stuff ...
	G()
	//... do some more stuff ...
	mu.Unlock()
}

func G() {
	mu.Lock()
	//... do some stuff ...
	mu.Unlock()
}
```

函数`F()`和`G()`使用了相同的互斥锁，并且都在各自函数内部进行了加锁，这要使用就会出现死锁，使用**可重入锁**可以解决这个问题，但是更好的方法是改变我们的代码结构，我们进行分解代码，如下：

```go

func call(){
  F()
  G()
}

func F() {
      mu.Lock()
      ... do some stuff
      mu.Unlock()
}

func g() {
     ... do some stuff ...
}

func G() {
     mu.Lock()
     g()
     mu.Unlock()
}
```

这样不仅避免了死锁，而且还对代码进行了解耦。这样的代码按照作用范围进行了分层，就像金字塔一样，上层调用下层的函数，越往上作用范围越大；各层有自己的锁。

总结：`Go`语言中完全没有必要使用可重入锁，如果我们发现我们的代码要使用到可重入锁了，那一定是我们写的代码有问题了，请检查代码结构，修改他！！！



## 总结

这篇文章我们知道了什么是可重入锁，并用`Go`语言实现了**可重入锁**，大家只需要知道这个概念就好了，实际开发中根本不需要。最后还是建议大家没事多思考一下自己的代码结构，好的代码都是经过深思熟虑的，最后希望大家都能写出漂亮的代码。

**好啦，这篇文章到此结束啦，素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！我是`asong`，我们下期见。**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E6%A0%87%E5%87%86%E8%89%B2%E7%89%88.png)

推荐往期文章：

- [Go看源码必会知识之unsafe包](https://mp.weixin.qq.com/s/nPWvqaQiQ6Z0TaPoqg3t2Q)
- [Go语言中new和make你使用哪个来分配内存？](https://mp.weixin.qq.com/s/XJ9O9O4KS3LbZL0jYnJHPg)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/mzSCWI8C_ByIPbb07XYFTQ)
- [空结构体引发的大型打脸现场](https://mp.weixin.qq.com/s/dNeCIwmPei2jEWGF6AuWQw)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/wURQFRt2ISz66icW7jbHFw)
- [面试官：两个nil比较结果是什么？](https://mp.weixin.qq.com/s/Dt46eoEXXXZc2ymr67_LVQ)
- [面试官：你能用Go写段代码判断当前系统的存储方式吗?](https://mp.weixin.qq.com/s/ffEsTpO-tyNZFR5navAbdA)
- [如何平滑切换线上Elasticsearch索引](https://mp.weixin.qq.com/s/8VQxK_Xh-bkVoOdMZs4Ujw)

