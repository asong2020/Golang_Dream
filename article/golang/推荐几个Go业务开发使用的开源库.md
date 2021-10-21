## 前言

> 哈喽，大家好，我是`asong`。拖更了好久，这周开始更新。
>
> 最近总有一些初学`Go`语言的小伙伴问我在业务开发中一般都使用什么web框架、开源中间件；所以我总结了我在日常开发中使用到的库，这些库不一定是特别完美的，但是基本可以解决日常工作需求，接下来我们就来看一下。



## `Gin`

`Gin`是一个用`Go`编写的`Web`框架，它是一个类似于`martini`但拥有更好性能的`API`框架。基本现在每个`Go`初学者学习的第一个`web`框架都是`Gin`。在网上看到一个关于对各个Go-web框架受欢迎的对比：

![来自网络](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/1.png)

我们可以看到`Gin`在社区受欢迎排第一，Gin 框架往往是进行 Web 应用开发的首选框架，许多公司都会选择采用`Gin`框架进行二次开发，加入日志，服务发现等功能，像Bilibili 开源的一套 Go 微服务框架 Kratos 就采用 Gin 框架进行了二次开发。

学习`Gin`通过他的官方文档就可以很快入手，不过文档是英文的，这个不用担心，我曾翻译了一份中文版，可以到我的公众号后台获取，回复【gin】即可获取。

github地址：https://github.com/gin-gonic/gin



## `zap`

`zap`是`uber`开源的日志库，选择`zap`他有两个优势：

- 它非常的快
- 它同时提供了结构化日志记录和printf风格的日志记录

大多数日志库基本都是基于反射的序列化和字符串格式化的，这样会导致在日志上占用大量`CPU`资源，不适用于业务开发场景，业务对性能敏感还是挺高的。`zap`采用了不同的方法，它设计了一个无反射、零分配的 JSON 编码器，并且基础 Logger 力求尽可能避免序列化开销和分配。 通过在此基础上构建高级 SugaredLogger，zap 允许用户选择何时需要计算每次分配以及何时更喜欢更熟悉的松散类型的 API。

`zap`的基准测试如下：

![来自官方文档](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/2.png)

可以看出`zap`的效率完全高于其他日志库，选谁不用我明说了吧！！！

github地址：https://github.com/uber-go/zap



## `jsoniter`

做业务开发离不开`json`的序列化与反序列化，标准库虽然提供了`encoding/json`，但是它主要是通过反射来实现的，所以性能消耗比较大。`jsoniter`可以解决这个痛点，其是一款快且灵活的 JSON 解析器，具有良好的性能并能100%兼容标准库，我们可以使用jsoniter替代encoding/json，官方文档称可以比标准库**快6倍**多，后来Go官方在go1.12版本对 json.Unmarshal 函数使用 sync.Pool 缓存了 decoder，性能较之前的版本有所提升，所以现在达不到**快6倍**多。

![来自官方文档](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/3.png)

github地址：https://github.com/json-iterator/go

对于`jsoniter`优化原理感兴趣的可以移步这里：http://jsoniter.com/benchmark.html#optimization-used



## `gorm`

`gorm`是一个使用`Go`语言编写的`ORM`框架，文档齐全，对开发者友好，并且支持主流的数据库：`MySQL`, `PostgreSQL`, `SQlite`, `SQL Server`。

个人觉得使用`gorm`最大的好处在于它是由国人开发，中文文档齐全，上手很快，目前大多数企业也都在使用`gorm`。我们来一下`gorm`的特性：

- 全功能 ORM
- 关联 (Has One，Has Many，Belongs To，Many To Many，多态，单表继承)
- Create，Save，Update，Delete，Find 中钩子方法
- 支持 `Preload`、`Joins` 的预加载
- 事务，嵌套事务，Save Point，Rollback To Saved Point
- Context、预编译模式、DryRun 模式
- 批量插入，FindInBatches，Find/Create with Map，使用 SQL 表达式、Context Valuer 进行 CRUD
- SQL 构建器，Upsert，数据库锁，Optimizer/Index/Comment Hint，命名参数，子查询
- 复合主键，索引，约束
- Auto Migration
- 自定义 Logger
- 灵活的可扩展插件 API：Database Resolver（多数据库，读写分离）、Prometheus…
- 每个特性都经过了测试的重重考验
- 开发者友好

github地址：https://github.com/go-gorm/gorm

官方文档：https://gorm.io/zh_CN/docs/index.html



## `robfig/cron`

github地址：https://github.com/robfig/cron

业务开发更离不开定时器的使用了，`cron`就是一个用于管理定时任务的库，用 Go 实现 Linux 中`crontab`这个命令的效果，与Linux 中`crontab`命令相似，`cron`库支持用 **5** 个空格分隔的域来表示时间。`cron`上手也是非常容易的，看一个官方的例子：

```go
package main

import (
  "fmt"
  "time"

  "github.com/robfig/cron/v3"
)

func main() {
  c := cron.New()

  c.AddFunc("@every 1s", func() {
    fmt.Println("tick every 1 second run once")
  })
  c.Start()
  time.Sleep(time.Second * 10)
}
```

针对`cron `的使用可以参考这篇文章：https://segmentfault.com/a/1190000023029219

之前我也写了一篇`cron`的基本使用，可以参考下：https://mp.weixin.qq.com/s/Z4B7Tn8ikFIkXVGhXNbsVA



## `wire`

都`1202`年了，应该不会有人不知道依赖注入的作用了吧。我们本身也可以自己实现依赖注入，但是这是在代码量少、结构不复杂的情况下，当结构之间的关系变得非常复杂的时候，这时候手动创建依赖，然后将他们组装起来就会变的异常繁琐，并且很容出错。Go语言社区有很多依赖注入的框架，可以分为两个类别：

- 依赖反射实现的运行时依赖注入：inject、uber、dig
- 使用代码生成实现的依赖注入：wire

个人觉得使用`wire`进行项目管理是最好的，在代码编译阶段就可以发现依赖注入的问题，在代码生成时即可报出来，不会拖到运行时才报，更便于` debug`。

`wire`的使用也是非常的简单，关于`wire`的使用我之前也写了一篇文章，可以参考一下：https://mp.weixin.qq.com/s/Z4B7Tn8ikFIkXVGhXNbsVA

github地址：https://github.com/google/wire



## `ants`

某些业务场景还会使用到`goroutine`池，`ants`就是一个广泛使用的goroute池，可以有效控制协程数量，防止协程过多影响程序性能。`ants`也是国人开发的，设计博文写的也很详细的，目前很多大厂也都在使用`ants`，经历过线上业务检验的，所以可以放心使用。

github地址：https://github.com/panjf2000/ants

`ants`源码不到`1k`行，建议大家赏析一下源码～。



## 总结

本文列举的几个库都是经常被使用的开源库，这几个库你都掌握了，基本的业务开发都没有啥问题了，一些初学者完全可以通过这几个库达到入门水平。还有一些库，比如：`go-redis`、`go-sql-driver`、`didi/gendry`、`golang/groupcache`、`olivere/elastic/v7 `等等，这些库也是经常使用的，入门都比较简单，就不这里详细介绍了。

如果大家也有经常使用的，比较好的开源库，欢迎推荐给我，我也学习学习！！！

**好啦，本文就到这里了，素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！我是`asong`，我们下期见。**

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)

推荐往期文章：

- [学习channel设计：从入门到放弃](https://mp.weixin.qq.com/s/E2XwSIXw1Si1EVSO1tMW7Q)
- [详解内存对齐](https://mp.weixin.qq.com/s/ig8LDNdpflEBWlypU1NRhw)
- [警惕请勿滥用goroutine](https://mp.weixin.qq.com/s/JC14dWffHub0nfPlPipsHQ)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/yJ05a6pNxr_G72eiWTJ-rw)
- [面试官：小松子来聊一聊内存逃逸](https://mp.weixin.qq.com/s/MepbrrSlGVhNrEkTQhfhhQ)
- [面试官：两个nil比较结果是什么？](https://mp.weixin.qq.com/s/CNOLLLRzHomjBnbZMnw0Gg)
- [Go语言如何操纵Kafka保证无消息丢失](https://mp.weixin.qq.com/s/XoSi3Cgp7ij-n9t4pvBoXQ)



