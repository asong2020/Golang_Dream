## 前言

> 哈喽，大家好，我是asong，这是我的第16篇原创文章，感谢各位的关注。今天给大家分享设计模式之单例模式，并使用go语言实现。熟悉java的同学对单例模式一定不陌生，单例模式，是一种很常见的软件设计模式，在他的核心结构中只包含一个被称为单例的特殊类。通过单例模式可以保证系统中一个类只有一个实例且该实例易于外界访问，从而方便对实例个数的控制并节约系统资源。下面我们就一起来看一看怎么使用go实现单例模式，这里有一个小坑，一定要注意一下，结尾告诉你哦～～～



## 什么是单例模式

单例模式确保某一个类只有一个实例。为什么要确保一个类只有一个实例？有什么时候才需要用到单例模式呢？听起来一个类只有一个实例好像没什么用呢！ 那我们来举个例子。比如我们的APP中有一个类用来保存运行时全局的一些状态信息，如果这个类实现不是单例的，那么App里面的组件能够随意的生成多个类用来保存自己的状态，等于大家各玩各的，那这个全局的状态信息就成了笑话了。而如果把这个类实现成单例的，那么不管App的哪个组件获取到的都是同一个对象（比如Application类，除了多进程的情况下）。

<img src="https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/singleton.png" style="zoom:50%;" />



## 饿汉模式

这里我们使用三种方式实现饿汉模式。先说一下什么是懒汉模式吧，从懒汉这两个字，我们就能知道，这个人很懒，所以他不可能在未使用实例时就创建了对象，他肯定会在使用时才会创建实例，这个好处的就在于，只有在使用的时候才会创建该实例。下面我们一起来看看他的实现：



- 不加锁

```go
package one

type singleton struct {

}

var  instance *singleton
func GetInstance() *singleton {
	if instance == nil{
		instance = new(singleton)
	}
	return instance
}
```

这种方法是会存在线程安全问题的，在高并发的时候会有多个线程同时掉这个方法，那么都会检测`instance`为nil，这样就会导致创建多个对象，所以这种方法是不推荐的，我再来看第二种写法。



- 整个方法加锁

```go
type singleton struct {

}

var instance *singleton
var lock sync.Mutex

func GetInstance() *singleton {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil{
		instance = new(singleton)
	}
	return instance
}
```

这里对整个方法进行了加锁，这种可以解决并发安全的问题，但是效率就会降下来，每一个对象创建时都是进行加锁解锁，这样就拖慢了速度，所以不推荐这种写法。

- 创建方法时进行锁定

```go
type singleton struct {

}

var instance *singleton
var lock sync.Mutex

func GetInstance() *singleton {
	if instance == nil{
		lock.Lock()
		instance = new(singleton)
		lock.Unlock()
	}
	return instance
}
```

这种方法也是线程不安全的，虽然我们加了锁，多个线程同样会导致创建多个实例，所以这种方式也不是推荐的。所以就有了下面的双重检索机制

- 双重检锁

```go
type singleton struct {
	
}

var instance *singleton
var lock sync.Mutex

func GetInstance() *singleton {
	if instance == nil{
		lock.Lock()
		if instance == nil{
			instance = new(singleton)
		}
		lock.Unlock()
	}
	return instance
}
```

这里在上面的代码做了改进，只有当对象未初始化的时候，才会有加锁和减锁的操作。但是又出现了另一个问题：每一次访问都要检查两次，为了解决这个问题，我们可以使用golang标准包中的方法进行原子性操作。

- 原子操作实现

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

这里使用了`sync.Once`的`Do`方法可以实现在程序运行过程中只运行一次其中的回调，这样就可以只创建了一个对象，这种方法是推荐的～～～。



## 饿汉模式

有懒汉模式，当然还要有饿汉模式啦，看了懒汉的模式，饿汉模式我们很好解释了，因为他饿呀，所以很着急的就创建了实例，不用等到使用时才创建，这样我们每次调用获取接口将不会重新创建新的对象，而是直接返回之前创建的对象。比较适用于：如果某个单例使用的次数少，并且创建单例消息的资源比较多，那么就需要实现单例的按需创建，这个时候懒汉模式就是一个不错的选择。不过也有缺点，饿汉模式将在包加载的时候就会创建单例对象，当程序中用不到该对象时，浪费了一部分空间，但是相对于懒汉模式，不需要进行了加锁操作，会更安全，但是会减慢启动速度。

下面我们一起来看看go实现饿汉模式：

```go
type singleton struct {

}

var instance = new(singleton)

func GetInstance()  *singleton{
	return instance
}

或者
type singleton struct {

}

var instance *singleton

func init()  {
	instance = new(singleton)
}

func GetInstance()  *singleton{
	return instance
}

```

这两种方法都可以，第一种我们采用创建一个全局变量的方式来实现，第二种我们使用`init`包加载的时候创建实例，这里两个都可以，不过根据golang的执行顺序，全局变量的初始化函数会比包的`init`函数先执行，没有特别的差距。



## 小坑

还记得我开头说的一句话，`go`语言中使用单例模式有一个小坑，如果不注意，就会导致我们的单例模式没有用，可以观察一下我写的代码，除了`GetInstance`方法外其他都使用的小写字母开头，知道这是为什么吗？

golang中根据首字母的大小写来确定可以访问的权限。无论是方法名、常量、变量名还是结构体的名称，如果首字母大写，则可以被其他的包访问；如果首字母小写，则只能在本包中使用。可以简单的理解成，首字母大写是公有的，首字母小写是私有的。这里`type singleton struct {`我们如果使用大写，那么我们写的这些方法就没有意义了，其他包可以通过`s  := &singleton{}`创建多个实例，单例模式就显得很没有意义了，所以这里一定要注意一下哦～～～



## 总结

>  这一篇就到此结束了，这里讲解了23种模式中最简单的单例模式，虽然他很简单，但是越简单的越容易犯错的呦，所以一定要细心对待每一件事情的呦～～
>
> 好啦，这一篇就到此结束了，我的代码已上传github：https://github.com/asong2020/Golang_Dream/tree/master/code_demo/singleton
>
> 欢迎star



**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**我是asong，一名普普通通的程序猿，让我一起慢慢变强吧。我自己建了一个`golang`交流群，有需要的小伙伴加我`vx`,我拉你入群。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [手把手教姐姐写消息队列](https://mp.weixin.qq.com/s/0MykGst1e2pgnXXUjojvhQ)

- [详解Context包，看这一篇就够了！！！](https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg)
- [go-ElasticSearch入门看这一篇就够了(一)](https://mp.weixin.qq.com/s/mV2hnfctQuRLRKpPPT9XRw)
- [面试官：go中for-range使用过吗？这几个问题你能解释一下原因吗](https://mp.weixin.qq.com/s/G7z80u83LTgLyfHgzgrd9g)
- [学会wire依赖注入、cron定时任务其实就这么简单！](https://mp.weixin.qq.com/s/qmbCmwZGmqKIZDlNs_a3Vw)
- [听说你还不会jwt和swagger-饭我都不吃了带着实践项目我就来了](https://mp.weixin.qq.com/s/z-PGZE84STccvfkf8ehTgA)
- [掌握这些Go语言特性，你的水平将提高N个档次(二)](https://mp.weixin.qq.com/s/7yyo83SzgQbEB7QWGY7k-w)

- [go实现多人聊天室，在这里你想聊什么都可以的啦！！！](https://mp.weixin.qq.com/s/H7F85CncQNdnPsjvGiemtg)
- [grpc实践-学会grpc就是这么简单](https://mp.weixin.qq.com/s/mOkihZEO7uwEAnnRKGdkLA)
- [go标准库rpc实践](https://mp.weixin.qq.com/s/d0xKVe_Cq1WsUGZxIlU8mw)
- [2020最新Gin框架中文文档 asong又捡起来了英语，用心翻译](https://mp.weixin.qq.com/s/vx8A6EEO2mgEMteUZNzkDg)
- [基于gin的几种热加载方式](https://mp.weixin.qq.com/s/CZvjXp3dimU-2hZlvsLfsw)
- [boss: 这小子还不会使用validator库进行数据校验，开了～～～](https://mp.weixin.qq.com/s?__biz=MzIzMDU0MTA3Nw==&mid=2247483829&idx=1&sn=d7cf4f46ea038a68e74a4bf00bbf64a9&scene=19&token=1606435091&lang=zh_CN#wechat_redirect)