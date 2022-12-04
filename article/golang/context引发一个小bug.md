## 背景

> 哈喽，everybody，小松子，再次回归，最近工作比较忙，好久都没有更新文章了，接下来会进行补更。今天与大家分享一个日常开发比较容易错误的点，那就是`contxt`误用导致的`bug`，我自己就因为误用导致异步更新缓存都失败了，究竟是因为什么呢？看这样一个例子，光看代码，你能看出来有什么`bug`吗？

```go
func AsyncAdd(run func() error)  {
	//TODO: 扔进异步协程池
	go run()
}

func GetInstance(ctx context.Context,id uint64) (string, error) {
	data,err := GetFromRedis(ctx,id)
	if err != nil && err != redis.Nil{
		return "", err
	}
	// 没有找到数据
	if err == redis.Nil {
		data,err = GetFromDB(ctx,id)
		if err != nil{
			return "", err
		}
		AsyncAdd(func() error{
			return UpdateCache(ctx,id,data)
		})
	}
	return data,nil
}

func GetFromRedis(ctx context.Context,id uint64) (string,error) {
	// TODO: 从redis获取信息
	return "",nil
}

func GetFromDB(ctx context.Context,id uint64) (string,error) {
	// TODO: 从DB中获取信息
	return "",nil
}

func UpdateCache(ctx context.Context,id interface{},data string) error {
	// TODO：更新缓存信息
	return nil
}

func main()  {
	ctx,cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	_,err := GetInstance(ctx,2021)
	if err != nil{
		return
	}
}
```



## 分析

我们先简单分析一下，这一段代码要干什么？其实很简单，我们想要获取一段信息，首先会从缓存中获取，如果缓存中获取不到，我们就从`DB`中获取，从DB中获取到信息后，在协程池中放入更新缓存的方法，异步去更新缓存。整个设计是不是很完美，但是在实际工作中，异步更新缓存就没有成功过？

导致失败的原因就在这一段代码：

```go
	AsyncAdd(func() error{
			return UpdateCache(ctx,id,data)
		})
```

错误的原因只有一个，就是这个`ctx`，如果改成这样，就啥事没有了。

```go
AsyncAdd(func() error{
			ctxAsync,cancel := context.WithTimeout(context.Background(),3 * time.Second)
			defer cancel()
			return UpdateCache(ctxAsync,id,data)
		})
```

看到这个，想必大家就已经知道为什么吧？

在这个`ctx`树中，根结点发生了`cancel()`，会将信号即时同步给下层，因为异步任务的`ctx`也在这棵树的节点上，所以当`main goroutine`取消了`ctx`时，异步任务也被取消了，导致了缓存更新一直失败。

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-03-20%20%E4%B8%8B%E5%8D%885.28.33.png)

因为我之前写过一篇关于[详解Context包，看这一篇就够了！！！](https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg)的文章，就不在这里细说其原理了，想知道其内部是怎么实现的，看以前[这篇](https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg)文章就可以了。在这里在与大家分享一下`context`的使用原则，避免踩坑。

- context.Background 只应用在最高等级，作为所有派生 context 的根。
- context 取消是建议性的，这些函数可能需要一些时间来清理和退出。
- 不要把`Context`放在结构体中，要以参数的方式传递。
- 以`Context`作为参数的函数方法，应该把`Context`作为第一个参数，放在第一位。
- 给一个函数方法传递Context的时候，不要传递nil，如果不知道传递什么，就使用context.TODO
- Context的Value相关方法应该传递必须的数据，不要什么数据都使用这个传递。context.Value 应该很少使用，它不应该被用来传递可选参数。这使得 API 隐式的并且可以引起错误。取而代之的是，这些值应该作为参数传递。
- Context是线程安全的，可以放心的在多个goroutine中传递。同一个Context可以传给使用其的多个goroutine，且Context可被多个goroutine同时安全访问。
- Context 结构没有取消方法，因为只有派生 context 的函数才应该取消 context。

Go 语言中的 `context.Context` 的主要作用还是在多个 Goroutine 组成的树中同步取消信号以减少对资源的消耗和占用，虽然它也有传值的功能，但是这个功能我们还是很少用到。在真正使用传值的功能时我们也应该非常谨慎，使用 `context.Context` 进行传递参数请求的所有参数一种非常差的设计，比较常见的使用场景是传递请求对应用户的认证令牌以及用于进行分布式追踪的请求 ID。



## 总结

写这篇文章的目的，就是把我日常写的`bug`分享出来，防止后人踩坑。已经踩过的坑就不要再踩了，把找`bug`的时间节省出来，多学点其他知识，他不香嘛～。

**好啦，这篇文章就到这里啦，素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：加我vx拉你入群，或者公众号获取入群二维码**

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**翻译了一份Machinery中文文档，会定期进行维护，有需要的小伙伴们后台回复[machinery]即可获取。**

**我是asong，一名普普通通的程序猿，让我们一起慢慢变强吧。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E7%99%BD%E8%89%B2%E7%89%88.png)

推荐往期文章：

- [Go看源码必会知识之unsafe包](https://mp.weixin.qq.com/s/nPWvqaQiQ6Z0TaPoqg3t2Q)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/mzSCWI8C_ByIPbb07XYFTQ)
- [详解并发编程基础之原子操作(atomic包)](https://mp.weixin.qq.com/s/PQ06eL8kMWoGXodpnyjNcA)
- [详解defer实现机制](https://mp.weixin.qq.com/s/FUmoBB8OHNSfy7STR0GsWw)
- [空结构体引发的大型打脸现场](https://mp.weixin.qq.com/s/dNeCIwmPei2jEWGF6AuWQw)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/wURQFRt2ISz66icW7jbHFw)
- [十张动图带你搞懂排序算法(附go实现代码)](https://mp.weixin.qq.com/s/rZBsoKuS-ORvV3kML39jKw)
- [go参数传递类型](https://mp.weixin.qq.com/s/JHbFh2GhoKewlemq7iI59Q)
- [手把手教姐姐写消息队列](https://mp.weixin.qq.com/s/0MykGst1e2pgnXXUjojvhQ)
- [常见面试题之缓存雪崩、缓存穿透、缓存击穿](https://mp.weixin.qq.com/s?__biz=MzIzMDU0MTA3Nw==&mid=2247483988&idx=1&sn=3bd52650907867d65f1c4d5c3cff8f13&chksm=e8b0902edfc71938f7d7a29246d7278ac48e6c104ba27c684e12e840892252b0823de94b94c1&token=1558933779&lang=zh_CN#rd)
- [详解Context包，看这一篇就够了！！！](https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg)
- [高并发系统的限流策略：漏桶和令牌桶(附源码剖析)](https://mp.weixin.qq.com/s/fURwiSTeEE_Wvc95Q_fHnA)
- [面试官：go中for-range使用过吗？这几个问题你能解释一下原因吗](https://mp.weixin.qq.com/s/G7z80u83LTgLyfHgzgrd9g)