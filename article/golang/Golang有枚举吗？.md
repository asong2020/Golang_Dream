## 前言

> 哈喽，大家好，我是`asong`。枚举是一种很重要的数据类型，在`java`、`C`语言等主流编程语言中都支持了枚举类型，但是在Go语言中却没有枚举类型，那有什么替代方案吗？ 本文我们来聊一聊这个事情；



## 为什么要有枚举

我们以`java`语言为例子，在`JDK1.5`之前没有枚举类型，我们通常会使用`int`常量来表示枚举，一般使用如下：

```java
public static final int COLOR_RED = 1;
public static final int COLOR_BLUE = 2;
public static final int COLOR_GREEN = 3;
```

使用`int`类型会存在以下隐患：

- 不具备安全性，声明时如果没有使用`final`就会造成值被篡改的风险；
- 语义不够明确，打印`int`型数字并不知道其具体含义

于是乎我们就想到用常量字符来表示，代码就变成了这样：

```java
public static final String COLOR_RED = "RED";
public static final String COLOR_BLUE = "BLUE";
public static final String COLOR_GREEN = "GREEN";
```

这样也同样存在问题，因为我们使用的常量字符，那么有些程序猿不按套路出牌就可以使用字符串的值进行比较，这样的代码会被不断模仿变得越来越多的，然后屎山就出现了；

所以我们迫切需要枚举类型的出现来起到约束的作用，假设使用一个枚举类型做入参，枚举类型就可以限定沙雕用户不按套路传参，这样就可以怼他了，哈哈～；

使用枚举的代码就可以变成这样，传了枚举之外的类型都不可以了；

```java
public class EnumClass {
    public static void main(String [] args){
        Color color = Color.RED;
        convert(color);
        System.out.println(color.name());

    }

    public static void convert(Color c){
        System.out.println(c.name());
    }

}

enum Color{
    RED,BLUE,GREEN;
}
```

Go语言就没有枚举类型，我们该使用什么方法来替代呢？



## 定义新类型实现枚举

枚举通常是一组相关的常量集合，`Go`语言中有提供常量类型，所以我们可以使用常量来声明枚举，但也同样会遇到上述的问题，起不到约束的作用，所以为了起到约束我们可以使用`Go`语言另外一个知识点 -- 类型定义，`Go`语言中可以使用`type`关键字定义不同的类型，我们可以为整型、浮点型、字符型等定义新的类型，新的类型与原类型转换需要显式转换，这样在一定程度上也起到了约束的作用，我们就可以用`Go`语言实现如下枚举：

```go
type OrderStatus int

const (
	CREATE OrderStatus = iota + 1
	PAID
	DELIVERING
	COMPLETED
	CANCELLED
)

func main() {
	a := 100
	IsCreated(a)
}
```

上面的代码就会报错：

```go
./main.go:19:12: cannot use a (variable of type int) as type OrderStatus in argument to IsCreated
```

定义新的类型可以起到约束作用，比如我们要检查状态机，入参限定了必须是OrderStatus类型，如果是`int`类型就会报错。

上面我们的枚举实现方式只能获取枚举值，获取不到其映射的字面意思，所以我们可以优化一下，实现`String`方法，使用官方提供的cmd/string来快速实现，代码如下：

```go
//go:generate stringer -type=OrderStatus
type OrderStatus int

const (
	CREATE OrderStatus = iota + 1
	PAID
	DELIVERING
	COMPLETED
	CANCELLED
)
```

执行命令`go generate ./...`生成orderstatus_string.go文件：

```go
import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CREATE-1]
	_ = x[PAID-2]
	_ = x[DELIVERING-3]
	_ = x[COMPLETED-4]
	_ = x[CANCELLED-5]
}

const _OrderStatus_name = "CREATEPAIDDELIVERINGCOMPLETEDCANCELLED"

var _OrderStatus_index = [...]uint8{0, 6, 10, 20, 29, 38}

func (i OrderStatus) String() string {
	i -= 1
	if i < 0 || i >= OrderStatus(len(_OrderStatus_index)-1) {
		return "OrderStatus(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _OrderStatus_name[_OrderStatus_index[i]:_OrderStatus_index[i+1]]
}
```



## protobuf中生成的枚举代码

`Go`语言使用protobuf会生成对应的枚举代码，我们发现其中也是使用定义新的类型的方式来实现的，然后在封装一些方法，我们来赏析一下protobuf生成的枚举代码：

```go
const (
	CREATED  OrderStatus = 1
	PAID OrderStatus = 2
	CANCELED OrderStatus = 3
)

var OrderStatus_name = map[int32]string{
	1: "CREATED",
	2: "PAID",
	3: "CANCELED",
}

var OrderStatus_value = map[string]int32{
	"CREATED":  1,
	"PAID": 2,
	"CANCELED": 3,
}

func (x OrderStatus) Enum() *OrderStatus {
	p := new(OrderStatus)
	*p = x
	return p
}

func (x OrderStatus) String() string {
	return proto.EnumName(OrderStatus_name, int32(x))
}

func (x *OrderStatus) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(OrderStatus_value, data, "OrderStatus")
	if err != nil {
		return err
	}
	*x = OrderStatus(value)
	return nil
}
```



## 总结

虽然Go语言没有提供枚举类型，但是我们也可以根据`Go`语言的两个特性：常量和定义新类型来实现枚举，方法总比困难多吗，开源库是优秀的，我们往往可以从高手那里里学习很多，记住，请永远保持一个学徒之心；

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**


![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)

