## 背景

哈喽，大家好，我是正在学习`PS`技术的`asong`，上周读者问了我一道题，觉得挺有意义的，就在这里分享一下，我们先来看一下这个题：

```go
type User struct {
}

func FPrint(u User) {
	fmt.Printf("FPrint %p\n", &u)
}

func main() {
	u := User{}
	FPrint(u)
	fmt.Printf("main: %p\n", &u)
}
// 运行结果
FPrint 0x118eff0
main: 0x118eff0
```

看了运行结果，大多数朋友应该和我一样，一脸懵逼？`Go`语言不是只有值传递嘛？之前我还写过一篇关于"[Go语言参数传递是传值还是传引用吗？](https://mp.weixin.qq.com/s/JHbFh2GhoKewlemq7iI59Q)"，已经得出明确的结论，`Go`语言的确是只有值传递，这不是打脸了嘛。。。

既然已经出现了这样的结果，那么就要给出一个合理的解释，不要再让气氛尴尬下去，于是我给出了我的猜想，如下：

- 猜想一：这是一个`bug`
- 猜想二：结构体的特殊特性导致的

猜想一有点天马行空的感觉，暂时也无法验证，所以我们先来验证猜想二，请开始我的表演，都坐下，我要装逼了。。。。



## 验证猜想二：结构体的特殊特性导致的

上面的那道题中传参是一个空结构体，如果改成一个带字段的结构体会是什么样呢？我们来看一下：

```go
type UserIsEmpty struct {
}

type UserHasField struct {
	Age uint64 `json:"age"`
}

func FPrint(uIsEmpty UserIsEmpty, uHasField UserHasField) {
	fmt.Printf("FPrint uIsEmpty:%p uHasField:%p\n", &uIsEmpty, &uHasField)
}

func main() {
	uIsEmpty := UserIsEmpty{}
	uHasField := UserHasField{
		Age: 10,
	}
	FPrint(uIsEmpty, uHasField)
	fmt.Printf("main: uIsEmpty:%p uHasField:%p\n", &uIsEmpty, &uHasField)
}
// 运行结果：
FPrint uIsEmpty:0x118fff0 uHasField:0xc0000ba008
main: uIsEmpty:0x118fff0 uHasField:0xc0000ba000
```

从结果我们可以看出来，带字段的结构体确实是值传递，那么就证明空结构体有猫腻，有进展了，带着这个线索，我们来看一看这段代码的汇编部分，执行`go tool compile -N -l -S test.go`，可以得到汇编部分，截取重要部分：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-03-01%20%E4%B8%8B%E5%8D%888.32.47.png)



从结果上我们看到有调用`runtime.newobject(SB)`来进行分配内存，顺着这个在`runtme/malloc.go`中找到了他的实现：

```go
func newobject(typ *_type) unsafe.Pointer {
	return mallocgc(typ.size, typ, true)
}
```

`newobject()`中主要是调用了`mallocgc()`方法，在这里我找到了答案。因为`mallocgc()`代码比较长，这里我截取关键部分：

```go
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
	if gcphase == _GCmarktermination {
		throw("mallocgc called with gcphase == _GCmarktermination")
	}

	if size == 0 {
		return unsafe.Pointer(&zerobase)
	}
..........
}
```

如果 `size` 为 `0` 的时候，统一返回的都是全局变量 `zerobase` 的地址。到这里可能还会有一些伙伴有疑惑，这个跟上面的题有什么关系？那是因为你还不知道一个知识点：正常`struct`是占用一小块内存的，并且结构体的大小是要经过边界，长度的对齐的，但是“空结构体”是不占内存的，`size `为` 0`。现在一切都可以说的清了，总结原因：

> 因为空结构体是不占用内存的，所以`size`为0，在内存分配时，`size`为`0`会统一返回`zerobase`的地址，所以空结构体在进行参数传递时，发生值拷贝后地址都是一样的，才造成了这个质疑`Go`不是值传递的假象。



## 空结构体特性延伸

既然说到了空结构体，就在这里补充一个关于空结构体的知识点：空结构体做为结构体内置字段时是否进行内存对齐。

先来看一个例子：

```go
func main(){	
	fmt.Println(unsafe.Sizeof(Test1{}))
	fmt.Println(unsafe.Sizeof(Test2{}))
	fmt.Println(unsafe.Sizeof(Test3{}))

}

type Test1 struct {
	s struct{}
	n byte
	m byte
}

type Test2 struct {
	n byte
	s struct{}
	c byte
}

type Test3 struct {
	b byte
	s struct{}
}
//运行结果
2
2
2
```

根据运行结果我们可以得出结论：

空结构体在结构体中的前面和中间时，是不占用空间的，但是当空结构体放到结构体中的最后时，会进行特殊填充，`struct { }` 作为最后一个字段，会被填充对齐到前一个字段的大小，地址偏移对齐规则不变；

## 总结

最后做一个全文总结吧：

1. 空结构体也是一个结构体，不过他的`size`为0，所有的空结构体内存分配都是同一个地址，都是`zerobase`的地址；
2. 空结构体作为内嵌字段时要注意放置的顺序，当作为最后一个字段时会进行特殊填充，会被填充对齐到前一个字段的大小，地址偏移对齐规则不变；

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
- [真的理解interface了嘛](https://mp.weixin.qq.com/s/sO6Phr9C5VwcSTQQjJux3g)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/wURQFRt2ISz66icW7jbHFw)
- [十张动图带你搞懂排序算法(附go实现代码)](https://mp.weixin.qq.com/s/rZBsoKuS-ORvV3kML39jKw)
- [go参数传递类型](https://mp.weixin.qq.com/s/JHbFh2GhoKewlemq7iI59Q)
- [手把手教姐姐写消息队列](https://mp.weixin.qq.com/s/0MykGst1e2pgnXXUjojvhQ)
- [常见面试题之缓存雪崩、缓存穿透、缓存击穿](https://mp.weixin.qq.com/s?__biz=MzIzMDU0MTA3Nw==&mid=2247483988&idx=1&sn=3bd52650907867d65f1c4d5c3cff8f13&chksm=e8b0902edfc71938f7d7a29246d7278ac48e6c104ba27c684e12e840892252b0823de94b94c1&token=1558933779&lang=zh_CN#rd)
- [详解Context包，看这一篇就够了！！！](https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg)
- [高并发系统的限流策略：漏桶和令牌桶(附源码剖析)](https://mp.weixin.qq.com/s/fURwiSTeEE_Wvc95Q_fHnA)
- [面试官：go中for-range使用过吗？这几个问题你能解释一下原因吗](https://mp.weixin.qq.com/s/G7z80u83LTgLyfHgzgrd9g)