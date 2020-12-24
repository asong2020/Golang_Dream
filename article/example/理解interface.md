## 前言

> 我想，对于各位使用面向对象编程的程序员来说，"接口"这个名词一定不陌生，比如java中的接口以及c++中的虚基类都是接口的实现。但是`golang`中的接口概念确与其他语言不同，有它自己的特点，下面我们就来一起解密。



## 定义

Go 语言中的接口是一组方法的签名，它是 Go 语言的重要组成部分。简单的说，interface是一组method签名的组合，我们通过interface来定义对象的一组行为。**interface 是一种类型**，定义如下：

```go
type Person interface {
    Eat(food string) 
}
```

它的定义可以看出来用了 type 关键字，更准确的说 interface 是一种**具有一组方法的类型**，这些方法定义了 interface 的行为。`golang`接口定义不能包含变量，但是允许不带任何方法，这种类型的接口叫`empty interface`。

**如果一个类型实现了一个`interface`中所有方法，我们就可以说该类型实现了该`interface`，所以我们我们的所有类型都实现了`empty interface`，因为任何一种类型至少实现了0个方法。并且`go`中并不像`java`中那样需要显式关键字来实现`interface`，只需要实现`interface`包含的方法即可。**



### 实现接口

这里先拿`java`语言来举例，在`java`中，我们要实现一个`interface`需要这样声明：

```java
public class MyWriter implments io.Writer{}
```

这就意味着对于接口的实现都需要显示声明，在代码编写方面有依赖限制，同时需要处理包的依赖，而在`Go`语言中实现接口就是隐式的，举例说明：

```go
type error interface {
	Error() string
}
type RPCError struct {
	Code    int64
	Message string
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("%s, code=%d", e.Message, e.Code)
}
```

上面的代码，并没有`error`接口的影子，我们只需要实现`Error() string`方法就实现了`error`接口。在`Go`中，实现接口的所有方法就隐式地实现了接口。我们使用上述 `RPCError` 结构体时并不关心它实现了哪些接口，Go 语言只会在传递参数、返回参数以及变量赋值时才会对某个类型是否实现接口进行检查。

`Go`语言的这种写法很方便，不用引入包依赖。但是`interface`底层实现的时候会动态检测也会引入一些问题：

- 性能下降。使用interface作为函数参数，runtime 的时候会动态的确定行为。使用具体类型则会在编译期就确定类型。
- 不能清楚的看出struct实现了哪些接口，需要借助ide或其它工具。

## 两种接口

这里大多数刚入门的同学肯定会有疑问，怎么会有两种接口，因为`Go`语言中接口会有两种表现形式，使用`runtime.iface`表示第一种接口，也就是我们上面实现的这种，接口中定义方法；使用`runtime.eface`表示第二种不包含任何方法的接口，第二种在我们日常开发中经常使用到，所以在实现时使用了特殊的类型。从编译角度来看，golang并不支持泛型编程。但还是可以用`interface{}`  来替换参数，而实现泛型。

### interface内部结构

Go 语言根据接口类型是否包含一组方法将接口类型分成了两类：

- 使用 [`runtime.iface`](https://draveness.me/golang/tree/runtime.iface) 结构体表示包含方法的接口
- 使用 [`runtime.eface`](https://draveness.me/golang/tree/runtime.eface) 结构体表示不包含任何方法的 `interface{}` 类型；

`runtime.iface`结构体在`Go`语言中的定义是这样的：

```go
type eface struct { // 16 字节
	_type *_type
	data  unsafe.Pointer
}
```

这里只包含指向底层数据和类型的两个指针，从这个`type`我们也可以推断出Go语言的任意类型都可以转换成`interface`。

另一个用于表示接口的结构体是 [`runtime.iface`](https://draveness.me/golang/tree/runtime.iface)，这个结构体中有指向原始数据的指针 `data`，不过更重要的是 [`runtime.itab`](https://draveness.me/golang/tree/runtime.itab) 类型的 `tab` 字段。

```go
type iface struct { // 16 字节
	tab  *itab
	data unsafe.Pointer
}
```

下面我们一起看看`interface`中这两个类型：

- `runtime_type`

`runtime_type`是 Go 语言类型的运行时表示。下面是运行时包中的结构体，其中包含了很多类型的元信息，例如：类型的大小、哈希、对齐以及种类等。

```go
type _type struct {
	size       uintptr
	ptrdata    uintptr
	hash       uint32
	tflag      tflag
	align      uint8
	fieldAlign uint8
	kind       uint8
	equal      func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata     *byte
	str        nameOff
	ptrToThis  typeOff
}
```

这里我只对几个比较重要的字段进行讲解：

- `size` 字段存储了类型占用的内存空间，为内存空间的分配提供信息；
- `hash` 字段能够帮助我们快速确定类型是否相等；
- `equal` 字段用于判断当前类型的多个对象是否相等，该字段是为了减少 Go 语言二进制包大小从 `typeAlg` 结构体中迁移过来的)；



- `runtime_itab`

`runtime.itab`结构体是接口类型的核心组成部分，每一个 `runtime.itab` 都占 32 字节，我们可以将其看成接口类型和具体类型的组合，它们分别用 `inter` 和 `_type` 两个字段表示：

```go
type itab struct { // 32 字节
	inter *interfacetype
	_type *_type
	hash  uint32
	_     [4]byte
	fun   [1]uintptr
}
```

`inter`和`_type`是用于表示类型的字段，`hash`是对`_type.hash`的拷贝，当我们想将 `interface` 类型转换成具体类型时，可以使用该字段快速判断目标类型和具体类型 `runtime._type`是否一致，`fun`是一个动态大小的数组，它是一个用于动态派发的虚函数表，存储了一组函数指针。虽然该变量被声明成大小固定的数组，但是在使用时会通过原始指针获取其中的数据，所以 `fun` 数组中保存的元素数量是不确定的；

内部结构就做一个简单介绍吧，有兴趣的同学可以自行深入学习。



### 空的interface（`runtime.eface`）

前文已经介绍了什么是空的`interface`，下面我们来看一看空的`interface`如何使用。定义函数入参如下：

```go
func doSomething(v interface{}){    
}
```

这个函数的入参是`interface`类型，要注意的是，`interface`类型不是任意类型，他与C语言中的`void *`不同，如果我们将类型转换成了 `interface{}` 类型，变量在运行期间的类型也会发生变化，获取变量类型时会得到 `interface{}`，之所以函数可以接受任何类型是在 go 执行时传递到函数的任何类型都被自动转换成 `interface{}`。

那么我们可以才来一个猜想，既然空的 interface 可以接受任何类型的参数，那么一个 `interface{}`类型的 slice 是不是就可以接受任何类型的 slice ？下面我们就来尝试一下：

```go

import (
	"fmt"
)

func printStr(str []interface{}) {
	for _, val := range str {
		fmt.Println(val)
	}
}

func main(){
	names := []string{"stanley", "david", "oscar"}
	printStr(names)
}
```

运行上面代码，会出现如下错误：`./main.go:15:10: cannot use names (type []string) as type []interface {} in argument to printStr`。

这里我也是很疑惑，为什么`Go`没有帮助我们自动把`slice`转换成`interface`类型的`slice`，之前做项目就想这么用，结果失败了。后来我终于找到了[答案](https://github.com/golang/go/wiki/InterfaceSlice)，有兴趣的可以看看原文，这里简单总结一下：`interface`会占用两个字长的存储空间，一个是自身的 methods 数据，一个是指向其存储值的指针，也就是 interface 变量存储的值，因而 slice []interface{} 其长度是固定的`N*2`，但是 []T 的长度是`N*sizeof(T)`，两种 slice 实际存储值的大小是有区别的。

既然这种方法行不通，那可以怎样解决呢？我们可以直接使用元素类型是interface的切片。

```go
var dataSlice []int = foo()
var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
for i, d := range dataSlice {
	interfaceSlice[i] = d
}
```



### 非空`interface`

`Go`语言实现接口时，既可以结构体类型的方法也可以是使用指针类型的方法。`Go`语言中并没有严格规定实现者的方法是值类型还是指针，那我们猜想一下，如果同时使用值类型和指针类型方法实现接口，会有什么问题吗？

先看这样一个例子：

```go
package main

import (
	"fmt"
)

type Person interface {
	GetAge () int
	SetAge (int)
}


type Man struct {
	Name string
	Age int
}

func(s Man) GetAge()int {
return s.Age
}

func(s *Man) SetAge(age int) {
	s.Age = age
}


func f(p Person){
	p.SetAge(10)
	fmt.Println(p.GetAge())
}

func main() {
	p := Man{}
	f(&p) 
}
```

看上面的代码，大家对`f(&p)`这里的入参是否会有疑问呢？如果不取地址，直接传过去会怎么样？试了一下，编译错误如下：`./main.go:34:3: cannot use p (type Man) as type Person in argument to f: Man does not implement Person (SetAge method has pointer receiver)`。透过注释我们可以看到，因为`SetAge`方法的`receiver`是指针类型，那么传递给`f`的是`P`的一份拷贝，在进行`p`的拷贝到`person`的转换时，`p`的拷贝是不满足`SetAge`方法的`receiver`是个指针类型，这也正说明一个问题**go中函数都是按值传递**。

上面的例子是因为发生了值传递才会导致出现这个问题。实际上不管接收者类型是值类型还是指针类型，都可以通过值类型或指针类型调用，这里面实际上通过语法糖起作用的。实现了接收者是值类型的方法，相当于自动实现了接收者是指针类型的方法；而实现了接收者是指针类型的方法，不会自动生成对应接收者是值类型的方法。

举个例子：

```go
type Animal interface {
	Walk()
	Eat()
}


type Dog struct {
	Name string
}

func (d *Dog)Walk()  {
	fmt.Println("go")
}

func (d *Dog)Eat()  {
	fmt.Println("eat shit")
}

func main() {
	var d Animal = &Dog{"nene"}
	d.Eat()
	d.Walk()
}
```

上面定义了一个接口`Animal`，接口定义了两个函数：

```go
Walk()
Eat()
```

接着定义了一个结构体`Dog`，他实现了两个方法，一个是值接受者，一个是指针接收者。我们通过接口类型的变量调用了定义的两个函数是没有问题的，如果我们改成这样呢：

```go
func main() {
	var d Animal = Dog{"nene"}
	d.Eat()
	d.Walk()
}
```

这样直接就会报错，我们只改了一部分，第一次将`&Dog{"nene"}`赋值给了`d`；第二次则将`Dog{"nene"}`赋值给了`d`。第二次报错是因为，`d`没有实现`Animal`。这正解释了上面的结论，所以，当实现了一个接收者是值类型的方法，就可以自动生成一个接收者是对应指针类型的方法，因为两者都不会影响接收者。但是，当实现了一个接收者是指针类型的方法，如果此时自动生成一个接收者是值类型的方法，原本期望对接收者的改变（通过指针实现），现在无法实现，因为值类型会产生一个拷贝，不会真正影响调用者。

总结一句话就是：**如果实现了接收者是值类型的方法，会隐含地也实现了接收者是指针类型的方法。**



## 类型断言

一个`interface`被多种类型实现时，有时候我们需要区分`interface`的变量究竟存储哪种类型的值，`go`可以使用`comma,ok`的形式做区分 `value, ok := em.(T)`：**em 是 interface 类型的变量，T代表要断言的类型，value 是 interface 变量存储的值，ok 是 bool 类型表示是否为该断言的类型 T**。总结出来语法如下：

```go
<目标类型的值>，<布尔参数> := <表达式>.( 目标类型 ) // 安全类型断言
<目标类型的值> := <表达式>.( 目标类型 )　　//非安全类型断言
```

看个简单的例子：

```go
type Dog struct {
	Name string
}

func main() {
	var d interface{} = new(Dog)
	d1,ok := d.(Dog)
	if !ok{
		return
	}
	fmt.Println(d1)
}
```

这种就属于安全类型断言，更适合在线上代码使用，如果使用非安全类型断言会怎么样呢？

```go
type Dog struct {
	Name string
}

func main() {
	var d interface{} = new(Dog)
	d1 := d.(Dog)
	fmt.Println(d1)
}
```

这样就会发生错误如下：

```go
panic: interface conversion: interface {} is *main.Dog, not main.Dog
```

断言失败。这里直接发生了 `panic`，所以不建议线上代码使用。

看过`fmt`源码包的同学应该知道，`fmt.println`内部就是使用到了类型断言，有兴趣的同学可以自行学习。





## 问题

上面介绍了`interface`的基本使用方法及可能会遇到的一些问题，下面出三个题，看看你们真的掌握了吗？



### 问题一

下面代码，哪一行存在编译错误？（多选）

```go
type Student struct {
}

func Set(x interface{}) {
}

func Get(x *interface{}) {
}

func main() {
	s := Student{}
	p := &s
	// A B C D
	Set(s)
	Get(s)
	Set(p)
	Get(p)
}
```

答案：B、D；解析：我们上文提到过，`interface`是所有`go`类型的父类，所以`Get`方法只能接口`*interface{}`类型的参数，其他任何类型都不可以。

### 问题二

这段代码的运行结果是什么？

```go
func PrintInterface(val interface{}) {
	if val == nil {
		fmt.Println("this is empty interface")
		return
	}
	fmt.Println("this is non-empty interface")
}
func main() {
	var pointer *string = nil
	PrintInterface(pointer)
}
```

答案：`this is non-empty interface`。解析：这里的`interface{}`是空接口类型，他的结构如下:

```go
type eface struct { // 16 字节
	_type *_type
	data  unsafe.Pointer
}
```

所以在调用函数`PrintInterface`时发生了**隐式的类型转换**，除了向方法传入参数之外，变量的赋值也会触发隐式类型转换。在类型转换时，`*string`类型会转换成`interface`类型，发生值拷贝，所以`eface struct{}`是不为`nil`，不过`data`指针指向的`poniter`为`nil`。

### 问题三

这段代码的运行结果是什么？

```go

type Animal interface {
	Walk()
}

type Dog struct{}

func (d *Dog) Walk() {
	fmt.Println("walk")
}

func NewAnimal() Animal {
	var d *Dog
	return d
}

func main() {
	if NewAnimal() == nil {
		fmt.Println("this is empty interface")
	} else {
		fmt.Println("this is non-empty interface")
	}
}
```

答案：`this is non-empty interface`. 解析：这里的`interface`是非空接口`iface`，他的结构如下：

```go
type iface struct { // 16 字节
	tab  *itab
	data unsafe.Pointer
}
```

`d`是一个指向nil的空指针，但是最后`return d` 会触发`匿名变量 Animal = p`值拷贝动作，所以最后`NewAnimal()`返回给上层的是一个`Animal interface{}`类型，也就是一个`iface struct{}`类型。 `p`为nil，只是`iface`中的data 为nil而已。 但是`iface struct{}`本身并不为nil.





## 总结

`interface`在我们日常开发中使用还是比较多，所以学好它还是很必要，希望这篇文章能让你对`Go`语言的接口有一个新的认识，这一篇到这里结束啦，我们下期见～～～。

**素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！**

**建了一个Golang交流群，欢迎大家的加入，第一时间观看优质文章，不容错过哦（公众号获取）**

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**翻译了一份Machinery中文文档，会定期进行维护，有需要的小伙伴们后台回复[machinery]即可获取。**

**我是asong，一名普普通通的程序猿，让gi我一起慢慢变强吧。我自己建了一个`golang`交流群，有需要的小伙伴加我`vx`,我拉你入群。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [machinery-go异步任务队列](https://mp.weixin.qq.com/s/4QG69Qh1q7_i0lJdxKXWyg)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/wURQFRt2ISz66icW7jbHFw)
- [十张动图带你搞懂排序算法(附go实现代码)](https://mp.weixin.qq.com/s/rZBsoKuS-ORvV3kML39jKw)
- [Go语言相关书籍推荐（从入门到放弃）](https://mp.weixin.qq.com/s/PaTPwRjG5dFMnOSbOlKcQA)
- [go参数传递类型](https://mp.weixin.qq.com/s/JHbFh2GhoKewlemq7iI59Q)
- [手把手教姐姐写消息队列](https://mp.weixin.qq.com/s/0MykGst1e2pgnXXUjojvhQ)
- [常见面试题之缓存雪崩、缓存穿透、缓存击穿](https://mp.weixin.qq.com/s?__biz=MzIzMDU0MTA3Nw==&mid=2247483988&idx=1&sn=3bd52650907867d65f1c4d5c3cff8f13&chksm=e8b0902edfc71938f7d7a29246d7278ac48e6c104ba27c684e12e840892252b0823de94b94c1&token=1558933779&lang=zh_CN#rd)
- [详解Context包，看这一篇就够了！！！](https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg)
- [go-ElasticSearch入门看这一篇就够了(一)](https://mp.weixin.qq.com/s/mV2hnfctQuRLRKpPPT9XRw)
- [面试官：go中for-range使用过吗？这几个问题你能解释一下原因吗](https://mp.weixin.qq.com/s/G7z80u83LTgLyfHgzgrd9g)
- [学会wire依赖注入、cron定时任务其实就这么简单！](https://mp.weixin.qq.com/s/qmbCmwZGmqKIZDlNs_a3Vw)
- [听说你还不会jwt和swagger-饭我都不吃了带着实践项目我就来了](https://mp.weixin.qq.com/s/z-PGZE84STccvfkf8ehTgA)

  

