`欢迎大家点击上方文字「Golang梦工厂」关注公众号，设为星标，第一时间接收推送文章。`

## 前言

> 哈喽，大家好，我是`asong`，今天给大家介绍一个并发编程包`errgroup`，其实这个包就是对`sync.waitGroup`的封装。我们在之前的文章—— [源码剖析sync.WaitGroup(文末思考题你能解释一下吗?)](https://mp.weixin.qq.com/s/hofXXzFhu-rk3_6i2X4m6A)，从源码层面分析了`sync.WaitGroup`的实现，使用`waitGroup`可以实现一个`goroutine`等待一组`goroutine`干活结束，更好的实现了任务同步，但是`waitGroup`却无法返回错误，当一组`Goroutine`中的某个`goroutine`出错时，我们是无法感知到的，所以`errGroup`对`waitGroup`进行了一层封装，封装代码仅仅不到`50`行，下面我们就来看一看他是如何封装的？



## `errGroup`如何使用

老规矩，我们先看一下`errGroup`是如何使用的，前面吹了这么久，先来验验货；

以下来自官方文档的例子：

```go
var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Result string
type Search func(ctx context.Context, query string) (Result, error)

func fakeSearch(kind string) Search {
	return func(_ context.Context, query string) (Result, error) {
		return Result(fmt.Sprintf("%s result for %q", kind, query)), nil
	}
}

func main() {
	Google := func(ctx context.Context, query string) ([]Result, error) {
		g, ctx := errgroup.WithContext(ctx)

		searches := []Search{Web, Image, Video}
		results := make([]Result, len(searches))
		for i, search := range searches {
			i, search := i, search // https://golang.org/doc/faq#closures_and_goroutines
			g.Go(func() error {
				result, err := search(ctx, query)
				if err == nil {
					results[i] = result
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return nil, err
		}
		return results, nil
	}

	results, err := Google(context.Background(), "golang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}

}
```

上面这个例子来自官方文档，代码量有点多，但是核心主要是在`Google`这个闭包中，首先我们使用`errgroup.WithContext`创建一个`errGroup`对象和`ctx`对象，然后我们直接调用`errGroup`对象的`Go`方法就可以启动一个协程了，`Go`方法中已经封装了`waitGroup`的控制操作，不需要我们手动添加了，最后我们调用`Wait`方法，其实就是调用了`waitGroup`方法。这个包不仅减少了我们的代码量，而且还增加了错误处理，对于一些业务可以更好的进行并发处理。



## 赏析`errGroup`

### 数据结构

我们先看一下`Group`的数据结构：

```go
type Group struct {
	cancel func() // 这个存的是context的cancel方法

	wg sync.WaitGroup // 封装sync.WaitGroup

	errOnce sync.Once // 保证只接受一次错误
	err     error // 保存第一个返回的错误
}
```



### 方法解析

```go
func WithContext(ctx context.Context) (*Group, context.Context)
func (g *Group) Go(f func() error)
func (g *Group) Wait() error
```

`errGroup`总共只有三个方法：

- `WithContext`方法

```go
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}
```

这个方法只有两步：

- 使用`context`的`WithCancel()`方法创建一个可取消的`Context`
- 创建`cancel()`方法赋值给`Group`对象



- `Go`方法

```go
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}
```

`Go`方法中运行步骤如下：

- 执行`Add()`方法增加一个计数器
- 开启一个协程，运行我们传入的函数`f`，使用`waitGroup`的`Done()`方法控制是否结束
- 如果有一个函数`f`运行出错了，我们把它保存起来，如果有`cancel()`方法，则执行`cancel()`取消其他`goroutine`

这里大家应该会好奇为什么使用`errOnce`，也就是`sync.Once`，这里的目的就是保证获取到第一个出错的信息，避免被后面的`Goroutine`的错误覆盖。



- `wait`方法

```go
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
```

总结一下`wait`方法的执行逻辑：

- 调用`waitGroup`的`Wait()`等待一组`Goroutine`的运行结束
- 这里为了保证代码的健壮性，如果前面赋值了`cancel`，要执行`cancel()`方法
- 返回错误信息，如果有`goroutine`出现了错误才会有值



### 小结

到这里我们就分析完了`errGroup`包，总共就`1`个结构体和`3`个方法，理解起来还是比较简单的，针对上面的知识点我们做一个小结：

- 我们可以使用`withContext`方法创建一个可取消的`Group`，也可以直接使用一个零值的`Group`或`new`一个`Group`，不过直接使用零值的`Group`和`new`出来的`Group`出现错误之后就不能取消其他`Goroutine`了。
- 如果多个`Goroutine`出现错误，我们只会获取到第一个出错的`Goroutine`的错误信息，晚于第一个出错的`Goroutine`的错误信息将不会被感知到。
- `errGroup`中没有做`panic`处理，我们在`Go`方法中传入`func() error`方法时要保证程序的健壮性



## 踩坑日记

使用`errGroup`也并不是一番风顺的，我之前在项目中使用`errGroup`就出现了一个`BUG`，把它分享出来，避免踩坑。

这个需求是这样的(并不是真实业务场景，由`asong`虚构的)：开启多个`Goroutine`去缓存中设置数据，同时开启一个`Goroutine`去异步写日志，很快我的代码就写出来了：

```go
func main()  {
	g, ctx := errgroup.WithContext(context.Background())

	// 单独开一个协程去做其他的事情，不参与waitGroup
	go WriteChangeLog(ctx)

	for i:=0 ; i< 3; i++{
		g.Go(func() error {
			return errors.New("访问redis失败\n")
		})
	}
	if err := g.Wait();err != nil{
		fmt.Printf("appear error and err is %s",err.Error())
	}
	time.Sleep(1 * time.Second)
}

func WriteChangeLog(ctx context.Context) error {
	select {
	case <- ctx.Done():
		return nil
	case <- time.After(time.Millisecond * 50):
		fmt.Println("write changelog")
	}
	return nil
}
// 运行结果
appear error and err is 访问redis失败
```

代码没啥问题吧，但是日志一直没有写入，排查了好久，终于找到问题原因。原因就是这个`ctx`。

因为这个`ctx`是`WithContext`方法返回的一个带取消的`ctx`，我们把这个`ctx`当作父`context`传入`WriteChangeLog`方法中了，如果`errGroup`取消了，也会导致上下文的`context`都取消了，所以`WriteChangelog`方法就一直执行不到。

这个点是我们在日常开发中想不到的，所以需要注意一下～。



## 总结

因为最近看很多朋友都不知道这个库，所以今天就把他分享出来了，封装代码仅仅不到`50`行，真的是很厉害，如果让你来封装，你能封装的更好吗？

**素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！我是`asong`，我们下期见。**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E7%99%BD%E8%89%B2%E7%89%88-20210717170231906.png)

推荐往期文章：

- [学习channel设计：从入门到放弃](https://mp.weixin.qq.com/s/E2XwSIXw1Si1EVSO1tMW7Q)
- [编程模式之Go如何实现装饰器](https://mp.weixin.qq.com/s/B_VYr3I525-vjHgzfW3Jhg)
- [Go语言中new和make你使用哪个来分配内存？](https://mp.weixin.qq.com/s/xNdnVXxC5Ji2ApgbfpRaXQ)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/yJ05a6pNxr_G72eiWTJ-rw)
- [空结构体引发的大型打脸现场](https://mp.weixin.qq.com/s/aHwGWWmnDFkcw2cyw5jmgw)
- [面试官：你能聊聊string和[]byte的转换吗？](https://mp.weixin.qq.com/s/jztwFH6thFdcySzowXOH_Q)
- [面试官：两个nil比较结果是什么？](https://mp.weixin.qq.com/s/CNOLLLRzHomjBnbZMnw0Gg)
- [面试官：你能用Go写段代码判断当前系统的存储方式吗?](https://mp.weixin.qq.com/s/DWMqzOi7wf79DoUUAJnr1w)
- [赏析Singleflight设计](https://mp.weixin.qq.com/s/JUkxGbx1Ufpup3Hx08tI2w)

