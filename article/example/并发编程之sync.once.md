> 哈喽，大家好，我是`asong`，这是我并发编程系列的第二篇文章. 上一篇我们一起分析了[`atomic`包](https://mp.weixin.qq.com/s/PQ06eL8kMWoGXodpnyjNcA)，今天我们一起来看一看`sync/once`的使用与实现.


## 什么是`sync.once`

Go语言标准库中的`sync.Once`可以保证`go`程序在运行期间的某段代码只会执行一次，作用与`init`类似，但是也有所不同：
-  `init`函数是在文件包首次被加载的时候执行，且只执行一次。
- `sync.Once`是在代码运行中需要的时候执行，且只执行一次。

还记得我之前写的一篇关于[`go`单例模式](https://mp.weixin.qq.com/s/7rgs9J5jlnMn-S7_GE_IzA),懒汉模式的一种实现就可以使用`sync.Once`，他可以解决双重检锁带来的每一次访问都要检查两次的问题，因为`sync.once`的内部实现可以完全解决这个问题(后面分析完源码就知道原因了)，下面我们来看一看这种懒汉模式怎么写：

```go
type singleton struct {
	
}

var instance *singleton
var once sync.Once
func GetInstance() *singleton {
	once.Do(func() {
		instance = new(singleton)
	})
	return instance
}
```
实现还是比较简单，就不细说了。


## 源码解析

`sync.Once`的源码还是很少的，首先我们看一下他的结构：

```go
// Once is an object that will perform exactly one action.
type Once struct {
	// done indicates whether the action has been performed.
	// It is first in the struct because it is used in the hot path.
	// The hot path is inlined at every call site.
	// Placing done first allows more compact instructions on some architectures (amd64/x86),
	// and fewer instructions (to calculate offset) on other architectures.
	done uint32
	m    Mutex
}
```
只有两个字段，字段`done`用来标识代码块是否执行过，字段`m`是一个互斥锁。

接下来我们一起来看一下代码实现：

```go
func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
```

这里把注释都省略了，反正都是英文，接下来咱用中文解释哈。`sync.Once`结构对外只提供了一个`Do()`方法，该方法的参数是一个入参为空的函数，这个函数也就是我们想要执行一次的代码块。接下来我们看一下代码流程：

- 首先原子性的读取`done`字段的值是否改变，没有改变则执行`doSlow()`方法.

- 一进入`doslow()`方法就开始执行加锁操作，这样在并发情况下可以保证只有一个线程会执行，在判断一次当前`done`字段是否发生改变(这里肯定有朋友会感到疑惑，为什么这里还要在判断一次`flag`？这里目的其实就是保证并发的情况下，代码块也只会执行一次，毕竟加锁是在`doslow()`方法内，不加这个判断的在并发情况下就会出现其他`goroutine`也能执行`f()`)，如果未发生改变，则开始执行代码块，代码块运行结束后会对`done`字段做原子操作，标识该代码块已经被执行过了.


## 优化sync.Once

如果让你自己写一个这样的库，你会考虑的这样全面吗？相信聪明的你们也一定会写出这样一段代码。如果要是我来写，上面的代码可能都一样，但是在`if o.done == 0 `这里我可能会采用`CAS`原子操作来代替这个判断，如下：

```go
type MyOnce struct {
	flag uint32
	lock sync.Mutex
}

func (m *MyOnce)Do(f func())  {
	if atomic.LoadUint32(&m.flag) == 0{
		m.lock.Lock()
		defer m.lock.Unlock()
		if atomic.CompareAndSwapUint32(&m.flag,0,1){
			f()
		}
	}
}

func testDo()  {
	mOnce := MyOnce{}
	for i := 0;i<10;i++{
		go func() {
			mOnce.Do(func() {
				fmt.Println("test my once only run once")
			})
		}()
	}
}

func main()  {
	testDo()
	time.Sleep(10 * time.Second)
}
// 运行结果：
test my once only run once
```

我就说原子操作是并发编程的基础吧，你看没有错吧～。


## 小试牛刀

上面我们也看了源码的实现，现在我们来看三道题，你认为他们的答案是多少？

### 问题一

`sync.Once()`方法中传入的函数发生了`panic`，重复传入还会执行吗？

```go
func panicDo()  {
	once := &sync.Once{}
	defer func() {
		if err := recover();err != nil{
			once.Do(func() {
				fmt.Println("run in recover")
			})
		}
	}()
	once.Do(func() {
		panic("panic i=0")
	})

}
```

### 问题二

`sync.Once()`方法传入的函数中再次调用`sync.Once()`方法会有什么问题吗？

```go
func nestedDo()  {
	once := &sync.Once{}
	once.Do(func() {
		once.Do(func() {
			fmt.Println("test nestedDo")
		})
	})
}
```

### 问题三

改成这样呢？

```go
func nestedDo()  {
	once1 := &sync.Once{}
	once2 := &sync.Once{}
	once1.Do(func() {
		once2.Do(func() {
			fmt.Println("test nestedDo")
		})
	})
}
```

## 总结

在本文的最把上面三道题的答案公布一下吧：

- 问题一：不会打印任何东西，`sync.Once.Do` 方法中传入的函数只会被执行一次，哪怕函数中发生了 `panic`；

- 问题二：发生死锁，根据源码实现我们可以知道在第二个`do`方法会一直等`doshow()`中锁的释放导致发生了死锁;

- 问题三：打印`test nestedDo`，once1，once2是两个对象，互不影响。所以`sync.Once`是使方法只执行一次对象的实现。

你们都做对了吗？

**好啦，这篇文章就到这里啦，素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：加我vx拉你入群，或者公众号获取入群二维码**

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**翻译了一份Machinery中文文档，会定期进行维护，有需要的小伙伴们后台回复[machinery]即可获取。**

**我是asong，一名普普通通的程序猿，让我们一起慢慢变强吧。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [machinery-go异步任务队列](https://mp.weixin.qq.com/s/4QG69Qh1q7_i0lJdxKXWyg)
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
