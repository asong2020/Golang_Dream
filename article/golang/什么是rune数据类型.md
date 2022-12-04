## 背景

> 哈喽，大家好，我是`asong`。今天我们一起来看看`Go`语言中的`rune`数据类型，首先从一道面试题入手，你能很快说出下面这道题的答案吗？

```go
func main()  {
	str := "Golang梦工厂"
	fmt.Println(len(str))
	fmt.Println(len([]rune(str)))
}
```

运行结果是`15`和`15`还是`15`和`9`呢？先思考一下，一会揭晓答案。

其实这并不是一道面试题，是我在日常开发中遇到的一个问题，当时场景是这样的：后端要对前端传来的字符串做字符校验，产品的需求是限制为200字符，然后我在后端做校验时直接使用`len(str) > 200`来做判断，结果出现了`bug`，前端字符校验没有超过`200`字符，调用后端接口确一直是参数错误，改成使用`len([]rune(str)) > 200`成功解决了这个问题。具体原因我们在文中揭晓。



## `Unicode`和字符编码

在介绍`rune`类型之前，我们还是要从一些基础知识开始。 ------ `Unicode`和字符编码。

- 什么是`Unicode`？

我们都知道计算机只能处理数字，如果想要处理文本需要转换为数字才能处理，早些时候，计算机在设计上采用`8bit`作为一个`byte`，一个`byte`表示的最大整数就是`255`，想表示更大的整数，就需要更多的`byte`。显然，一个字节表示中文，是不够的，至少需要两个字节，而且还不能和ASCII编码冲突，所以，我国制定了`GB2312`编码，用来把中文编进去。但是世界上有很多语言，不同语言制定一个编码，就会不可避免地出现冲突，所以`unicode`字符就是来解决这个痛点的。`Unicode`把所有语言都统一到一套编码里。总结来说："**unicode其实就是对字符的一种编码方式，可以理解为一个字符---数字的映射机制，利用一个数字即可表示一个字符。**"

- 什么是字符编码？

虽然`unicode`把所有语言统一到一套编码里了，但是他却没有规定字符对应的二进制码是如何存储。以汉字“汉”为例，它的 `Unicode` 码点是 `0x6c49`，对应的二进制数是 `110110001001001`，二进制数有 `15 `位，这也就说明了它至少需要 `2 `个字节来表示。可以想象，在` Unicode` 字典中往后的字符可能就需要 `3 `个字节或者 `4 `个字节，甚至更多字节来表示了。

这就导致了一些问题，计算机怎么知道你这个` 2 `个字节表示的是一个字符，而不是分别表示两个字符呢？这里我们可能会想到，那就取个最大的，假如 `Unicode `中最大的字符用` 4` 字节就可以表示了，那么我们就将所有的字符都用` 4 `个字节来表示，不够的就往前面补` 0`。这样确实可以解决编码问题，但是却造成了空间的极大浪费，如果是一个英文文档，那文件大小就大出了` 3` 倍，这显然是无法接受的。

于是，为了较好的解决` Unicode` 的编码问题， `UTF-8` 和` UTF-16` 两种当前比较流行的编码方式诞生了。`UTF-8` 是目前互联网上使用最广泛的一种` Unicode `编码方式，它的最大特点就是可变长。它可以使用 `1 - 4 `个字节表示一个字符，根据字符的不同变换长度。在`UTF-8`编码中，一个英文为一个字节，一个中文为三个字节。



## `Go`语言中的字符串

### 基本概念

先来看一下官方对`string`的定义：

```go
// string is the set of all strings of 8-bit bytes, conventionally but not
// necessarily representing UTF-8-encoded text. A string may be empty, but
// not nil. Values of string type are immutable.
type string string
```

人工翻译：

> `string`是`8`位字节的集合，通常但不一定代表`UTF-8`编码的文本。`string`可以为空，但不能为`nil`。**`string`的值是不能改变的**

说得通俗一点，其实字符串实际上是只读的字节切片，对于字符串底层而言就是一个`byte`数组，不过这个数组是只读的，不允许修改。

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-04-25%20%E4%B8%8B%E5%8D%883.07.31.png)

写个例子验证一下：

```go
func main()  {
	byte1 := []byte("Hl Asong!")
	byte1[1] = 'i'

	str1 := "Hl Asong!"
	str1[1] = 'i'
}
```

对于`byte`的操作是可行，而`string`操作会直接报错：

> cannot assign to str1[1]

所以说`string`修改操作是不允许的，仅仅支持替换操作。

根据前面的分析，我们也可以得出我们将字符存储在字符串中时，也就是按字节进行存储的，所以最后存储的其实是一个数值。



### `Go`语言的字符串编码

上面我们介绍了字符串的基本概念，接下来我们看一下`Go`语言中的字符串编码是怎样的。

`Go `源代码为 `UTF-8` 编码格式的，源代码中的字符串直接量是 `UTF-8` 文本。所以`Go`语言中字符串是`UTF-8`编码格式的。



###  `Go`语言字符串循环

`Go`语言中字符串可以使用`range`循环和下标循环。我们写一个例子，看一下两种方式循环有什么区别：

```go
func main()  {
	str := "Golang梦工厂"
	for k,v := range str{
		fmt.Printf("v type: %T index,val: %v,%v \n",v,k,v)
	}
	for i:=0 ; i< len(str) ; i++{
		fmt.Printf("v type: %T index,val:%v,%v \n",str[i],i,str[i])
	}
}
```

运行结果：

```go
v type: int32 index,val: 0,71 
v type: int32 index,val: 1,111 
v type: int32 index,val: 2,108 
v type: int32 index,val: 3,97 
v type: int32 index,val: 4,110 
v type: int32 index,val: 5,103 
v type: int32 index,val: 6,26790 
v type: int32 index,val: 9,24037 
v type: int32 index,val: 12,21378 
v type: uint8 index,val:0,71 
v type: uint8 index,val:1,111 
v type: uint8 index,val:2,108 
v type: uint8 index,val:3,97 
v type: uint8 index,val:4,110 
v type: uint8 index,val:5,103 
v type: uint8 index,val:6,230 
v type: uint8 index,val:7,162 
v type: uint8 index,val:8,166 
v type: uint8 index,val:9,229 
v type: uint8 index,val:10,183 
v type: uint8 index,val:11,165 
v type: uint8 index,val:12,229 
v type: uint8 index,val:13,142 
v type: uint8 index,val:14,130
```

根据运行结果我们可以得出如下结论：

> 使用下标遍历获取的是`ASCII`字符，而使用`Range`遍历获取的是`Unicode`字符。



## 什么是`rune`数据类型

官方对`rune`的定义如下：

```go
// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32
```

人工翻译：

> `rune`是`int32`的别名，在所有方面都等同于`int32`，按照约定，它用于区分字符值和整数值。

说的通俗一点就是`rune`一个值代表的就是一个`Unicode`字符，因为一个`Go`语言中字符串编码为`UTF-8`，使用`1-4`字节就可以表示一个字符，所以使用`int32`类型范围就可以完美适配。



## 答案揭晓

前面说了这么多知识点，确实有点乱了，我们现在就根据开始的那道题来做一个总结。为了方便查看，在贴一下这道题：

```go
func main()  {
	str := "Golang梦工厂"
	fmt.Println(len(str))
	fmt.Println(len([]rune(str)))
}
```

这道题的正确答案是`15`和`9`。

具体原因：

> `len()`函数是用来获取字符串的字节长度，`rune`一个值代表的就是一个`Unicode`字符，所以求`rune`切片的长度就是字符个数。因为在`utf-8`编码中，英文占`1`个字节，中文占`3`个字节，所以最终结果就是`15`和`9`。

贴个图，方便理解：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-04-25%20%E4%B8%8B%E5%8D%886.08.35.png)



## `unicode/utf8`库

如果大家对`rune`的使用不是很明确，可以学习使用一下`Go`标准库`unicode/utf8`，其中提供了多种关于`rune`的使用方法。比如上面这道题，我们就可以使用`utf8.RuneCountInString`方法获取字符个数。更多库函数使用方法请自行解锁，本篇就不做过多介绍了。



## 总结

针对全文，我们做一个总结：

- Go语言源代码始终为`UTF-8`
- `Go`语言的字符串可以包含任意字节，字符底层是一个只读的`byte`数组。
- `Go`语言中字符串可以进行循环，使用下表循环获取的`acsii`字符，使用`range`循环获取的`unicode`字符。
- `Go`语言中提供了`rune`类型用来区分字符值和整数值，一个值代表的就是一个`Unicode`字符。
- `Go`语言中获取字符串的字节长度使用`len()`函数，获取字符串的字符个数使用`utf8.RuneCountInString`函数或者转换为`rune`切片求其长度，这两种方法都可以达到预期结果。

**好啦，这篇文章就到这里啦，素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

**我是asong，一名普普通通的程序猿，让我们一起慢慢变强吧。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E7%99%BD%E8%89%B2%E7%89%88-20210425182909397.png)

推荐往期文章：

- [Go看源码必会知识之unsafe包](https://mp.weixin.qq.com/s/nPWvqaQiQ6Z0TaPoqg3t2Q)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/mzSCWI8C_ByIPbb07XYFTQ)
- [空结构体引发的大型打脸现场](https://mp.weixin.qq.com/s/dNeCIwmPei2jEWGF6AuWQw)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/wURQFRt2ISz66icW7jbHFw)
- [面试官：两个nil比较结果是什么？](https://mp.weixin.qq.com/s/Dt46eoEXXXZc2ymr67_LVQ)
- [面试官：你能用Go写段代码判断当前系统的存储方式吗?](https://mp.weixin.qq.com/s/ffEsTpO-tyNZFR5navAbdA)
- [如何平滑切换线上Elasticsearch索引](https://mp.weixin.qq.com/s/8VQxK_Xh-bkVoOdMZs4Ujw)