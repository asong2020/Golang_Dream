## 前言

哈喽，大家好，我是asong。这是我的第十篇原创文章。这周在公司做项目，在做API部分开发时，需要对请求参数的校验，防止用户的恶意请求。例如日期格式，用户年龄，性别等必须是正常的值，不能随意设置。最开始在做这一部分的时候，我采用老方法，自己编写参数检验方法，统一进行参数验证。后来在同事CR的时候，说GIN有更好的参数检验方法，gin框架使用[github.com/go-playground/validator](https://github.com/go-playground/validator)进行参数校验，我们只需要在定义结构体时使用`binding`或`validate`tag标识相关校验规则，就可以进行参数校验了，很方便。相信也有很多小伙伴不知道这个功能，今天就来介绍一下这部分。

`自己翻译了一份gin官方中文文档。关注公众号[Golang梦工厂]（扫描下方二维码），后台回复：gin，即可获取。`

## 快速安装

使用之前，我们先要获取`validator`这个库。

```shell
# 第一次安装使用如下命令
$ go get github.com/go-playground/validator/v10
# 项目中引入包
import "github.com/go-playground/validator/v10"
```



## 简单示例

安装还是很简单的，下面我先来一个官方样例，看看是怎么使用的，然后展开分析。

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Age      uint8  `json:"age" binding:"gte=1,lte=120"`
}

func main() {

	router := gin.Default()

	router.POST("register", Register)

	router.Run(":9999")
}

func Register(c *gin.Context) {
	var r RegisterRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		fmt.Println("register failed")
		c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}
	//验证 存储操作省略.....
	fmt.Println("register success")
	c.JSON(http.StatusOK, "successful")
}

```



- 测试

```javascript
curl --location --request POST 'http://localhost:9999/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "asong",
    "nickname": "golang梦工厂",
    "email": "7418.com",
    "password": "123",
    "age": 140
}'
```

- 返回结果

```json
{
    "msg": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag\nKey: 'RegisterRequest.Age' Error:Field validation for 'Age' failed on the 'lte' tag"
}
```

看这个输出结果，我们可以看到`validator`的检验生效了，email字段不是一个合法邮箱，age字段超过了最大限制。我们只在结构体中添加tag就解决了这个问题，是不是很方便，下面我们就来学习一下具体使用。



## validator库

gin框架是使用validator.v10这个库来进行参数验证的，所以我们先来看看这个库的使用。

先安装这个库：

```shell
$ go get github.com/go-playground/validator/v10
```

然后先写一个简单的示例：

```go
package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Username string `validate:"min=6,max=10"`
	Age      uint8  `validate:"gte=1,lte=10"`
	Sex      string `validate:"oneof=female male"`
}

func main() {
	validate := validator.New()

	user1 := User{Username: "asong", Age: 11, Sex: "null"}
	err := validate.Struct(user1)
	if err != nil {
		fmt.Println(err)
	}

	user2 := User{Username: "asong111", Age: 8, Sex: "male"}
	err = validate.Struct(user2)
	if err != nil {
		fmt.Println(err)
	}

}
```

我们在结构体定义validator标签的tag，使用`validator.New()`创建一个验证器，这个验证器可以指定选项、添加自定义约束，然后在调用他的`Struct()`方法来验证各种结构对象的字段是否符合定义的约束。

上面的例子，我们在User结构体中，有三个字段：

- Name：通过min和max来进行约束，Name的字符串长度为[6,10]之间。
- Age：通过gte和lte对年轻的范围进行约束，age的大小大于1，小于10。
- Sex：通过oneof对值进行约束，只能是所列举的值，oneof列举出性别为男士🚹和女士🚺(不是硬性规定奥，可能还有别的性别)。

所以`user1`会进行报错，错误信息如下：

```shell
Key: 'User.Name' Error:Field validation for 'Name' failed on the 'min' tag
Key: 'User.Age' Error:Field validation for 'Age' failed on the 'lte' tag
Key: 'User.Sex' Error:Field validation for 'Sex' failed on the 'oneof' tag
```

各个字段违反了什么约束，一眼我们便能从错误信息中看出来。看完了简单示例，下面我就来看一看都有哪些tag，我们都可以怎么使用。本文不介绍所有的tag，更多使用方法，请到[官方文档](https://github.com/go-playground/validator)自行学习。



#### 字符串约束

- `excludesall`：不包含参数中任意的 UNICODE 字符，例如`excludesall=ab`；

- `excludesrune`：不包含参数表示的 rune 字符，`excludesrune=asong`；

- `startswith`：以参数子串为前缀，例如`startswith=hi`；

- `endswith`：以参数子串为后缀，例如`endswith=bye`。

- `contains=`：包含参数子串，例如`contains=email`；

- `containsany`：包含参数中任意的 UNICODE 字符，例如`containsany=ab`；

- `containsrune`：包含参数表示的 rune 字符，例如`containsrune=asong；

- `excludes`：不包含参数子串，例如`excludes=email`；



#### 范围约束

范围约束的字段类型分为三种：

- 对于数值，我们则可以约束其值
- 对于切片、数组和map，我们则可以约束其长度
- 对于字符串，我们则可以约束其长度

常用tag介绍：

- `ne`：不等于参数值，例如`ne=5`；
- `gt`：大于参数值，例如`gt=5`；
- `gte`：大于等于参数值，例如`gte=50`；
- `lt`：小于参数值，例如`lt=50`；
- `lte`：小于等于参数值，例如`lte=50`；
- `oneof`：只能是列举出的值其中一个，这些值必须是数值或字符串，以空格分隔，如果字符串中有空格，将字符串用单引号包围，例如`oneof=male female`。
- `eq`：等于参数值，注意与`len`不同。对于字符串，`eq`约束字符串本身的值，而`len`约束字符串长度。例如`eq=10`；
- `len`：等于参数值，例如`len=10`；
- `max`：小于等于参数值，例如`max=10`；
- `min`：大于等于参数值，例如`min=10`



#### Fields约束

- `eqfield`：定义字段间的相等约束，用于约束同一结构体中的字段。例如：`eqfield=Password`
- `eqcsfield`：约束统一结构体中字段等于另一个字段（相对），确认密码时可以使用，例如：`eqfiel=ConfirmPassword`
- `nefield`：用来约束两个字段是否相同，确认两种颜色是否一致时可以使用，例如：`nefield=Color1`
- `necsfield`：约束两个字段是否相同（相对）



#### 常用约束

- `unique`：指定唯一性约束，不同类型处理不同：

  - 对于map，unique约束没有重复的值
  - 对于数组和切片，unique没有重复的值
  - 对于元素类型为结构体的碎片，unique约束结构体对象的某个字段不重复，使用`unique=field`指定字段名
- `email`：使用`email`来限制字段必须是邮件形式，直接写`eamil`即可，无需加任何指定。
- `omitempty`：字段未设置，则忽略
- `-`：跳过该字段，不检验；
- `|`：使用多个约束，只需要满足其中一个，例如`rgb|rgba`；
- `required`：字段必须设置，不能为默认值；



好啦，就介绍这些常用的约束，更多约束学习请到文档自行学习吧，都有example供你学习，很快的。



## gin中的参数校验

学习了validator，我们也就知道了怎么在gin中使用参数校验了。这些约束是都没有变的，在`validator`中，我们直接结构体中将约束放到`validate` tag中，同样道理，在gin中我们只需将约束放到`binding`tag中就可以了。是不是很简单。

但是有些时候，并不是所有的参数校验都能满足我们的需求，所以我们可以定义自己的约束。自定义约束支持自定义结构体校验、自定义字段校验等。这里来介绍一下自定义结构体校验。

### 自定义结构体校验

当涉及到一些复杂的校验规则，这些已有的校验规则就不能满足我们的需求了。例如现在有一个需求，存在db的用户信息中创建时间与更新时间都要大于某一时间，假设是从前端传来的（当然不可能，哈哈）。现在我们来写一个简单示例，学习一下怎么对这个参数进行校验。

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Info struct {
	CreateTime time.Time `form:"create_time" binding:"required,timing" time_format:"2006-01-02"`
	UpdateTime time.Time `form:"update_time" binding:"required,timing" time_format:"2006-01-02"`
}

// 自定义验证规则断言
func timing(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

func main() {
	route := gin.Default()
	// 注册验证
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("timing", timing)
		if err != nil {
			fmt.Println("success")
		}
	}

	route.GET("/time", getTime)
	route.Run(":8080")
}

func getTime(c *gin.Context) {
	var b Info
	// 数据模型绑定查询字符串验证
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "time are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
```

写好了，下面我就来测试验证一下：

```shell
$ curl "localhost:8080/time?create_time=2020-10-11&update_time=2020-10-11"
# 结果
{"message":"time are valid!"}%
$ curl "localhost:8080/time?create_time=1997-10-11&update_time=1997-10-11"
# 结果
{"error":"Key: 'Info.CreateTime' Error:Field validation for 'CreateTime' failed on the 'timing' tag\nKey: 'Info.UpdateTime' Error:Field validation for 'UpdateTime' failed on the 'timing' tag"}%
```

这里我们看到虽然参数验证成功了，但是这里返回的错误显示的也太全了，在项目开发中不可以给前端返回这么详细的信息的，所以我们需要改造一下：

```go
func getTime(c *gin.Context) {
	var b Info
	// 数据模型绑定查询字符串验证
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "time are valid!"})
	} else {
		_, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "param is error"})
	}
}
```

这里在出现错误时返回固定错误即可。这里你也可以使用一个方法封装一下，对错误进行处理在进行返回，更多使用方法等你发觉哟。



## 小彩蛋

我们返回错误时都是英文的，当错误很长的时候，对于我这种英语渣渣，就要借助翻译软件了。所以要是能返回的错误直接是中文的就好了。`validator`库本身是支持国际化的，借助相应的语言包可以实现校验错误提示信息的自动翻译。下面就写一个代码演示一下啦。

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

// loca 通常取决于 http 请求头的 'Accept-Language'
func transInit(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() //chinese
		enT := en.New() //english
		uni := ut.New(enT, zhT, enT)

		var o bool
		trans, o = uni.GetTranslator(local)
		if !o {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}
		//register translate
		// 注册翻译器
		switch local {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,max=16,min=6"`
}

func main() {
	if err := transInit("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}
	router := gin.Default()

	router.POST("/user/login", login)

	err := router.Run(":8888")
	if err != nil {
		log.Println("failed")
	}
}

func login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		c.JSON(http.StatusOK, gin.H{
			"msg": errs.Translate(trans),
		})
		return
	}
	//login 操作省略
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}
```

我这里请求参数中限制密码的长度，来验证一下吧。

```shell
curl --location --request POST 'http://localhost:8888/user/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "asong",
    "password": "11122222222222222222"
}'
# 返回
{
    "msg": {
        "loginRequest.Password": "Password长度不能超过16个字符"
    }
}
```



看，直接显示中文了，是不是很棒，我们可以在测试的时候使用这个，上线项目不建议使用呦！！！





## 总结

好啦，这一篇文章到这里结束啦。这一篇干货还是满满的。学会这些知识点，提高我们的开发效率，省去了一些没必要写的代码。能用的轮子我们还是不要错过滴。

**我是asong，一名普普通通的程序猿，让我一起慢慢变强吧。欢迎各位的关注，我们下期见~~~**
![公众号图片](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)
推荐往期文章：

- [听说你还不会jwt和swagger-饭我都不吃了带着实践项目我就来了](https://mp.weixin.qq.com/s/z-PGZE84STccvfkf8ehTgA)
- [掌握这些Go语言特性，你的水平将提高N个档次(二)](https://mp.weixin.qq.com/s/7yyo83SzgQbEB7QWGY7k-w)
- [go实现多人聊天室，在这里你想聊什么都可以的啦！！！](https://mp.weixin.qq.com/s/H7F85CncQNdnPsjvGiemtg)
- [grpc实践-学会grpc就是这么简单](https://mp.weixin.qq.com/s/mOkihZEO7uwEAnnRKGdkLA)
- [go标准库rpc实践](https://mp.weixin.qq.com/s/d0xKVe_Cq1WsUGZxIlU8mw)
- [2020最新Gin框架中文文档 asong又捡起来了英语，用心翻译](https://mp.weixin.qq.com/s/vx8A6EEO2mgEMteUZNzkDg)

- [基于gin的几种热加载方式](https://mp.weixin.qq.com/s/CZvjXp3dimU-2hZlvsLfsw)

