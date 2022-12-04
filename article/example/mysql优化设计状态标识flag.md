## 前言

> 哈喽，everybody，我是asong。今天asong教你们一个`mysql`优化设计状态标识。学会了他，我们的DB结构看起来更清晰，也避免了DB结构过大的问题，具体怎么设计，下面你就看我怎么操作就好了～～～



## 背景

我们在很多应用场景中，通常是需要给数据加上一些标识，已表明这条数据的某个特性。比如标识用户的支付渠道，标识商家的结算方式、商品的类型等等。对于这样的具有有限固定的几个值的标识，我们通过枚举的方式来标识就可以了，但是对于一些同时具有多个属性且变化比较大的就显然不合适了，举个很简单的例子，我们在某宝上想买一个平板，这个平板的商品类型可标识为电子商品、二手商品、、手机、数码等等，对于这种场景，一个商品对应多种类型，不确定性很大，这种就不是简单的通过几个值标识就能解决的了。本文就是针对这个问题，给出了自己的一些思考。

### 

## 问题与分析

我们就拿最近刚过去的双11举个例子，在双11要开始之前，某宝就会通过各种优惠的方式发放优惠卷、积分抵扣等等福利，这样我们在双11清空购物车时享受这些优惠。这种场景其实对我们程序员来说并不是简单的实现优惠减免这么简单，这种场景更多是标识优惠以计算用户实际所需支付金额，以及为后续业绩统计、制定促销计划、提高用户活跃度等提供数据依据。下面我们根据例子进行分析：

假设当前某宝平台可以使用的优惠方式如下：

| 序号 | 优惠内容  | 使用条件               | 是否长期有效 | 备注                                              |
| ---- | --------- | ---------------------- | ------------ | ------------------------------------------------- |
| 1    | 账户余额  | 直接抵扣现金           | 是           | 用户充值获得(平台奖励吸引的充值，如：充100送10元) |
| 2    | 平台积分  | 100积分抵扣1元         | 是           | 通过参与平台活动、购物行为积累获取                |
| 3    | 满减卷5元 | 满100减5元             | 否           | 平台活动促销发放                                  |
| 4    | 免邮费    | 订单总金额符合条件即可 | 是           | 平台单笔订单总金额满199元免邮费                   |

当用户进行下单时，只要满足各优惠的使用条件时，就可以使用各种优惠。这时我们思考一个问题，数据库是怎么存储这些优惠的呢？

根据上面的举例，用户下单时可以同时使用上面4种优惠抵扣方式，也就说用户可能出现的组合有`2^4 - 1=15`种，如果我们的表结构设计成单独用一个普通标识字段来标识存储，实现起来是比较简单，但是其需要标识的组合种类实在有点多，不太利于编码与后续扩展，想一想，优惠政策会随着平台发展不断推出的，如果新加了一种优惠类型，其需要添加多少种组合标识啊，且呈指数式爆长，这种方式显然不太合理。那么有没有什么解决方案呢？

方案一：

采用另外引入一张关联表的方式，专门用一张关联表来存储订单使用的优惠组合信息，每使用一种优惠就添加一条关联记录，相比单独使用普通字段标识，这在一定程度上减少了设置标识的繁琐性，增加了灵活性（每多使用一种优惠就添加一条关联记录），但是，同时也带来了另一些问题，其中主要问题是：新增一张关联表后，数据维护起来麻烦。在互联网场景下，数据量通常是非常大的，像订单数据一般都需要进行数据库sharding，以应对数据量暴涨后数据库的读写性能瓶颈，增加系统的水平扩展能力。因此，另外增加一张数据量是订单数据本身数据量几倍的关联表也显然不太合适。

方案二：

这就是本文的重点了，也就是我们使用“特殊标识位”的方式来实现，具体思路如下：

- 我们不再直接使用十进制数字来标识存储优惠信息，而是存储一个二进制数转化后的十进制数，这些1、2、3之类的优惠数字表示占二进制数的第几位（从右至左数）；
- 具体数据的存储、读取判断通过一个通用方法进行转换。

现在我们假设使用`int32`数据类型进行存储，共32位，除去符号位，可用于标识的位数有31位，即最多可以标识31种优惠情况。

| 优惠项    | 占第几位 | 二进制数  | 十进制数 |
| --------- | -------- | --------- | -------- |
| 账户余额  | 1        | 0000 0001 | 1        |
| 平台积分  | 2        | 0000 0010 | 2        |
| 满减卷5元 | 3        | 0000 0100 | 4        |
| 免邮费    | 4        | 0000 1000 | 8        |

说明：若用户使用了账户余额,则使用二进制数 00000001 标识，若使用了平台积分，则使用二进制数 00000010 标识，存储到DB时，转换成对应十进制数分别对应1、2；若同时使用了账户余额、平台积分，则使用二进制数 00000011 标识，最终存储到DB的对应十进制数是3。其它优惠项，所占的二进制位依次类推。



## 代码样例

先看代码

```go
package main

import (
	"fmt"
)

// golang没有enum 使用const代替
const (
	TYPE_BALANCE      = 1 // type = 1
	TYPE_INTEGRAL     = 2 // type = 2
	TYPE_COUPON       = 3 // type = 3
	TYPE_FREEPOSTAGE  = 4 // type = 4
)

// 是否使用有优惠卷
func IsUseDiscount(discountType , value uint32) bool {
	return (value & (1<< (discountType-1))) > 0
}


// 设置使用
func SetDiscountValue(discountType ,value uint32) uint32{
	return value | (1 << (discountType-1))
}

func main()  {
	// 测试1 不设置优惠类型
	var flag1 uint32 = 0
	fmt.Println(IsUseDiscount(TYPE_BALANCE,flag1))
	fmt.Println(IsUseDiscount(TYPE_INTEGRAL,flag1))
	fmt.Println(IsUseDiscount(TYPE_COUPON,flag1))
	fmt.Println(IsUseDiscount(TYPE_FREEPOSTAGE,flag1))


	// 测试2 只设置一个优惠类型
	var flag2 uint32 = 0
	flag2 = SetDiscountValue(TYPE_BALANCE,flag2)
	fmt.Println(IsUseDiscount(TYPE_BALANCE,flag2))
	fmt.Println(IsUseDiscount(TYPE_INTEGRAL,flag2))
	fmt.Println(IsUseDiscount(TYPE_COUPON,flag2))
	fmt.Println(IsUseDiscount(TYPE_FREEPOSTAGE,flag2))

	// 测试3 设置两个优惠类型
	var flag3 uint32 = 0
	flag3 = SetDiscountValue(TYPE_BALANCE,flag3)
	flag3 = SetDiscountValue(TYPE_INTEGRAL,flag3)
	fmt.Println(IsUseDiscount(TYPE_BALANCE,flag3))
	fmt.Println(IsUseDiscount(TYPE_INTEGRAL,flag3))
	fmt.Println(IsUseDiscount(TYPE_COUPON,flag3))
	fmt.Println(IsUseDiscount(TYPE_FREEPOSTAGE,flag3))

	// 测试4 设置三个优惠类型
	var flag4 uint32 = 0
	flag4 = SetDiscountValue(TYPE_BALANCE,flag4)
	flag4 = SetDiscountValue(TYPE_INTEGRAL,flag4)
	flag4 = SetDiscountValue(TYPE_COUPON,flag4)
	fmt.Println(IsUseDiscount(TYPE_BALANCE,flag4))
	fmt.Println(IsUseDiscount(TYPE_INTEGRAL,flag4))
	fmt.Println(IsUseDiscount(TYPE_COUPON,flag4))
	fmt.Println(IsUseDiscount(TYPE_FREEPOSTAGE,flag4))

	// 测试5 设置四个优惠类型
	var flag5 uint32 = 0
	flag5 = SetDiscountValue(TYPE_BALANCE,flag5)
	flag5 = SetDiscountValue(TYPE_INTEGRAL,flag5)
	flag5 = SetDiscountValue(TYPE_COUPON,flag5)
	flag5 = SetDiscountValue(TYPE_FREEPOSTAGE,flag5)
	fmt.Println(IsUseDiscount(TYPE_BALANCE,flag5))
	fmt.Println(IsUseDiscount(TYPE_INTEGRAL,flag5))
	fmt.Println(IsUseDiscount(TYPE_COUPON,flag5))
	fmt.Println(IsUseDiscount(TYPE_FREEPOSTAGE,flag5))
}
```



运行结果：

```go
false
false
false
false
true
false
false
false
true
true
false
false
true
true
true
false
true
true
true
true
```

因为`go`没有枚举，所以我们使用`const`声明常量的方式来实现，定义四个常量，代表四种优惠种类，这个并不是最最终存储到DB的值，而是表示占二进制数的第几位（从右至左数，从1开始）；当需要存储优惠种类到DB中，或者从DB中查询对应的优惠种类时，通过`SetDiscountValue`和`IsUseDiscount`这两个方法对值进行设置（项目中可以封装一个文件中作为工具类）。

`SetDiscountValue`方法的实现：通过位运算来实现，`(1 << (discountType-1))`通过位移的方法来找到其在二进制中的位置，然后通过与`value`位或的方法设定所占二进制位数，最终返回设置占位后的十进制数。

`IsUseDiscount`方法的实现：`(1<< (discountType-1))`通过位移的方法来找到其在二进制中的位置，然后通过与`value`位与的方法来判断优惠项应占位是否有占位，返回判断结果。

上面就是一个使用`特殊标识位`的一个简单代码样例，这个程序还可以进行扩展与完善，等待你们的开发呦～～～。



## 总结

在这里简单总结一下使用特殊标识位的优缺点：

- 优点
  - 方便扩展，易于维护；当业务场景迅速扩展时，这种方式可以方便的标识新增的业务场景，数据也易于维护。要知道，在互联网场景下，业务的变化是非常快的，新加字段并不是那么方便。
  - 方便标识存储，一个字段就可以标识多种业务场景。
- 缺点
  - 数据的存储、查询需要转换，不够直观；相对普通的标识方式，没接触过的人需要一点时间理解这种使用特殊标识位的方式。
  - DB数据查询时，稍显繁琐。

你们学废了嘛？反正我学废了，哈哈哈哈哈～～～～～。

**好啦，这一篇文章到这就结束了，我们下期见～～。希望对你们有用，又不对的地方欢迎指出，可添加我的golang交流群，我们一起学习交流。**

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**翻译了一份Machinery中文文档，会定期进行维护，有需要的小伙伴们后台回复[machinery]即可获取。**

**我是asong，一名普普通通的程序猿，让gi我一起慢慢变强吧。我自己建了一个`golang`交流群，有需要的小伙伴加我`vx`,我拉你入群。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [machinery-go异步任务队列](https://mp.weixin.qq.com/s/4QG69Qh1q7_i0lJdxKXWyg)
- [go参数传递类型](https://mp.weixin.qq.com/s/JHbFh2GhoKewlemq7iI59Q)
- [手把手教姐姐写消息队列](https://mp.weixin.qq.com/s/0MykGst1e2pgnXXUjojvhQ)
- [常见面试题之缓存雪崩、缓存穿透、缓存击穿](https://mp.weixin.qq.com/s?__biz=MzIzMDU0MTA3Nw==&mid=2247483988&idx=1&sn=3bd52650907867d65f1c4d5c3cff8f13&chksm=e8b0902edfc71938f7d7a29246d7278ac48e6c104ba27c684e12e840892252b0823de94b94c1&token=1558933779&lang=zh_CN#rd)
- [详解Context包，看这一篇就够了！！！](https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg)
- [go-ElasticSearch入门看这一篇就够了(一)](https://mp.weixin.qq.com/s/mV2hnfctQuRLRKpPPT9XRw)
- [面试官：go中for-range使用过吗？这几个问题你能解释一下原因吗](https://mp.weixin.qq.com/s/G7z80u83LTgLyfHgzgrd9g)

