## 前言

> 哈喽，大家好，我是`asong`。最近一个群里看到一个有趣的八股文，问题是：使用`context`携带的`value`是线程安全的吗？这道题其实就是考察面试者对`context`实现原理的理解，如果不知道`context`的实现原理，很容易答错这道题，所以本文我们就借着这道题，再重新理解一遍`context`携带`value`的实现原理。



## `context`携带`value`是线程安全的吗？

先说答案，`context`本身就是线程安全的，所以`context`携带`value`也是线程安全的，写个简单例子验证一下：

```go
func main()  {
	ctx := context.WithValue(context.Background(), "asong", "test01")
	go func() {
		for {
			_ = context.WithValue(ctx, "asong", "test02")
		}
	}()
	go func() {
		for {
			_ = context.WithValue(ctx, "asong", "test03")
		}
	}()
	go func() {
		for {
			fmt.Println(ctx.Value("asong"))
		}
	}()
	go func() {
		for {
			fmt.Println(ctx.Value("asong"))
		}
	}()
	time.Sleep(10 * time.Second)
}
```

程序正常运行，没有任何问题，接下来我们就来看一下为什么`context`是线程安全的！！！



## 为什么线程安全？

`context`包提供两种创建根`context`的方式：

- `context.Backgroud()`
- `context.TODO()`

又提供了四个函数基于父`Context`衍生，其中使用`WithValue`函数来衍生`context`并携带数据，每次调用`WithValue`函数都会基于当前`context`衍生一个新的子`context`，`WithValue`内部主要就是调用`valueCtx`类：

```go
func WithValue(parent Context, key, val interface{}) Context {
 if parent == nil {
  panic("cannot create context from nil parent")
 }
 if key == nil {
  panic("nil key")
 }
 if !reflectlite.TypeOf(key).Comparable() {
  panic("key is not comparable")
 }
 return &valueCtx{parent, key, val}
}
```

`valueCtx`结构如下：

```go
type valueCtx struct {
 Context
 key, val interface{}
}
```

`valueCtx`继承父`Context`，这种是采用匿名接口的继承实现方式，`key,val`用来存储携带的键值对。

通过上面的代码分析，可以看到添加键值对不是在原`context`结构体上直接添加，而是以此`context`作为父节点，重新创建一个新的`valueCtx`子节点，将键值对添加在子节点上，由此形成一条`context`链。

获取键值过程也是层层向上调用直到最终的根节点，中间要是找到了`key`就会返回，否会就会找到最终的`emptyCtx`返回`nil`。

画个图表示一下：

![image-20220207214507921](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/image-20220207214507921.png)

总结：`context`添加的键值对一个链式的，会不断衍生新的`context`，所以`context`本身是不可变的，因此是线程安全的。



## 总结

本文主要是想带大家回顾一下`context`的实现原理，面试中面试官都喜欢隐晦提出问题，所以这就需要我们有很扎实的基本功，一不小心就会掉入面试官的陷阱，要处处小心哦～

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)

