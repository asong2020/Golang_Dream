## 前言

> 哈喽，大家好，我是拖更好久的鸽子`asong`。因为`5.1`去找女朋友，所以一直没有时间写文章啦，想着回来就抓紧学习，无奈，依然沉浸在5.1的甜蜜生活中，一拖再拖，就到现在啦。果然女人影响了我拔刀的速度，但是我很喜欢，略略略。
>
> 好啦，不撒狗粮了，开始进入正题，今天我们就来探讨一下`Go`语言中的`make`和`new`到底怎么使用？它们又有什么不同？



## 分配内存之`new`

官方文档定义：
```go
// The new built-in function allocates memory. The first argument is a type,
// not a value, and the value returned is a pointer to a newly
// allocated zero value of that type.
func new(Type) *Type
```
翻译出来就是：`new`是一个分配内存的内置函数，第一个参数是类型，而不是值，返回的值是指向该类型新分配的零值的指针。
我们平常在使用指针的时候是需要分配内存空间的，未分配内存空间的指针直接使用会使程序崩溃，比如这样：
```go
var a *int64
*a = 10
```
我们声明了一个指针变量，直接就去使用它，就会使用程序触发`panic`，因为现在这个指针变量`a`在内存中没有块地址属于它，就无法直接使用该指针变量，所以`new`函数的作用就出现了，通过`new`来分配一下内存，就没有问题了：
```go
var a *int64 = new(int64)
	*a = 10
```

上面的例子，我们是针对普通类型`int64`进行`new`处理的，如果是复合类型，使用`new`会是什么样呢？来看一个示例：

```go
func main(){
	// 数组
	array := new([5]int64)
	fmt.Printf("array: %p %#v \n", &array, array)// array: 0xc0000ae018 &[5]int64{0, 0, 0, 0, 0}
	(*array)[0] = 1
	fmt.Printf("array: %p %#v \n", &array, array)// array: 0xc0000ae018 &[5]int64{1, 0, 0, 0, 0}
	
	// 切片
	slice := new([]int64)
	fmt.Printf("slice: %p %#v \n", &slice, slice) // slice: 0xc0000ae028 &[]int64(nil)
	(*slice)[0] = 1
	fmt.Printf("slice: %p %#v \n", &slice, slice) // panic: runtime error: index out of range [0] with length 0

	// map
	map1 := new(map[string]string)
	fmt.Printf("map1: %p %#v \n", &map1, map1) // map1: 0xc00000e038 &map[string]string(nil)
	(*map1)["key"] = "value"
	fmt.Printf("map1: %p %#v \n", &map1, map1) // panic: assignment to entry in nil map

	// channel
	channel := new(chan string)
	fmt.Printf("channel: %p %#v \n", &channel, channel) // channel: 0xc0000ae028 (*chan string)(0xc0000ae030) 
	channel <- "123" // Invalid operation: channel <- "123" (send to non-chan type *chan string) 
}
```

从运行结果可以看出，我们使用`new`函数分配内存后，只有数组在初始化后可以直接使用，`slice`、`map`、`chan`初始化后还是不能使用，会触发`panic`，这是因为`slice`、`map`、`chan`基本数据结构是一个`struct`，也就是说他里面的成员变量仍未进行初始化，所以他们初始化要使用`make`来进行，`make`会初始化他们的内部结构，我们下面一节细说。还是回到`struct`初始化的问题上，先看一个例子：

```go
type test struct {
	A *int64
}

func main(){
	t := new(test)
	*t.A = 10  // panic: runtime error: invalid memory address or nil pointer dereference
             // [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x10a89fd]
	fmt.Println(t.A)
}
```

从运行结果得出使用`new()`函数初始化结构体时，我们只是初始化了`struct`这个类型的，而它的成员变量是没有初始化的，所以初始化结构体不建议使用`new`函数，使用键值对进行初始化效果更佳。

其实 `new` 函数在日常工程代码中是比较少见的，因为它是可以被代替，使用`T{}`方式更加便捷方便。




## 初始化内置结构之`make`

在上一节我们说到了，`make`函数是专门支持 `slice`、`map`、`channel` 三种数据类型的内存创建，其官方定义如下：

```go
// The make built-in function allocates and initializes an object of type
// slice, map, or chan (only). Like new, the first argument is a type, not a
// value. Unlike new, make's return type is the same as the type of its
// argument, not a pointer to it. The specification of the result depends on
// the type:
//	Slice: The size specifies the length. The capacity of the slice is
//	equal to its length. A second integer argument may be provided to
//	specify a different capacity; it must be no smaller than the
//	length. For example, make([]int, 0, 10) allocates an underlying array
//	of size 10 and returns a slice of length 0 and capacity 10 that is
//	backed by this underlying array.
//	Map: An empty map is allocated with enough space to hold the
//	specified number of elements. The size may be omitted, in which case
//	a small starting size is allocated.
//	Channel: The channel's buffer is initialized with the specified
//	buffer capacity. If zero, or the size is omitted, the channel is
//	unbuffered.
func make(t Type, size ...IntegerType) Type
```

大概翻译最上面一段：`make`内置函数分配并初始化一个`slice`、`map`或`chan`类型的对象。像`new`函数一样，第一个参数是类型，而不是值。与`new`不同，`make`的返回类型与其参数的类型相同，而不是指向它的指针。结果的取决于传入的类型。

使用`make`初始化传入的类型也是不同的，具体可以这样区分：

```go
Func             Type T     res
make(T, n)       slice      slice of type T with length n and capacity n
make(T, n, m)    slice      slice of type T with length n and capacity m

make(T)          map        map of type T
make(T, n)       map        map of type T with initial space for approximately n elements

make(T)          channel    unbuffered channel of type T
make(T, n)       channel    buffered channel of type T, buffer size n
```

不同的类型初始化可以使用不同的姿势，主要区别主要是长度（len）和容量（cap）的指定，有的类型是没有容量这一说法，因此自然也就无法指定。如果确定长度和容量大小，能很好节省内存空间。

写个简单的示例：

```go
func main(){
	slice := make([]int64, 3, 5)
	fmt.Println(slice) // [0 0 0]
	map1 := make(map[int64]bool, 5)
	fmt.Println(map1) // map[]
	channel := make(chan int, 1)
	fmt.Println(channel) // 0xc000066070
}
```

这里有一个需要注意的点，就是`slice`在进行初始化时，默认会给零值，在开发中要注意这个问题，我就犯过这个错误，导致数据不一致。





## `new`和`make`区别总结

- `new`函数主要是为类型申请一片内存空间，返回执行内存的指针
- `make`函数能够分配并初始化类型所需的内存空间和结构，返回复合类型的本身。
- `make`函数仅支持 `channel`、`map`、`slice` 三种类型，其他类型不可以使用使用`make`。
- `new`函数在日常开发中使用是比较少的，可以被替代。
- `make`函数初始化`slice`会初始化零值，日常开发要注意这个问题。



## `make`函数底层实现

我还是比较好奇`make`底层实现是怎样的，所以执行汇编指令：`go tool compile -N -l -S file.go`，我们可以看到`make`函数初始化`slice`、`map`、`chan`分别调用的是`runtime.makeslice`、`runtime.makemap_small`、`runtime.makechan`这三个方法，因为不同类型底层数据结构不同，所以初始化方式也不同，我们只看一下`slice`的内部实现就好了，其他的交给大家自己去看，其实都是大同小异的。

```go
func makeslice(et *_type, len, cap int) unsafe.Pointer {
	mem, overflow := math.MulUintptr(et.size, uintptr(cap))
	if overflow || mem > maxAlloc || len < 0 || len > cap {
		// NOTE: Produce a 'len out of range' error instead of a
		// 'cap out of range' error when someone does make([]T, bignumber).
		// 'cap out of range' is true too, but since the cap is only being
		// supplied implicitly, saying len is clearer.
		// See golang.org/issue/4085.
		mem, overflow := math.MulUintptr(et.size, uintptr(len))
		if overflow || mem > maxAlloc || len < 0 {
			panicmakeslicelen()
		}
		panicmakeslicecap()
	}

	return mallocgc(mem, et, true)
}
```

这个函数功能其实也比较简单：

- 检查切片占用的内存空间是否溢出。
- 调用`mallocgc`在堆上申请一片连续的内存。

检查内存空间这里是根据切片容量进行计算的，根据当前切片元素的大小与切片容量的乘积得出当前内存空间的大小，检查溢出的条件有四个：

- 内存空间大小溢出了
- 申请的内存空间大于最大可分配的内存
- 传入的`len`小于`0`，`cap`的大小只小于`len`

`mallocgc`函数实现比较复杂，我暂时还没有看懂，不过也不是很重要，大家有兴趣可以自行学习。



## `new`函数底层实现

`new`函数底层主要是调用`runtime.newobject`：

```go
// implementation of new builtin
// compiler (both frontend and SSA backend) knows the signature
// of this function
func newobject(typ *_type) unsafe.Pointer {
	return mallocgc(typ.size, typ, true)
}
```

内部实现就是直接调用``mallocgc``函数去堆上申请内存，返回值是指针类型。



## 总结

今天这篇文章我们主要介绍了`make`和`new`的使用场景、以及其不同之处，其实他们都是用来分配内存的，只不过`make`函数为`slice`、`map`、`chan`这三种类型服务。日常开发中使用`make`初始化`slice`时要注意零值问题，否则又是一个`p0`事故。

**好啦，这篇文章到此结束啦，素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！我是`asong`，我们下期见。**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

推荐往期文章：

- [Go看源码必会知识之unsafe包](https://mp.weixin.qq.com/s/nPWvqaQiQ6Z0TaPoqg3t2Q)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/mzSCWI8C_ByIPbb07XYFTQ)
- [空结构体引发的大型打脸现场](https://mp.weixin.qq.com/s/dNeCIwmPei2jEWGF6AuWQw)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/wURQFRt2ISz66icW7jbHFw)
- [面试官：两个nil比较结果是什么？](https://mp.weixin.qq.com/s/Dt46eoEXXXZc2ymr67_LVQ)
- [面试官：你能用Go写段代码判断当前系统的存储方式吗?](https://mp.weixin.qq.com/s/ffEsTpO-tyNZFR5navAbdA)
- [如何平滑切换线上Elasticsearch索引](https://mp.weixin.qq.com/s/8VQxK_Xh-bkVoOdMZs4Ujw)