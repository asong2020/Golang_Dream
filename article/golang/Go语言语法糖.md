## 前言

> 哈喽，大家好，我是`asong`。
>
> 每门语言都有自己的语法糖，像`java`的语法糖就有方法变长参数、拆箱与装箱、枚举、`for-each`等等，`Go`语言也不例外，其也有自己的语法糖，掌握这些语法糖可以助我们提高开发的效率，所以本文就来介绍一些`Go`语言的语法糖，总结的可能不能全，欢迎评论区补充。



## 可变长参数

`Go`语言允许一个函数把任意数量的值作为参数，`Go`语言内置了**...**操作符，在函数的最后一个形参才能使用**...**操作符，使用它必须注意如下事项：

- 可变长参数必须在函数列表的最后一个；
- 把可变长参数当切片来解析，可变长参数没有没有值时就是`nil`切片
- 可变长参数的类型必须相同

```go
func test(a int, b ...int){
  return
}
```

既然我们的函数可以接收可变长参数，那么我们在传参的时候也可以传递切片使用**...**进行解包转换为参数列表，`append`方法就是最好的例子：

```go
var sl []int
sl = append(sl, 1)
sl = append(sl, sl...)
```

append方法定义如下：

```go
//	slice = append(slice, elem1, elem2)
//	slice = append(slice, anotherSlice...)
func append(slice []Type, elems ...Type) []Type
```



## 声明不定长数组

数组是有固定长度的，我们在声明数组时一定要声明长度，因为数组在编译时就要确认好其长度，但是有些时候对于想偷懒的我，就是不想写数组长度，有没有办法让他自己算呢？当然有，使用**...**操作符声明数组时，你只管填充元素值，其他的交给编译器自己去搞就好了；

```go
a := [...]int{1, 3, 5} // 数组长度是3，等同于 a := [3]{1, 3, 5}
```

有时我们想声明一个大数组，但是某些`index`想设置特别的值也可以使用**...**操作符搞定：

```go
a := [...]int{1: 20, 999: 10} // 数组长度是100, 下标1的元素值是20，下标999的元素值是10，其他元素值都是0
```



## `init`函数

`Go`语言提供了先于`main`函数执行的`init`函数，初始化每个包后会自动执行`init`函数，每个包中可以有多个`init`函数，每个包中的源文件中也可以有多个`init`函数，加载顺序如下：

> 从当前包开始，如果当前包包含多个依赖包，则先初始化依赖包，层层递归初始化各个包，在每一个包中，按照源文件的字典序从前往后执行，每一个源文件中，优先初始化常量、变量，最后初始化`init`函数，当出现多个`init`函数时，则按照顺序从前往后依次执行，每一个包完成加载后，递归返回，最后在初始化当前包！

`init`函数实现了`sync.Once`，无论包被导入多少次，`init`函数只会被执行一次，所以使用`init`可以应用在服务注册、中间件初始化、实现单例模式等等，比如我们经常使用的`pprof`工具，他就使用到了`init`函数，在`init`函数里面进行路由注册：

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



## 忽略导包

Go语言在设计师有代码洁癖，在设计上尽可能避免代码滥用，所以`Go`语言的导包必须要使用，如果导包了但是没有使用的话就会产生编译错误，但有些场景我们会遇到只想导包，但是不使用的情况，比如上文提到的`init`函数，我们只想初始化包里的`init`函数，但是不会使用包内的任何方法，这时就可以使用  **_**   操作符号重命名导入一个不使用的包：

```go
import _ "github.com/asong"
```



## 忽略字段

在我们日常开发中，一般都是在屎上上堆屎，遇到可以用的方法就直接复用了，但是这个方法的返回值我们并不一定都使用，还要绞尽脑汁的给他想一个命名，有没有办法可以不处理不要的返回值呢？当然有，还是 **_** 操作符，将不需要的值赋给空标识符：

```go
_, ok := test(a, b int)
```



## json序列化忽略某个字段

大多数业务场景我们都会对`struct`做序列化操作，但有些时候我们想要`json`里面的某些字段不参加序列化，**-**操作符可以帮我们处理，`Go`语言的结构体提供标签功能，在结构体标签中使用 **-** 操作符就可以对不需要序列化的字段做特殊处理，使用如下：

```go
type Person struct{
  name string `json:"-"`
  age string `json: "age"`
}
```



## json序列化忽略空值字段

我们使用`json.Marshal`进行序列化是不会忽略`struct`中的空值，默认输出字段的类型零值（`string`类型零值是""，对象类型的零值是`nil`...），如果我们想在序列化是忽略掉这些没有值的字段时，可以在结构体标签中中添加`omitempty` tag：

```go
type User struct {
	Name  string   `json:"name"`
	Email string   `json:"email,omitempty"`
  Age int        `json: "age"`
}

func test() {
	u1 := User{
		Name: "asong",
	}
	b, err := json.Marshal(u1)
	if err != nil {
		fmt.Printf("json.Marshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b)
}
```

运行结果：

```go
str:{"name":"asong","Age":0}
```

`Age`字段我们没有添加`omitempty` tag在`json`序列化结果就是带空值的，`email`字段就被忽略掉了；



## 短变量声明

每次使用变量时都要先进行函数声明，对于我这种懒人来说是真的不想写，因为写`python`写惯了，那么在`Go`语言是不是也可以不进行变量声明直接使用呢？我们可以使用 **name := expression** 的语法形式来声明和初始化局部变量，相比于使用`var`声明的方式可以减少声明的步骤：

```go
var a int = 10
等用于
a := 10
```

使用短变量声明时有两个注释事项：

- 短变量声明只能在函数内使用，不能用于初始化全局变量
- 短变量声明代表引入一个新的变量，不能在同一作用域重复声明变量
- 多变量声明中如果其中一个变量是新变量，那么可以使用短变量声明，否则不可重复声明变量；



## 类型断言

我们通常都会使用`interface`，一种是带方法的`interface`，一种是空的`interface`，`Go1.18`之前是没有泛型的，所以我们可以用空的`interface{}`来作为一种伪泛型使用，当我们使用到空的`interface{}`作为入参或返回值时，就会使用到类型断言，来获取我们所需要的类型，在Go语言中类型断言的语法格式如下：

```go
value, ok := x.(T)
or
value := x.(T)
```

x是`interface`类型，T是具体的类型，方式一是安全的断言，方式二断言失败会触发panic；这里类型断言需要区分`x`的类型，如果`x`是空接口类型：

**空接口类型断言实质是将`eface`中`_type`与要匹配的类型进行对比，匹配成功在内存中组装返回值，匹配失败直接清空寄存器，返回默认值。**

如果`x`是非空接口类型：

**非空接口类型断言的实质是 iface 中 `*itab` 的对比。`*itab` 匹配成功会在内存中组装返回值。匹配失败直接清空寄存器，返回默认值。**

具体源码剖析可以看这篇文章：[源码剖析类型断言是如何实现的！附性能损耗测试](https://mp.weixin.qq.com/s/JqjxV8Jej3t89KdvsvGZhg_)



## 切片循环

切片/数组是我们经常使用的操作，在`Go`语言中提供了`for range`语法来快速迭代对象，数组、切片、字符串、map、channel等等都可以进行遍历，总结起来总共有三种方式：

```go
// 方式一：只遍历不关心数据，适用于切片、数组、字符串、map、channel
for range T {}

// 方式二：遍历获取索引或数组，切片，数组、字符串就是索引，map就是key，channel就是数据
for key := range T{}

// 方式三：遍历获取索引和数据，适用于切片、数组、字符串，第一个参数就是索引，第二个参数就是对应的元素值，map 第一个参数就是key，第二个参数就是对应的值；
for key, value := range T{}
```



## 判断map的key是否存在

Go语言提供语法 `value, ok := m[key]`来判断`map`中的`key`是否存在，如果存在就会返回key所对应的值，不存在就会返回空值：

```go
import "fmt"

func main() {
    dict := map[string]int{"asong": 1}
    if value, ok := dict["asong"]; ok {
        fmt.Printf(value)
    } else {
      fmt.Println("key:asong不存在")
    }
}
```



## select控制结构

`Go`语言提供了`select`关键字，`select`配合`channel`能够让`Goroutine`同时等待多个`channel`读或者写，在`channel`状态未改变之前，`select`会一直阻塞当前线程或`Goroutine`。先看一个例子：

```go
func fibonacci(ch chan int, done chan struct{}) {
 x, y := 0, 1
 for {
  select {
  case ch <- x:
   x, y = y, x+y
  case <-done:
   fmt.Println("over")
   return
  }
 }
}
func main() {
 ch := make(chan int)
 done := make(chan struct{})
 go func() {
  for i := 0; i < 10; i++ {
   fmt.Println(<-ch)
  }
  done <- struct{}{}
 }()
 fibonacci(ch, done)
}
```

`select`与`switch`具有相似的控制结构，与`switch`不同的是，`select`中的`case`中的表达式必须是`channel`的收发操作，当`select`中的两个`case`同时被触发时，会随机执行其中的一个。为什么是随机执行的呢？随机的引入就是为了避免饥饿问题的发生，如果我们每次都是按照顺序依次执行的，若两个`case`一直都是满足条件的，那么后面的`case`永远都不会执行。

上面例子中的`select`用法是阻塞式的收发操作，直到有一个`channel`发生状态改变。我们也可以在`select`中使用`default`语句，那么`select`语句在执行时会遇到这两种情况：

- 当存在可以收发的`Channel`时，直接处理该`Channel` 对应的 `case`；
- 当不存在可以收发的`Channel` 时，执行 `default` 中的语句；

**注意：`nil channel`上的操作会一直被阻塞，如果没有`default case`,只有`nil channel`的`select`会一直被阻塞。**



## 总结

本文介绍了`Go`语言中的一些开发技巧，也就是`Go`语言的语法糖，掌握好这些可以提高我们的开发效率，你都学会了吗？

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**


![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)
