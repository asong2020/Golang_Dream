## 前言

> 哈喽，大家好，我是`asong`。在上一篇文章：[小白也能看懂的context包详解：从入门到精通](https://mp.weixin.qq.com/s/_5gBIwvtXKJME7AV2W2bqQ) 分析`context`的源码时，我们看到了一种编程方法，在结构体里内嵌匿名接口，这种写法对于大多数初学`Go`语言的朋友看起来是懵逼的，其实在结构体里内嵌匿名接口、匿名结构体都是在面向对象编程中继承和重写的一种实现方式，之前写过`java`、`python`对面向对象编程中的继承和重写应该很熟悉，但是转`Go`语言后写出的代码都是面向过程式的代码，所以本文就一起来分析一下如何在`Go`语言中写出面向对象的代码。



面向对象程序设计是一种计算机编程架构，英文全称：Object Oriented Programming，简称OOP。OOP的一条基本原则是计算机程序由单个能够起到子程序作用的单元或对象组合而成，OOP达到了软件工程的三个主要目标：重用性、灵活性和扩展性。OOP=对象+类+继承+多态+消息，其中核心概念就是类和对象。

这一段话在网上介绍什么是面向对象编程时经常出现，大多数学习`Go`语言的朋友应该也都是从`C++`、`python`、`java`转过来的，所以对面向对象编程的理解应该很深了，所以本文就没必要介绍概念了，重点来看一下如何使用`Go`语言来实现面向对象编程。



## 类

`Go`语言本身就不是一个面向对象的编程语言，所以`Go`语言中没有类的概念，但是他是支持类型的，因此我们可以使用`struct`类型来提供类似于`java`中的类的服务，可以定义属性、方法、还能定义构造器。来看个例子：

```go
type Hero struct {
	Name string
	Age uint64
}

func NewHero() *Hero {
	return &Hero{
		Name: "盖伦",
		Age: 18,
	}
}

func (h *Hero) GetName() string {
	return h.Name
}

func (h *Hero) GetAge() uint64 {
	return h.Age
}


func main()  {
	h := NewHero()
	print(h.GetName())
	print(h.GetAge())
}
```

这就一个简单的 "类"的使用，这个类名就是`Hero`，其中`Name`、`Age`就是我们定义的属性，`GetName`、`GetAge`这两个就是我们定义的类的方法，`NewHero`就是定义的构造器。因为`Go`语言的特性问题，构造器只能够依靠我们手动来实现。

这里方法的实现是依赖于结构体的值接收者、指针接收者的特性来实现的。



## 封装

封装是把一个对象的属性私有化，同时提供一些可以被外界访问的属性和方法，如果不想被外界方法，我们大可不必提供方法给外界访问。在`Go`语言中实现封装我们可以采用两种方式：

- `Go`语言支持包级别的封装，小写字母开头的名称只能在该程序中可见，所以我们如果不想暴露一些方法，可以通过这种方式私有包中的内容，这个理解比较简单，就不举例子了。
- `Go`语言可以通过 `type` 关键字创建新的类型，所以我们为了不暴露一些属性和方法，可以采用创建一个新类型的方式，自己手写构造器的方式实现封装，举个例子：

```go
type IdCard string

func NewIdCard(card string) IdCard {
	return IdCard(card)
}

func (i IdCard) GetPlaceOfBirth() string {
	return string(i[:6])
}

func (i IdCard) GetBirthDay() string {
	return string(i[6:14])
}
```

声明一个新类型`IdCard`，本质是一个`string`类型，`NewIdCard`用来构造对象，

`GetPlaceOfBirth`、`GetBirthDay`就是封装的方法。



## 继承

`Go`并没有原生级别的继承支持，不过我们可以使用组合的方式来实现继承，通过结构体内嵌类型的方式实现继承，典型的应用是内嵌匿名结构体类型和内嵌匿名接口类型，这两种方式还有点细微差别：

- 内嵌匿名结构体类型：将父结构体嵌入到子结构体中，子结构体拥有父结构体的属性和方法，但是这种方式不能支持参数多态。
- 内嵌匿名接口类型：将接口类型嵌入到结构体中，该结构体默认实现了该接口的所有方法，该结构体也可以对这些方法进行重写，这种方式可以支持参数多态。

### 内嵌匿名结构体类型实现继承的一个例子

```go
type Base struct {
	Value string
}

func (b *Base) GetMsg() string {
	return b.Value
}


type Person struct {
	Base
	Name string
	Age uint64
}

func (p *Person) GetName() string {
	return p.Name
}

func (p *Person) GetAge() uint64 {
	return p.Age
}

func check(b *Base)  {
	b.GetMsg()
}

func main()  {
	m := Base{Value: "I Love You"}
	p := &Person{
		Base: m,
		Name: "asong",
		Age: 18,
	}
	fmt.Print(p.GetName(), "  ", p.GetAge(), " and say ",p.GetMsg())
	//check(p)
}
```

上面注释掉的方法就证明了不能进行参数多态。

### 内嵌匿名接口类型实现继承的例子

直接拿一个业务场景举例子，假设现在我们现在要给用户发一个通知，`web`、`app`端发送的通知内容都是一样的，但是点击后的动作是不一样的，所以我们可以进行抽象一个接口`OrderChangeNotificationHandler`来声明出三个公共方法：`GenerateMessage`、`GeneratePhotos`、`generateUrl`，所有类都会实现这三个方法，因为`web`、`app`端发送的内容是一样的，所以我们可以抽相出一个父类`OrderChangeNotificationHandlerImpl`来实现一个默认的方法，然后在写两个子类`WebOrderChangeNotificationHandler`、`AppOrderChangeNotificationHandler`去继承父类重写`generateUrl`方法即可，后面如果不同端的内容有做修改，直接重写父类方法就可以了，来看例子：

```go
type Photos struct {
	width uint64
	height uint64
	value string
}

type OrderChangeNotificationHandler interface {
	GenerateMessage() string
	GeneratePhotos() Photos
	generateUrl() string
}


type OrderChangeNotificationHandlerImpl struct {
	url string
}

func NewOrderChangeNotificationHandlerImpl() OrderChangeNotificationHandler {
	return OrderChangeNotificationHandlerImpl{
		url: "https://base.test.com",
	}
}

func (o OrderChangeNotificationHandlerImpl) GenerateMessage() string {
	return "OrderChangeNotificationHandlerImpl GenerateMessage"
}

func (o OrderChangeNotificationHandlerImpl) GeneratePhotos() Photos {
	return Photos{
		width: 1,
		height: 1,
		value: "https://www.baidu.com",
	}
}

func (w OrderChangeNotificationHandlerImpl) generateUrl() string {
	return w.url
}

type WebOrderChangeNotificationHandler struct {
	OrderChangeNotificationHandler
	url string
}

func (w WebOrderChangeNotificationHandler) generateUrl() string {
	return w.url
}

type AppOrderChangeNotificationHandler struct {
	OrderChangeNotificationHandler
	url string
}

func (a AppOrderChangeNotificationHandler) generateUrl() string {
	return a.url
}

func check(handler OrderChangeNotificationHandler)  {
	fmt.Println(handler.GenerateMessage())
}

func main()  {
	base := NewOrderChangeNotificationHandlerImpl()
	web := WebOrderChangeNotificationHandler{
		OrderChangeNotificationHandler: base,
		url: "http://web.test.com",
	}
	fmt.Println(web.GenerateMessage())
	fmt.Println(web.generateUrl())

	check(web)
}
```

因为所有组合都实现了`OrderChangeNotificationHandler`类型，所以可以处理任何特定类型以及是该特定类型的派生类的通配符。



## 多态

多态是面向对象编程的本质，多态是支代码可以根据类型的具体实现采取不同行为的能力，在`Go`语言中任何用户定义的类型都可以实现任何接口，所以通过不同实体类型对接口值方法的调用就是多态，举个例子：

```go
type SendEmail interface {
	send()
}

func Send(s SendEmail)  {
	s.send()
}

type user struct {
	name string
	email string
}

func (u *user) send()  {
	fmt.Println(u.name + " email is " + u.email + "already send")
}

type admin struct {
	name string
	email string
}

func (a *admin) send()  {
	fmt.Println(a.name + " email is " + a.email + "already send")
}

func main()  {
	u := &user{
		name: "asong",
		email: "你猜",
	}
	a := &admin{
		name: "asong1",
		email: "就不告诉你",
	}
	Send(u)
	Send(a)
}
```



## 总结

归根结底面向对象编程就是一种编程思想，只不过有些语言在语法特性方面更好的为这种思想提供了支持，对于写出面向对象的代码更容易，但是写代码的还是我们自己的，并不是我们用了`java`就一定会写出更抽象的代码，在工作中我看到用`java`写出面向过程式的代码不胜其数，所以这就是一种编程思想，无论用什么语言，我们都应该思考如何写好一份代码，大量的抽象接口帮助我们精简代码，代码是优雅了，但也会面临着可读性的问题，什么事都是有两面性的，写出好代码的路还很长，还需要不断探索............。

文中示例代码已经上传`github`：

**好啦，本文到这里就结束了，我是`asong`，我们下期见。**

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)