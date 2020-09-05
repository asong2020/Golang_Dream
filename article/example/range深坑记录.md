## 前言

读者A：不会吧，阿Sir，这周这么高产～～～

asong：当然啦，为了你们，一切都值得～～～

读者B：净放臭屁屁，就你戏多～～～

asong：你凶人家，坏坏～～～

哈哈哈，戏太足了奥。自导自演可还行。今日分享之前，先放松放松嘛，毕竟接下来的知识，还是需要我们思考的。今天给大家分享的是go中的range，这个我们在实际开发中，是经常使用，但是他有一个坑，使用不好，是要被开除的。但是，今天你恰好看了我这一篇文章，就避免了这个坑，开心嘛～～～。直接笑，别克制，我知道你嘴角已经上扬了。

废话结束，我们直接开始。



### 正文

#### 1. 指针数据坑

range到底有什么坑呢，我们先来运行一个例子吧。

```go
package main

import (
	"fmt"
)

type user struct {
	name string
	age uint64
}

func main()  {
	u := []user{
		{"asong",23},
		{"song",19},
		{"asong2020",18},
	}
	n := make([]*user,0,len(u))
	for _,v := range u{
		n = append(n, &v)
	}
	fmt.Println(n)
	for _,v := range n{
		fmt.Println(v)
	}
}
```

这个例子的目的是，通过`u`这个slice构造成新的slice。我们预期应该是显示`u`slice的内容，但是运行结果如下：

```go
[0xc0000a6040 0xc0000a6040 0xc0000a6040]
&{asong2020 18}
&{asong2020 18}
&{asong2020 18}
```

这里我们看到`n`这个slice打印出来的三个同样的数据，并且他们的内存地址相同。这是什么原因呢？先别着急，再来看这一段代码，我给他改正确他，对比之后我们再来分析，你们才会恍然大悟。

```go
package main

import (
	"fmt"
)

type user struct {
	name string
	age uint64
}

func main()  {
	u := []user{
		{"asong",23},
		{"song",19},
		{"asong2020",18},
	}
	n := make([]*user,0,len(u))
	for _,v := range u{
		o := v
		n = append(n, &o)
	}
	fmt.Println(n)
	for _,v := range n{
		fmt.Println(v)
	}
}
```

细心的你们看到，我改动了哪一部分代码了嘛？对，没错，我就加了一句话，他就成功了，我在`for range`里面引入了一个中间变量，每次迭代都重新声明一个变量`o`，赋值后再将`v`的地址添加`n`切片中，这样成功解决了刚才的问题。

现在来解释一下原因：在`for range`中，变量`v`是用来保存迭代切片所得的值，因为`v`只被声明了一次，每次迭代的值都是赋值给`v`，该变量的内存地址始终未变，这样讲他的地址追加到新的切片中，该切片保存的都是同一个地址，这肯定无法达到预期效果的。这里还需要注意一点，变量`v`的地址也并不是指向原来切片`u[2]`的，因我在使用`range`迭代的时候，变量`v`的数据是切片的拷贝数据，所以直接`copy`了结构体数据。

上面的问题还有一种解决方法，直接引用数据的内存，这个方法比较好，不需要开辟新的内存空间，看代码：

```go
......略
for k,_ := range u{
		n = append(n, &u[k])
	}
......略
```



#### 2. 迭代修改变量问题

还是刚才的例子，我们做一点改动，现在我们要对切片中保存的每个用户的年龄进行修改，因为我们都是永远18岁，嘎嘎嘎～～～。

```go
package main

import (
	"fmt"
)

type user struct {
	name string
	age uint64
}

func main()  {
	u := []user{
		{"asong",23},
		{"song",19},
		{"asong2020",18},
	}
	for _,v := range u{
		if v.age != 18{
			v.age = 20
		}
	}
	fmt.Println(u)
}
```

来看一下运行结果：

```go
[{asong 23} {song 19} {asong2020 18}]
```

哎呀，怎么回事。怎么没有更改呢。其实道理都是一样，还记得，我在上文说的一个知识点嘛。对，就是这个，想起来了吧。`v`变量是拷贝切片中的数据，修改拷贝数据怎么会对原切片有影响呢，还是这个问题，`copy`这个知识点很重要，一不注意，就会出现问题。知道问题了，我们现在来把这个问题解决吧。

```go
package main

import (
	"fmt"
)

type user struct {
	name string
	age uint64
}

func main()  {
	u := []user{
		{"asong",23},
		{"song",19},
		{"asong2020",18},
	}
	for k,v := range u{
		if v.age != 18{
			u[k].age = 18
		}
	}
	fmt.Println(u)
}
```

可以看到，我们直接对切片的值进行修改，这样就修改成功了。所以这里还是要注意一下的，防止以后出现`bug`。



#### 3. 是否会造成死循环

来看一段代码：

```go
func main() {
	v := []int{1, 2, 3}
	for i := range v {
		v = append(v, i)
	}
}
```

这一段代码会造成死循环吗？答案：当然不会，前面都说了`range`会对切片做拷贝，新增的数据并不在拷贝内容中，并不会发生死循环。这种题一般会在面试中问，可以留意下的。



### 你不知道的`range`用法

#### `delete`

没看错，删除，在`range`迭代时，可以删除`map`中的数据，第一次见到这么使用的，我刚听到确实不太相信，所以我就去查了一下官方文档，确实有这个写法：

```go
for key := range m {
    if key.expired() {
        delete(m, key)
    }
}
```

看看官方的解释：

```
The iteration order over maps is not specified and is not guaranteed to be the same from one iteration to the next. If map entries that have not yet been reached are removed during iteration, the corresponding iteration values will not be produced. If map entries are created during iteration, that entry may be produced during the iteration or may be skipped. The choice may vary for each entry created and from one iteration to the next. If the map is nil, the number of iterations is 0.

翻译：
未指定`map`的迭代顺序，并且不能保证每次迭代之间都相同。 如果在迭代过程中删除了尚未到达的映射条目，则不会生成相应的迭代值。 如果映射条目是在迭代过程中创建的，则该条目可能在迭代过程中产生或可以被跳过。 对于创建的每个条目以及从一个迭代到下一个迭代，选择可能有所不同。 如果映射为nil，则迭代次数为0。
```

看这个代码：

```go
func main()  {
	d := map[string]string{
		"asong": "帅",
		"song": "太帅了",
	}
	for k := range d{
		if k == "asong"{
			delete(d,k)
		}
	}
	fmt.Println(d)
}

# 运行结果
map[song:太帅了]
```

从运行结果我们可以看出，key为`asong`的这位帅哥被从帅哥`map`中删掉了，哇哦，可气呀。这个方法，相信很多小伙伴都不知道，今天教给你们了，以后可以用起来了。



#### add

上面是删除，那肯定会有新增呀，直接看代码吧。

```go
func main()  {
	d := map[string]string{
		"asong": "帅",
		"song": "太帅了",
	}
	for k,v := range d{
		d[v] = k
		fmt.Println(d)
	}
}
```

这里我把打印放到了`range`里，你们思考一下，新增的元素，在遍历时能够遍历到呢。我们来验证一下。

```
func main()  {
	var addTomap = func() {
		var t = map[string]string{
			"asong": "太帅",
			"song": "好帅",
			"asong1": "非常帅",
		}
		for k := range t {
			t["song2020"] = "真帅"
			fmt.Printf("%s%s ", k, t[k])
		}
	}
	for i := 0; i < 10; i++ {
		addTomap()
		fmt.Println()
	}
}
```

运行结果：

```go
asong太帅 song好帅 asong1非常帅 song2020真帅 
asong太帅 song好帅 asong1非常帅 
asong太帅 song好帅 asong1非常帅 song2020真帅 
asong1非常帅 song2020真帅 asong太帅 song好帅 
asong太帅 song好帅 asong1非常帅 song2020真帅 
asong太帅 song好帅 asong1非常帅 song2020真帅 
asong太帅 song好帅 asong1非常帅 
asong1非常帅 song2020真帅 asong太帅 song好帅 
asong太帅 song好帅 asong1非常帅 song2020真帅 
asong太帅 song好帅 asong1非常帅 song2020真帅
```

从运行结果，我们可以看出来，每一次的结果并不是确定的。这是为什么呢？这就来揭秘，map内部实现是一个链式hash表，为了保证无顺序，初始化时会随机一个遍历开始的位置，所以新增的元素被遍历到就变的不确定了，同样删除也是一个道理，但是删除元素后边就不会出现，所以一定不会被遍历到。



## 总结

怎么样，伙伴们，收获不小吧。一个小小的`range`就会引发这么多的问题，所以说写代码一定要实践，光靠想是没有用的，有些问题只有在实践中才会有所提高。希望今天的分享对你们有用，好啦，这一期就结束啦。我们下期见。打个预告：下期将介绍go-elastic的使用，有需要的小伙伴留意一下。

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**我是asong，一名普普通通的程序猿，让我一起慢慢变强吧。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [学会wire依赖注入、cron定时任务其实就这么简单！](https://mp.weixin.qq.com/s/qmbCmwZGmqKIZDlNs_a3Vw)

- [听说你还不会jwt和swagger-饭我都不吃了带着实践项目我就来了](https://mp.weixin.qq.com/s/z-PGZE84STccvfkf8ehTgA)
- [掌握这些Go语言特性，你的水平将提高N个档次(二)](https://mp.weixin.qq.com/s/7yyo83SzgQbEB7QWGY7k-w)
- [go实现多人聊天室，在这里你想聊什么都可以的啦！！！](https://mp.weixin.qq.com/s/H7F85CncQNdnPsjvGiemtg)
- [grpc实践-学会grpc就是这么简单](https://mp.weixin.qq.com/s/mOkihZEO7uwEAnnRKGdkLA)
- [go标准库rpc实践](https://mp.weixin.qq.com/s/d0xKVe_Cq1WsUGZxIlU8mw)
- [2020最新Gin框架中文文档 asong又捡起来了英语，用心翻译](https://mp.weixin.qq.com/s/vx8A6EEO2mgEMteUZNzkDg)
- [基于gin的几种热加载方式](https://mp.weixin.qq.com/s/CZvjXp3dimU-2hZlvsLfsw)
- [boss: 这小子还不会使用validator库进行数据校验，开了～～～](https://mp.weixin.qq.com/s?__biz=MzIzMDU0MTA3Nw==&mid=2247483829&idx=1&sn=d7cf4f46ea038a68e74a4bf00bbf64a9&scene=19&token=1606435091&lang=zh_CN#wechat_redirect)

