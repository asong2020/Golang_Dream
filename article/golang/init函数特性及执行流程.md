## 前言

> 哈喽，兄弟们，我是`asong`。今天与大家聊一聊`Go`语言中的神奇函数`init`，为什么叫他神奇函数呢？因为该函数可以在所有程序执行开始前被调用，并且每个包下可以有多个`init`函数。这个函数使用起来比较简单，但是你们知道他的执行顺序是怎样的嘛？本文我们就一起来解密。



## `init`函数的特性
先简单介绍一下`init`函数的基本特性：

- `init `函数先于`main`函数自动执行
- 每个包中可以有多个`init`函数，每个包中的源文件中也可以有多个`init`函数
- `init`函数没有输入参数、返回值，也未声明，所以无法引用
- 不同包的`init`函数按照包导入的依赖关系决定执行顺序
- 无论包被导入多少次，`init`函数只会被调用一次，也就是只执行一次

## `init`函数的执行顺序

我在刚学习`init`函数时就对他的执行顺序很好奇，在谷歌上搜了几篇文章，他们都有一样的图：

下图来源于网络：

![截屏2021-06-05 上午9.55.15](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-06-05%20%E4%B8%8A%E5%8D%889.55.15.png)

这张图片很清晰的反应了`init`函数的加载顺序：

- 包加载优先级排在第一位，先层层递归进行包加载
- 每个包中加载顺序为：`const` > `var` > `init`，首先进行初始化的是常量，然后是变量，最后才是`init`函数。针对包级别的变量初始化顺序，[`Go`官方文档](https://golang.org/ref/spec#Package_initialization)给出这样一个例子：

```go
var (
	a = c + b  // == 9
	b = f()    // == 4
	c = f()    // == 5
	d = 3      // == 5 after initialization has finished
)

func f() int {
	d++
	return d
}
```

变量的初始化按出现的顺序从前往后进行，假若某个变量需要依赖其他变量，则被依赖的变量先初始化。所以这个例子中，初始化顺序是 `d` -> `b` -> `c` -> `a`。

上图只是表达了`init`函数大概的加载顺序，有些细节我们还是不知道的，比如：当前包下有多个`init`函数，按照什么顺序执行，当前源文件下有多个`init`函数，这又按照什么顺序执行呢？本来想写个例子挨个验证一下的，后来一看[`Go`官方文档](https://golang.org/ref/spec#Package_initialization)中都有说明，也就没有必要再写一个例子啦，直接说结论吧：

- 如果当前包下有多个`init`函数，首先按照源文件名的字典序从前往后执行。
- 若一个文件中出现多个`init`函数，则按照出现顺序从前往后执行。



前面说的有点乱，对`init`函数的加载顺序做一个小结：

> 从当前包开始，如果当前包包含多个依赖包，则先初始化依赖包，层层递归初始化各个包，在每一个包中，按照源文件的字典序从前往后执行，每一个源文件中，优先初始化常量、变量，最后初始化`init`函数，当出现多个`init`函数时，则按照顺序从前往后依次执行，每一个包完成加载后，递归返回，最后在初始化当前包！





## `init`函数的使用场景

还记得我之前的这篇文章吗：[go解锁设计模式之单例模式](https://mp.weixin.qq.com/s/7rgs9J5jlnMn-S7_GE_IzA)，借用`init`函数的加载机制我们可以实现单例模式中的饿汉模式，具体怎么实现可以参考这篇文章，这里就不在写一遍了。

`init`函数的使用场景还是挺多的，比如进行服务注册、进行数据库或各种中间件的初始化连接等。`Go`的标准库中也有许多地方使用到了`init`函数，比如我们经常使用的`pprof`工具，他就使用到了`init`函数，在`init`函数里面进行路由注册：

```go
//go/1.15.7/libexec/src/cmd/trace/pprof.go
func init() {
	http.HandleFunc("/io", serveSVGProfile(pprofByGoroutine(computePprofIO)))
	http.HandleFunc("/block", serveSVGProfile(pprofByGoroutine(computePprofBlock)))
	http.HandleFunc("/syscall", serveSVGProfile(pprofByGoroutine(computePprofSyscall)))
	http.HandleFunc("/sched", serveSVGProfile(pprofByGoroutine(computePprofSched)))

	http.HandleFunc("/regionio", serveSVGProfile(pprofByRegion(computePprofIO)))
	http.HandleFunc("/regionblock", serveSVGProfile(pprofByRegion(computePprofBlock)))
	http.HandleFunc("/regionsyscall", serveSVGProfile(pprofByRegion(computePprofSyscall)))
	http.HandleFunc("/regionsched", serveSVGProfile(pprofByRegion(computePprofSched)))
}
```

这里就不扩展太多了，更多标准库中的使用方法大家可以自己去探索一下。

在这最后总结一下使用`init`要注意的问题吧：

- 编程时不要依赖`init`的顺序
- 一个源文件下可以有多个`init`函数，代码比较长时可以考虑分多个`init`函数
- 复杂逻辑不建议使用`init`函数，会增加代码的复杂性，可读性也会下降
- 在`init`函数中也可以启动`goroutine`，也就是在初始化的同时启动新的`goroutine`，这并不会影响初始化顺序
- `init`函数不应该依赖任何在`main`函数里创建的变量，因为`init`函数的执行是在`main`函数之前的
- `init`函数在代码中不能被显示的调用，不能被引用（赋值给函数变量），否则会出现编译错误。
- 导入包不要出现循环依赖，这样会导致程序编译失败
- `Go`程序仅仅想要用一个`package`的`init`执行，我们可以这样使用：`import _ "test_xxxx"`，导入包的时候加上下划线就`ok`了
- 包级别的变量初始化、`init`函数执行，这两个操作都是在同一个`goroutine`中调用的，按顺序调用，一次一个包





## 总结

好啦，这篇文章到这里就结束了，本身`init`函数就很好理解，写这篇文章的目的就是让大家了解他的执行顺序，这样在日常开发中才不会写出`bug`。希望本文对大家有所帮助，我们下期见！

**素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！我是`asong`，我们下期见。**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E6%A0%87%E5%87%86%E8%89%B2%E7%89%88.png)

推荐往期文章：

- [Go语言如何实现可重入锁？](https://mp.weixin.qq.com/s/wBp4k7pJLNeSzyLVhGHLEA)
- [Go语言中new和make你使用哪个来分配内存？](https://mp.weixin.qq.com/s/XJ9O9O4KS3LbZL0jYnJHPg)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/mzSCWI8C_ByIPbb07XYFTQ)
- [空结构体引发的大型打脸现场](https://mp.weixin.qq.com/s/dNeCIwmPei2jEWGF6AuWQw)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/wURQFRt2ISz66icW7jbHFw)
- [面试官：两个nil比较结果是什么？](https://mp.weixin.qq.com/s/Dt46eoEXXXZc2ymr67_LVQ)
- [面试官：你能用Go写段代码判断当前系统的存储方式吗?](https://mp.weixin.qq.com/s/ffEsTpO-tyNZFR5navAbdA)
- [如何平滑切换线上Elasticsearch索引](https://mp.weixin.qq.com/s/8VQxK_Xh-bkVoOdMZs4Ujw)