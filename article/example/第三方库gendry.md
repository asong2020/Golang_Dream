## 前言

> 哈喽，我是`asong`。
>
> 今天给大家推荐一个第三方库`gendry`，这个库是用于辅助操作数据库的`Go`包。其是基于`go-sql-driver/mysql`，它提供了一系列的方法来为你调用标准库`database/sql`中的方法准备参数。对于我这种不喜欢是使用`orm`框架的选手，真的是爱不释手，即使不使用`orm`框架，也可以写出动态`sql`。下面我就带大家看一看这个库怎么使用！
>
> github地址：https://github.com/didi/gendry



## 初始化连接

既然要使用数据库，那么第一步我们就来进行数据库连接，我们先来看一下直接使用标准库进行连接库是怎样写的：

```go
func NewMysqlClient(conf *config.Server) *sql.DB {
	connInfo := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.Mysql.Username, conf.Mysql.Password, conf.Mysql.Host, conf.Mysql.Db)
	var err error
	db, err := sql.Open("mysql", connInfo)
	if err != nil {
		fmt.Printf("init mysql err %v\n", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("ping mysql err: %v", err)
	}
	db.SetMaxIdleConns(conf.Mysql.Conn.MaxIdle)
	db.SetMaxOpenConns(conf.Mysql.Conn.Maxopen)
	db.SetConnMaxLifetime(5 * time.Minute)
	fmt.Println("init mysql successc")
	return db
}
```

从上面的代码可以看出，我们需要自己拼接连接参数，这就需要我们时刻记住连接参数（对于我这种记忆白痴，每回都要去度娘一下，很难受）。`Gendry`为我们提供了一个`manager`库，主要用来初始化连接池，设置其各种参数，你可以设置任何`go-sql-driver/mysql`驱动支持的参数，所以我们的初始化代码可以这样写：

```go
func MysqlClient(conf *config.Mysql) *sql.DB {

	db, err := manager.
		New(conf.Db,conf.Username,conf.Password,conf.Host).Set(
		manager.SetCharset("utf8"),
		manager.SetAllowCleartextPasswords(true),
		manager.SetInterpolateParams(true),
		manager.SetTimeout(1 * time.Second),
		manager.SetReadTimeout(1 * time.Second),
			).Port(conf.Port).Open(true)

	if err != nil {
		fmt.Printf("init mysql err %v\n", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("ping mysql err: %v", err)
	}
	db.SetMaxIdleConns(conf.Conn.MaxIdle)
	db.SetMaxOpenConns(conf.Conn.Maxopen)
	db.SetConnMaxLifetime(5 * time.Minute)
	//scanner.SetTagName("json")  // 全局设置，只允许设置一次
	fmt.Println("init mysql successc")
	return db
}
```

`manager`做的事情就是帮我们生成`datasourceName`，并且它支持了几乎所有该驱动支持的参数设置，我们完全不需要管`datasourceName`的格式是怎样的，只管配置参数就可以了。



## 如何使用？

下面我就带着大家一起来几个`demo`学习，更多使用方法可以看源代码解锁（之所以没说看官方文档解决的原因：文档不是很详细，还不过看源码来的实在）。



### 数据库准备

既然是写示例代码，那么一定要先有一个数据表来提供测试呀，测试数据表如下：

```sql
create table users
(
    id       bigint unsigned auto_increment
        primary key,
    username varchar(64)  default '' not null,
    nickname varchar(255) default '' null,
    password varchar(256) default '' not null,
    salt     varchar(48)  default '' not null,
    avatar   varchar(128)            null,
    uptime   bigint       default 0  not null,
    constraint username
        unique (username)
)
    charset = utf8mb4;
```

好了数据表也有了，下面就开始展示吧，以下按照增删改查的顺序依次展示～。

### 插入数据

`gendry`提供了三种方法帮助你构造插入sql，分别是：

```sql
// BuildInsert work as its name says
func BuildInsert(table string, data []map[string]interface{}) (string, []interface{}, error) {
	return buildInsert(table, data, commonInsert)
}

// BuildInsertIgnore work as its name says
func BuildInsertIgnore(table string, data []map[string]interface{}) (string, []interface{}, error) {
	return buildInsert(table, data, ignoreInsert)
}

// BuildReplaceInsert work as its name says
func BuildReplaceInsert(table string, data []map[string]interface{}) (string, []interface{}, error) {
	return buildInsert(table, data, replaceInsert)
}

// BuildInsertOnDuplicateKey builds an INSERT ... ON DUPLICATE KEY UPDATE clause.
func BuildInsertOnDuplicate(table string, data []map[string]interface{}, update map[string]interface{}) (string, []interface{}, error) {
	return buildInsertOnDuplicate(table, data, update)
}
```

看命名想必大家就已经知道他们代表的是什么意思了吧，这里就不一一解释了，这里我们以`buildInsert`为示例，写一个小demo：

```go
func (db *UserDB) Add(ctx context.Context,cond map[string]interface{}) (int64,error) {
	sqlStr,values,err := builder.BuildInsert(tplTable,[]map[string]interface{}{cond})
	if err != nil{
		return 0,err
	}
	// TODO:DEBUG
	fmt.Println(sqlStr,values)
	res,err := db.cli.ExecContext(ctx,sqlStr,values...)
	if err != nil{
		return 0,err
	}
	return res.LastInsertId()
}
// 单元测试如下：
func (u *UserDBTest) Test_Add()  {
	cond := map[string]interface{}{
		"username": "test_add",
		"nickname": "asong",
		"password": "123456",
		"salt": "oooo",
		"avatar": "http://www.baidu.com",
		"uptime": 123,
	}
	s,err := u.db.Add(context.Background(),cond)
	u.Nil(err)
	u.T().Log(s)
}
```

我们把要插入的数据放到`map`结构中，`key`就是要字段，`value`就是我们要插入的值，其他都交给` builder.BuildInsert`就好了，我们的代码大大减少。大家肯定很好奇这个方法是怎样实现的呢？别着急，后面我们一起解密。



### 删除数据

我最喜欢删数据了，不知道为什么，删完数据总有一种快感。。。。

删除数据可以直接调用` builder.BuildDelete`方法，比如我们现在我们要删除刚才插入的那条数据：

```go
func (db *UserDB)Delete(ctx context.Context,where map[string]interface{}) error {
	sqlStr,values,err := builder.BuildDelete(tplTable,where)
	if err != nil{
		return err
	}
	// TODO:DEBUG
	fmt.Println(sqlStr,values)
	res,err := db.cli.ExecContext(ctx,sqlStr,values...)
	if err != nil{
		return err
	}
	affectedRows,err := res.RowsAffected()
	if err != nil{
		return err
	}
	if affectedRows == 0{
		return errors.New("no record delete")
	}
	return nil
}

// 单测如下：
func (u *UserDBTest)Test_Delete()  {
	where := map[string]interface{}{
		"username in": []string{"test_add"},
	}
	err := u.db.Delete(context.Background(),where)
	u.Nil(err)
}
```

这里在传入`where`条件时，`key`使用的`username in`，这里使用空格加了一个操作符`in`，这是`gendry`库所支持的写法，当我们的`SQL`存在一些操作符时，就可以通过这样方法进行书写，形式如下：

```go
where := map[string]interface{}{
    "field 操作符": "value",
}
```

官文文档给出的支持操作如下：

```go
=
>
<
=
<=
>=
!=
<>
in
not in
like
not like
between
not between
```

既然说到了这里，顺便把`gendry`支持的关键字也说一下吧，官方文档给出的支持如下：

```go
_or
_orderby
_groupby
_having
_limit
_lockMode
```

参考示例：

```go
where := map[string]interface{}{
    "age >": 100,
    "_or": []map[string]interface{}{
        {
            "x1":    11,
            "x2 >=": 45,
        },
        {
            "x3":    "234",
            "x4 <>": "tx2",
        },
    },
    "_orderby": "fieldName asc",
    "_groupby": "fieldName",
    "_having": map[string]interface{}{"foo":"bar",},
    "_limit": []uint{offset, row_count},
    "_lockMode": "share",
}
```

这里有几个需要注意的问题：

- 如果`_groupby`没有被设置将忽略`_having`
- `_limit`可以这样写：
  - `"_limit": []uint{a,b}` => `LIMIT a,b`
  - `"_limit": []uint{a}` => `LIMIT 0,a`
- `_lockMode`暂时只支持`share`和`exclusive`
  - `share`代表的是`SELECT ... LOCK IN SHARE MODE`.不幸的是，当前版本不支持`SELECT ... FOR SHARE`.
  - `exclusive`代表的是`SELECT ... FOR UPDATE`.



## 更新数据

更新数据可以使用`builder.BuildUpdate`方法进行构建`sql`语句，不过要注意的是，他不支持`_orderby`、`_groupby`、`_having`.只有这个是我们所需要注意的，其他的正常使用就可以了。

```go

func (db *UserDB) Update(ctx context.Context,where map[string]interface{},data map[string]interface{}) error {
	sqlStr,values,err := builder.BuildUpdate(tplTable,where,data)
	if err != nil{
		return err
	}
	// TODO:DEBUG
	fmt.Println(sqlStr,values)
	res,err := db.cli.ExecContext(ctx,sqlStr,values...)
	if err != nil{
		return err
	}
	affectedRows,err := res.RowsAffected()
	if err != nil{
		return err
	}
	if affectedRows == 0{
		return errors.New("no record update")
	}
	return nil
}
// 单元测试如下：
func (u *UserDBTest) Test_Update()  {
	where := map[string]interface{}{
		"username": "asong",
	}
	data := map[string]interface{}{
		"nickname": "shuai",
	}
	err := u.db.Update(context.Background(),where,data)
	u.Nil(err)
}
```

这里入参变成了两个，一个是用来指定`where`条件的，另一个就是来放我们要更新的数据的。



### 查询数据

查询使用的是`builder.BuildSelect`方法来构建`sql`语句，先来一个示例，看看怎么用？

```go

func (db *UserDB) Query(ctx context.Context,cond map[string]interface{}) ([]*model.User,error) {
	sqlStr,values,err := builder.BuildSelect(tplTable,cond,db.getFiledList())
	if err != nil{
		return nil, err
	}
	rows,err := db.cli.QueryContext(ctx,sqlStr,values...)
	defer func() {
		if rows != nil{
			_ = rows.Close()
		}
	}()
	if err != nil{
		if err == sql.ErrNoRows{
			return nil,errors.New("not found")
		}
		return nil,err
	}
	user := make([]*model.User,0)
	err = scanner.Scan(rows,&user)
	if err != nil{
		return nil,err
	}
	return user,nil
}
// 单元测试
func (u *UserDBTest) Test_Query()  {
	cond := map[string]interface{}{
		"id in": []int{1,2},
	}
	s,err := u.db.Query(context.Background(),cond)
	u.Nil(err)
	for k,v := range s{
		u.T().Log(k,v)
	}
}
```

`BuildSelect(table string, where map[string]interface{}, selectField []string)`总共有三个入参，`table`就是数据表名，`where`里面就是我们的条件参数，`selectFiled`就是我们要查询的字段，如果传`nil`，对应的`sql`语句就是`select * ...`。看完上面的代码，系统的朋友应该会对`scanner.Scan`，这个就是`gendry`提供一个映射结果集的方法，下面我们来看一看这个库怎么用。



### scanner

执行了数据库操作之后，要把返回的结果集和自定义的struct进行映射。Scanner提供一个简单的接口通过反射来进行结果集和自定义类型的绑定，上面的`scanner.Scan`方法就是来做这个，scanner进行反射时会使用结构体的tag。默认使用的tagName是`ddb:"xxx"`，你也可以自定义。使用`scanner.SetTagName("json")`进行设置，**scaner.SetTagName是全局设置，为了避免歧义，只允许设置一次，一般在初始化DB阶段进行此项设置**.

有时候我们可能不太想定义一个结构体去存中间结果，那么`gendry`还提供了`scanMap`可以使用：

```go
rows,_ := db.Query("select name,m_age from person")
result,err := scanner.ScanMap(rows)
for _,record := range result {
	fmt.Println(record["name"], record["m_age"])
}
```

在使用`scanner`是有以下几点需要注意：

- 如果是使用Scan或者ScanMap的话，你必须在之后手动close rows
- 传给Scan的必须是引用
- ScanClose和ScanMapClose不需要手动close rows

### 

### 手写`SQL`

对于一些比较复杂的查询，`gendry`方法就不能满足我们的需求了，这就可能需要我们自定义`sql`了，`gendry`提供了`NamedQuery`就是这么使用的，具体使用如下：

```go

func (db *UserDB) CustomizeGet(ctx context.Context,sql string,data map[string]interface{}) (*model.User,error) {
	sqlStr,values,err := builder.NamedQuery(sql,data)
	if err != nil{
		return nil, err
	}
	// TODO:DEBUG
	fmt.Println(sql,values)
	rows,err := db.cli.QueryContext(ctx,sqlStr,values...)
	if err != nil{
		return nil,err
	}
	defer func() {
		if rows != nil{
			_ = rows.Close()
		}
	}()
	user := model.NewEmptyUser()
	err = scanner.Scan(rows,&user)
	if err != nil{
		return nil,err
	}
	return user,nil
}
// 单元测试
func (u *UserDBTest) Test_CustomizeGet()  {
	sql := "SELECT * FROM users WHERE username={{username}}"
	data := map[string]interface{}{
		"username": "test_add",
	}
	user,err := u.db.CustomizeGet(context.Background(),sql,data)
	u.Nil(err)
	u.T().Log(user)
}
```

这种就是纯手写`sql`了，一些复杂的地方可以这么使用。



### 聚合查询

`gendry`还为我们提供了聚合查询，例如：count,sum,max,min,avg。这里就拿`count`来举例吧，假设我们现在要统计密码相同的用户有多少，就可以这么写：

```go

func (db *UserDB) AggregateCount(ctx context.Context,where map[string]interface{},filed string) (int64,error) {
	res,err := builder.AggregateQuery(ctx,db.cli,tplTable,where,builder.AggregateCount(filed))
	if err != nil{
		return 0, err
	}
	numberOfRecords := res.Int64()
	return numberOfRecords,nil
}
// 单元测试
func (u *UserDBTest) Test_AggregateCount()  {
	where := map[string]interface{}{
		"password": "123456",
	}
	count,err := u.db.AggregateCount(context.Background(),where,"*")
	u.Nil(err)
	u.T().Log(count)
}
```

到这里，所有的基本用法基本演示了一遍，更多的使用方法可以自行解锁。

## cli工具

除了上面这些`API`以外，`Gendry`还提供了一个命令行来进行代码生成，可以显著减少你的开发量，`gforge`是基于[gendry](https://github.com/caibirdme/gforge/blob/master/github.com/didi/gendry)的cli工具，它根据表名生成golang结构，这可以减轻您的负担。甚至gforge都可以为您生成完整的DAO层。

### 安装

```go
go get -u github.com/caibirdme/gforge
```

使用`gforge -h`来验证是否安装成功，同时会给出使用提示。

### 生成表结构

使用`gforge`生成的表结构是可以通过`golint`和` govet`的。生成指令如下：

```go
gforge table -uroot -proot1997 -h127.0.0.1 -dasong -tusers

// Users is a mapping object for users table in mysql
type Users struct {
	ID uint64 `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Salt string `json:"salt"`
	Avatar string `json:"avatar"`
	Uptime int64 `json:"uptime"`
}
```

这样就省去了我们自定义表结构的时间，或者更方便的是直接把`dao`层生成出来。

### 生成`dao`文件

运行指令如下：

```go
gforge dao -uroot -proot1997 -h127.0.0.1 -dasong -tusers | gofmt > dao.go
```

这里我把生成的`dao`层直接丢到了文件里了，这里就不贴具体代码了，没有意义，知道怎么使用就好了。



## 解密

想必大家一定都跟我一样特别好奇`gendry`是怎么实现的呢？下面就以`builder.buildSelect`为例子，我们来看一看他是怎么实现的。其他原理相似，有兴趣的童鞋可以看源码学习。我们先来看一下`buildSelect`这个方法的源码：

```go
func BuildSelect(table string, where map[string]interface{}, selectField []string) (cond string, vals []interface{}, err error) {
	var orderBy string
	var limit *eleLimit
	var groupBy string
	var having map[string]interface{}
	var lockMode string
	if val, ok := where["_orderby"]; ok {
		s, ok := val.(string)
		if !ok {
			err = errOrderByValueType
			return
		}
		orderBy = strings.TrimSpace(s)
	}
	if val, ok := where["_groupby"]; ok {
		s, ok := val.(string)
		if !ok {
			err = errGroupByValueType
			return
		}
		groupBy = strings.TrimSpace(s)
		if "" != groupBy {
			if h, ok := where["_having"]; ok {
				having, err = resolveHaving(h)
				if nil != err {
					return
				}
			}
		}
	}
	if val, ok := where["_limit"]; ok {
		arr, ok := val.([]uint)
		if !ok {
			err = errLimitValueType
			return
		}
		if len(arr) != 2 {
			if len(arr) == 1 {
				arr = []uint{0, arr[0]}
			} else {
				err = errLimitValueLength
				return
			}
		}
		begin, step := arr[0], arr[1]
		limit = &eleLimit{
			begin: begin,
			step:  step,
		}
	}
	if val, ok := where["_lockMode"]; ok {
		s, ok := val.(string)
		if !ok {
			err = errLockModeValueType
			return
		}
		lockMode = strings.TrimSpace(s)
		if _, ok := allowedLockMode[lockMode]; !ok {
			err = errNotAllowedLockMode
			return
		}
	}
	conditions, err := getWhereConditions(where, defaultIgnoreKeys)
	if nil != err {
		return
	}
	if having != nil {
		havingCondition, err1 := getWhereConditions(having, defaultIgnoreKeys)
		if nil != err1 {
			err = err1
			return
		}
		conditions = append(conditions, nilComparable(0))
		conditions = append(conditions, havingCondition...)
	}
	return buildSelect(table, selectField, groupBy, orderBy, lockMode, limit, conditions...)
}
```

- 首先会对几个关键字进行处理。
- 然后会调用`getWhereConditions`这个方法去构造`sql`，看一下内部实现(摘取部分)：

```go
for key, val := range where {
		if _, ok := ignoreKeys[key]; ok {
			continue
		}
		if key == "_or" {
			var (
				orWheres          []map[string]interface{}
				orWhereComparable []Comparable
				ok                bool
			)
			if orWheres, ok = val.([]map[string]interface{}); !ok {
				return nil, errOrValueType
			}
			for _, orWhere := range orWheres {
				if orWhere == nil {
					continue
				}
				orNestWhere, err := getWhereConditions(orWhere, ignoreKeys)
				if nil != err {
					return nil, err
				}
				orWhereComparable = append(orWhereComparable, NestWhere(orNestWhere))
			}
			comparables = append(comparables, OrWhere(orWhereComparable))
			continue
		}
		field, operator, err = splitKey(key)
		if nil != err {
			return nil, err
		}
		operator = strings.ToLower(operator)
		if !isStringInSlice(operator, opOrder) {
			return nil, ErrUnsupportedOperator
		}
		if _, ok := val.(NullType); ok {
			operator = opNull
		}
		wms.add(operator, field, val)
	}
```

这一段就是遍历`slice`，之前处理过的关键字部分会被忽略，`_or`关键字会递归处理得到所有条件数据。之后就没有特别要说明的地方了。我自己返回到`buildSelect`方法中，在处理了`where`条件之后，如果有`having`条件还会在进行一次过滤，最后所有的数据构建好了后，会调用`buildSelect`方法来构造最后的`sql`语句。



## 总结

看过源码以后，只想说：大佬就是大佬。源码其实很容易看懂，这就没有做详细的解析，主要是这样思想值得大家学习，建议大家都可以看一遍`gendry`的源码，涨知识～～。

**好啦，这篇文章就到这里啦，素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！**

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

