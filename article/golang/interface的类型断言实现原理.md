## 前言

> 哈喽，everyBody，我是`asong`，今天我们一起来探索一下`interface`的类型断言是如何实现的。我们通常使用`interface`有两种方式，一种是带方法的`interface`，一种是空的`interface`。因为`Go`中是没有泛型，所以我们可以用空的`interface{}`来作为一种伪泛型使用，当我们使用到空的`interface{}`作为入参或返回值时，就会使用到类型断言，来获取我们所需要的类型，所以平常我们会在代码中看到大量的类型断言使用，你就不好奇它是怎么实现的嘛？你就不好奇它的性能损耗是多少嘛？反正我很好奇，略～。



## 类型断言的基本使用

`Type Assertion`（断言）是用于`interface value`的一种操作，语法是`x.(T)`，`x`是`interface type`的表达式，而`T`是`asserted type`，被断言的类型。举个例子看一下基本使用：

```go
func main() {
	var demo interface{} = "Golang梦工厂"
	str := demo.(string)
	fmt.Printf("value: %v", str)
}
```

上面我们声明了一个接口对象`demo`，通过类型断言的方式断言一个接口对象`demo`是不是`nil`，并判断接口对象`demo`存储的值的类型是`T`，如果断言成功，就会返回值给`str`，如果断言失败，就会触发`panic`。这段代码加上如果这样写，就会触发`panic`：

```go
number := demo.(int64)
fmt.Printf("value： %v\n", number)
```

所以为了安全起见，我们还可以这样使用：

```go
func main() {
	var demo interface{} = "Golang梦工厂"
	number, ok := demo.(int64)
	if !ok {
		fmt.Printf("assert failed")
		return
	}
	fmt.Printf("value： %v\n", number)
}
运行结果：assert failed
```

这里使用的表达式是`t,ok:=i.(T)`，这个表达式也是可以断言一个接口对象`（i）`里不是` nil`，并且接口对象`（i）`存储的值的类型是 `T`，如果断言成功，就会返回其类型给` t`，并且此时 `ok` 的值 为` true`，表示断言成功。如果接口值的类型，并不是我们所断言的 `T`，就会断言失败，但和第一种表达式不同的是这个不会触发 `panic`，而是将 `ok` 的值设为` false `，表示断言失败，此时`t `为` T `的零值。所以推荐使用这种方式，可以保证代码的健壮性。

如果我们想要区分多种类型，可以使用`type switch`断言，使用这种方法就不需要我们按上面的方式去一个一个的进行类型断言了，更简单，更高效。上面的代码我们可以改成这样：

```go
func main() {
	var demo interface{} = "Golang梦工厂"

	switch demo.(type) {
	case nil:
		fmt.Printf("demo type is nil\n")
	case int64:
		fmt.Printf("demo type is int64\n")
	case bool:
		fmt.Printf("demo type is bool\n")
	case string:
		fmt.Printf("demo type is string\n")
	default:
		fmt.Printf("demo type unkonwn\n")
	}
}
```

`type switch`的一个典型应用是在`go.uber.org/zap`库中的`zap.Any()`方法，里面就用到了类型断言，把所有的类型的`case`都列举出来了，`default`分支使用的是`Reflect`，也就是当所有类型都不匹配时使用反射获取相应的值，具体大家可以去看一下源码。



## 类型断言实现源码剖析

非空接口和空接口都可以使用类型断言，我们分两种进行剖析。

### 空接口

我们先来写一段测试代码：

```go
type User struct {
	Name string
}

func main() {
	var u interface{} = &User{Name: "asong"}
	val, ok := u.(int)
	if !ok {
		fmt.Printf("%v\n", val)
	}
}
```

老样子，我们将上述代码转换成汇编代码看一下：

```go
go tool compile -S -N -l main.go > main.s4 2>&1
```

截取部分重要汇编代码如下：

```go
	0x002f 00047 (main.go:12)	XORPS	X0, X0
	0x0032 00050 (main.go:12)	MOVUPS	X0, ""..autotmp_8+136(SP)
	0x003a 00058 (main.go:12)	PCDATA	$2, $1
	0x003a 00058 (main.go:12)	PCDATA	$0, $0
	0x003a 00058 (main.go:12)	LEAQ	""..autotmp_8+136(SP), AX
	0x0042 00066 (main.go:12)	MOVQ	AX, ""..autotmp_7+96(SP)
	0x0047 00071 (main.go:12)	TESTB	AL, (AX)
	0x0049 00073 (main.go:12)	MOVQ	$5, ""..autotmp_8+144(SP)
	0x0055 00085 (main.go:12)	PCDATA	$2, $2
	0x0055 00085 (main.go:12)	LEAQ	go.string."asong"(SB), CX
	0x005c 00092 (main.go:12)	PCDATA	$2, $1
	0x005c 00092 (main.go:12)	MOVQ	CX, ""..autotmp_8+136(SP)
	0x0064 00100 (main.go:12)	MOVQ	AX, ""..autotmp_3+104(SP)
	0x0069 00105 (main.go:12)	PCDATA	$2, $2
	0x0069 00105 (main.go:12)	PCDATA	$0, $2
	0x0069 00105 (main.go:12)	LEAQ	type.*"".User(SB), CX
	0x0070 00112 (main.go:12)	PCDATA	$2, $1
	0x0070 00112 (main.go:12)	MOVQ	CX, "".u+120(SP)
	0x0075 00117 (main.go:12)	PCDATA	$2, $0
	0x0075 00117 (main.go:12)	MOVQ	AX, "".u+128(SP)
```

上面这段汇编代码的作用就是赋值给空接口，数据都存在栈上，因为空`interface{}`的结构是`eface`，所以就是组装了一个`eface`在内存中，内存布局如下：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-04-05%20%E4%B8%8B%E5%8D%884.45.44.png)

我们知道空接口的数据结构中只有两个字段，一个`_type`字段，一个`data`字段，从上图中，我们可以看出来，`eface`的`_type`存储在内存的`+120(SP)`处，`unsafe.Pointer`存在了`+128（SP）`处，现在我们知道了他是怎么存的了，接下来我们看一下空接口的类型断言汇编是怎么实现的：

```go
	0x007d 00125 (main.go:13)	PCDATA	$2, $1
	0x007d 00125 (main.go:13)	MOVQ	"".u+128(SP), AX
	0x0085 00133 (main.go:13)	PCDATA	$0, $0
	0x0085 00133 (main.go:13)	MOVQ	"".u+120(SP), CX
	0x008a 00138 (main.go:13)	PCDATA	$2, $3
	0x008a 00138 (main.go:13)	LEAQ	type.int(SB), DX
	0x0091 00145 (main.go:13)	PCDATA	$2, $1
	0x0091 00145 (main.go:13)	CMPQ	CX, DX
	0x0094 00148 (main.go:13)	JEQ	155
	0x0096 00150 (main.go:13)	JMP	395
	0x009b 00155 (main.go:13)	PCDATA	$2, $0
	0x009b 00155 (main.go:13)	MOVQ	(AX), AX
	0x009e 00158 (main.go:13)	MOVL	$1, CX
	0x00a3 00163 (main.go:13)	JMP	165
	0x00a5 00165 (main.go:13)	MOVQ	AX, ""..autotmp_4+80(SP)
	0x00aa 00170 (main.go:13)	MOVB	CL, ""..autotmp_5+71(SP)
	0x00ae 00174 (main.go:13)	MOVQ	""..autotmp_4+80(SP), AX
	0x00b3 00179 (main.go:13)	MOVQ	AX, "".val+72(SP)
	0x00b8 00184 (main.go:13)	MOVBLZX	""..autotmp_5+71(SP), AX
	0x00bd 00189 (main.go:13)	MOVB	AL, "".ok+70(SP)
	0x00c1 00193 (main.go:14)	CMPB	"".ok+70(SP), $0
```

从上面这段汇编我们可以看出来，空接口的类型断言是通过判断`eface`中的`_type`字段和比较的类型进行对比，相同就会去准备接下来的返回值，如果类型断言正确，经过中间临时变量的传递，最终`val`保存在内存中`+72(SP)`处。`ok`保存在内存`+70(SP)`处。

```go
	0x018b 00395 (main.go:15)	XORL	AX, AX
	0x018d 00397 (main.go:15)	XORL	CX, CX
	0x018f 00399 (main.go:13)	JMP	165
	0x0194 00404 (main.go:13)	NOP
```

如果断言失败，就会清空`AX`和`CX`寄存器，因为`AX`和`CX`中存的是`eface`结构体里面的字段。

**最后总结一下空接口类型断言实现流程：空接口类型断言实质是将`eface`中`_type`与要匹配的类型进行对比，匹配成功在内存中组装返回值，匹配失败直接清空寄存器，返回默认值。**



### 非空接口

老样子，还是先写一个例子，然后我们在看他的汇编实现：

```go
type Basic interface {
	GetName() string
	SetName(name string) error
}

type User struct {
	Name string
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) SetName(name string) error {
	u.Name = name
	return nil
}

func main() {
	var u Basic = &User{Name: "asong"}
	switch u.(type) {
	case *User:
		u1 := u.(*User)
		fmt.Println(u1.Name)
	default:
		fmt.Println("failed to match")
	}
}
```

使用汇编指令看一下他的汇编代码如下：

```go
	0x002f 00047 (main.go:26)	PCDATA	$2, $0
	0x002f 00047 (main.go:26)	PCDATA	$0, $1
	0x002f 00047 (main.go:26)	XORPS	X0, X0
	0x0032 00050 (main.go:26)	MOVUPS	X0, ""..autotmp_5+152(SP)
	0x003a 00058 (main.go:26)	PCDATA	$2, $1
	0x003a 00058 (main.go:26)	PCDATA	$0, $0
	0x003a 00058 (main.go:26)	LEAQ	""..autotmp_5+152(SP), AX
	0x0042 00066 (main.go:26)	MOVQ	AX, ""..autotmp_4+64(SP)
	0x0047 00071 (main.go:26)	TESTB	AL, (AX)
	0x0049 00073 (main.go:26)	MOVQ	$5, ""..autotmp_5+160(SP)
	0x0055 00085 (main.go:26)	PCDATA	$2, $2
	0x0055 00085 (main.go:26)	LEAQ	go.string."asong"(SB), CX
	0x005c 00092 (main.go:26)	PCDATA	$2, $1
	0x005c 00092 (main.go:26)	MOVQ	CX, ""..autotmp_5+152(SP)
	0x0064 00100 (main.go:26)	MOVQ	AX, ""..autotmp_2+72(SP)
	0x0069 00105 (main.go:26)	PCDATA	$2, $2
	0x0069 00105 (main.go:26)	PCDATA	$0, $2
	0x0069 00105 (main.go:26)	LEAQ	go.itab.*"".User,"".Basic(SB), CX
	0x0070 00112 (main.go:26)	PCDATA	$2, $1
	0x0070 00112 (main.go:26)	MOVQ	CX, "".u+104(SP)
	0x0075 00117 (main.go:26)	PCDATA	$2, $0
	0x0075 00117 (main.go:26)	MOVQ	AX, "".u+112(SP)
```

上面这段汇编代码作用就是赋值给非空接口的`iface`结构，组装了`iface`的内存布局，因为上面分析了非空接口的，这里就不细讲了，理解他的意思就好。接下来我们看一下他是如何进行类型断言的。

```go
	0x00df 00223 (main.go:29)	PCDATA	$2, $1
	0x00df 00223 (main.go:29)	PCDATA	$0, $2
	0x00df 00223 (main.go:29)	MOVQ	"".u+112(SP), AX
	0x00e4 00228 (main.go:29)	PCDATA	$0, $0
	0x00e4 00228 (main.go:29)	MOVQ	"".u+104(SP), CX
	0x00e9 00233 (main.go:29)	PCDATA	$2, $3
	0x00e9 00233 (main.go:29)	LEAQ	go.itab.*"".User,"".Basic(SB), DX
	0x00f0 00240 (main.go:29)	PCDATA	$2, $1
	0x00f0 00240 (main.go:29)	CMPQ	CX, DX
	0x00f3 00243 (main.go:29)	JEQ	250
	0x00f5 00245 (main.go:29)	JMP	583
	0x00fa 00250 (main.go:29)	MOVQ	AX, "".u1+56(SP)
```

上面代码我们可以看到调用`iface`结构中的`itab`字段，这里为什么这么调用呢？因为我们类型推断的是一个具体的类型，编译器会直接构造出`iface`，不会去调用已经在`runtime/iface.go`实现好的断言方法。上述代码中，先构造出` iface`，其中` *itab `存在内存 `+104(SP) `中，`unsafe.Pointer` 存在 `+112(SP)` 中。然后在类型推断的时候又重新构造了一遍 `*itab`，最后将新的 `*itab` 和前一次 `+104(SP)` 里的` *itab` 进行对比。

后面的赋值操作也就不再细说了，没有什么特别的。

这里还有一个要注意的问题，如果我们类型断言的是接口类型，那么我们在就会看到这样的汇编代码：

```go
// 代码修改
func main() {
	var u Basic = &User{Name: "asong"}
	v, ok := u.(Basic)
	if !ok {
		fmt.Printf("%v\n", v)
	}
}
	// 部分汇编代码
	0x008c 00140 (main.go:27)	MOVUPS	X0, ""..autotmp_4+168(SP)
	0x0094 00148 (main.go:27)	PCDATA	$2, $1
	0x0094 00148 (main.go:27)	MOVQ	"".u+128(SP), AX
	0x009c 00156 (main.go:27)	PCDATA	$0, $0
	0x009c 00156 (main.go:27)	MOVQ	"".u+120(SP), CX
	0x00a1 00161 (main.go:27)	PCDATA	$2, $4
	0x00a1 00161 (main.go:27)	LEAQ	type."".Basic(SB), DX
	0x00a8 00168 (main.go:27)	PCDATA	$2, $1
	0x00a8 00168 (main.go:27)	MOVQ	DX, (SP)
	0x00ac 00172 (main.go:27)	MOVQ	CX, 8(SP)
	0x00b1 00177 (main.go:27)	PCDATA	$2, $0
	0x00b1 00177 (main.go:27)	MOVQ	AX, 16(SP)
	0x00b6 00182 (main.go:27)	CALL	runtime.assertI2I2(SB)
```

我们可以看到，直接调用的是`runtime.assertI2I2()`方法进行类型断言，这个方法的实现代码如下：

```go
func assertI2I(inter *interfacetype, i iface) (r iface) {
	tab := i.tab
	if tab == nil {
		// explicit conversions require non-nil interface value.
		panic(&TypeAssertionError{nil, nil, &inter.typ, ""})
	}
	if tab.inter == inter {
		r.tab = tab
		r.data = i.data
		return
	}
	r.tab = getitab(inter, tab._type, false)
	r.data = i.data
	return
}
```

上述代码逻辑很简单，如果 `iface` 中的` itab.inter` 和第一个入参 `*interfacetype` 相同，说明类型相同，直接返回入参 `iface `的相同类型，布尔值为 `true`；如果` iface` 中的` itab.inter` 和第一个入参 `*interfacetype` 不相同，则重新根据 `*interfacetype` 和 `iface.tab` 去构造` tab`。构造的过程会查找` itabTable`。如果类型不匹配，或者不是属于同一个 `interface `类型，都会失败。`getitab() `方法第三个参数是 `canfail`，这里传入了` true`，表示构建 `*itab `允许失败，失败以后返回 `nil`。

**差异**：如果我们断言的类型是具体类型，编译器会直接构造出`iface`，不会去调用已经在`runtime/iface.go`实现好的断言方法。如果我们断言的类型是接口类型，将会去调用相应的断言方法进行判断。

**小结**：**非空接口类型断言的实质是 iface 中 `*itab` 的对比。`*itab` 匹配成功会在内存中组装返回值。匹配失败直接清空寄存器，返回默认值。**



## 类型断言的性能损耗

前面我们已经分析了断言的底层原理，下面我们来看一下不同场景下进行断言的代价。

针对不同的场景可以写出测试文件如下（截取了部分代码，全部代码获取[戳这里](https://github.com/asong2020/Golang_Dream/tree/master/code_demo/assert_test)）: 

```go
var dst int64

// 空接口类型直接类型断言具体的类型
func Benchmark_efaceToType(b *testing.B) {
	b.Run("efaceToType", func(b *testing.B) {
		var ebread interface{} = int64(666)
		for i := 0; i < b.N; i++ {
			dst = ebread.(int64)
		}
	})
}

// 空接口类型使用TypeSwitch 只有部分类型
func Benchmark_efaceWithSwitchOnlyIntType(b *testing.B) {
	b.Run("efaceWithSwitchOnlyIntType", func(b *testing.B) {
		var ebread interface{} = 666
		for i := 0; i < b.N; i++ {
			OnlyInt(ebread)
		}
	})
}

// 空接口类型使用TypeSwitch 所有类型
func Benchmark_efaceWithSwitchAllType(b *testing.B) {
	b.Run("efaceWithSwitchAllType", func(b *testing.B) {
		var ebread interface{} = 666
		for i := 0; i < b.N; i++ {
			Any(ebread)
		}
	})
}

//直接使用类型转换
func Benchmark_TypeConversion(b *testing.B) {
	b.Run("typeConversion", func(b *testing.B) {
		var ebread int32 = 666

		for i := 0; i < b.N; i++ {
			dst = int64(ebread)
		}
	})
}

// 非空接口类型判断一个类型是否实现了该接口 两个方法
func Benchmark_ifaceToType(b *testing.B) {
	b.Run("ifaceToType", func(b *testing.B) {
		var iface Basic = &User{}
		for i := 0; i < b.N; i++ {
			iface.GetName()
			iface.SetName("1")
		}
	})
}

// 非空接口类型判断一个类型是否实现了该接口 12个方法
func Benchmark_ifaceToTypeWithMoreMethod(b *testing.B) {
	b.Run("ifaceToTypeWithMoreMethod", func(b *testing.B) {
		var iface MoreMethod = &More{}
		for i := 0; i < b.N; i++ {
			iface.Get()
			iface.Set()
			iface.One()
			iface.Two()
			iface.Three()
			iface.Four()
			iface.Five()
			iface.Six()
			iface.Seven()
			iface.Eight()
			iface.Nine()
			iface.Ten()
		}
	})
}

// 直接调用方法
func Benchmark_DirectlyUseMethod(b *testing.B) {
	b.Run("directlyUseMethod", func(b *testing.B) {
		m := &More{
			Name: "asong",
		}
		m.Get()
	})
}
```

运行结果：

```go
goos: darwin
goarch: amd64
pkg: asong.cloud/Golang_Dream/code_demo/assert_test
Benchmark_efaceToType/efaceToType-16            1000000000               0.507 ns/op
Benchmark_efaceWithSwitchOnlyIntType/efaceWithSwitchOnlyIntType-16              384958000                3.00 ns/op
Benchmark_efaceWithSwitchAllType/efaceWithSwitchAllType-16                      351172759                3.33 ns/op
Benchmark_TypeConversion/typeConversion-16                                      1000000000               0.473 ns/op
Benchmark_ifaceToType/ifaceToType-16                                            355683139                3.38 ns/op
Benchmark_ifaceToTypeWithMoreMethod/ifaceToTypeWithMoreMethod-16                85421563                12.8 ns/op
Benchmark_DirectlyUseMethod/directlyUseMethod-16                                1000000000               0.000000 ns/op
PASS
ok      asong.cloud/Golang_Dream/code_demo/assert_test  7.797s
```

从结果我们可以分析一下：

- 空接口类型的类型断言代价并不高，与直接类型转换几乎没有性能差异
- 空接口类型使用`type switch`进行类型断言时，随着`case`的增多性能会直线下降
- 非空接口类型进行类型断言时，随着接口中方法的增多，性能会直线下降
- 直接进行方法调用要比非接口类型进行类型断言要高效很多

好啦，现在我们也知道怎样使用类型断言能提高性能啦，又可以和同事吹水一手啦。



## 总结

好啦，本文到这里就已经接近尾声了，在最后做一个小小的总结：

- 空接口类型断言实现流程：空接口类型断言实质是将`eface`中`_type`与要匹配的类型进行对比，匹配成功在内存中组装返回值，匹配失败直接清空寄存器，返回默认值。
- 非空接口类型断言的实质是 iface 中 `*itab` 的对比。`*itab` 匹配成功会在内存中组装返回值。匹配失败直接清空寄存器，返回默认值

- 泛型是在编译期做的事情，使用类型断言会消耗一点性能，类型断言使用方式不同，带来的性能损耗也不同，具体请看上面的章节。

**文中代码已上传`github`：https://github.com/asong2020/Golang_Dream/tree/master/code_demo/assert_test，欢迎`star`**

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
- [空结构体引发的大型打脸现场](https://mp.weixin.qq.com/s/dNeCIwmPei2jEWGF6AuWQw)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/wURQFRt2ISz66icW7jbHFw)
- [十张动图带你搞懂排序算法(附go实现代码)](https://mp.weixin.qq.com/s/rZBsoKuS-ORvV3kML39jKw)
- [go参数传递类型](https://mp.weixin.qq.com/s/JHbFh2GhoKewlemq7iI59Q)
- [手把手教姐姐写消息队列](https://mp.weixin.qq.com/s/0MykGst1e2pgnXXUjojvhQ)
- [常见面试题之缓存雪崩、缓存穿透、缓存击穿](https://mp.weixin.qq.com/s?__biz=MzIzMDU0MTA3Nw==&mid=2247483988&idx=1&sn=3bd52650907867d65f1c4d5c3cff8f13&chksm=e8b0902edfc71938f7d7a29246d7278ac48e6c104ba27c684e12e840892252b0823de94b94c1&token=1558933779&lang=zh_CN#rd)
- [详解Context包，看这一篇就够了！！！](https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg)
- [面试官：你能用Go写段代码判断当前系统的存储方式吗?](https://mp.weixin.qq.com/s/ffEsTpO-tyNZFR5navAbdA)
- [如何平滑切换线上Elasticsearch索引](https://mp.weixin.qq.com/s/8VQxK_Xh-bkVoOdMZs4Ujw)

