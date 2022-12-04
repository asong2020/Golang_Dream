## 前言

> 今天逛`github`时，发现了这款对 SQL 进行优化和改写的自动化工具`sora`。感觉挺不错的，就下载学习了一下。这个工具支持的功能比较多，可以作为我们日常开发中的一款辅助工具，现在我就把它推荐给你们～～～
>
> github传送门：https://github.com/XiaoMi/soar



## 背景

在我们日常开发中，优化SQL总是我们日常开发任务之一。例行 SQL 优化，不仅可以提升程序性能，还能够降低线上故障的概率。

目前常用的 SQL 优化方式包括但不限于：业务层优化、SQL逻辑优化、索引优化等。其中索引优化通常通过调整索引或新增索引从而达到 SQL 优化的目的。索引优化往往可以在短时间内产生非常巨大的效果。如果能够将索引优化转化成工具化、标准化的流程，减少人工介入的工作量，无疑会大大提高我们的工作效率。

SOAR(SQL Optimizer And Rewriter) 是一个对 SQL 进行优化和改写的自动化工具。 由小米人工智能与云平台的数据库团队开发与维护。

与业内其他优秀产品对比如下：

| SOAR         | sqlcheck | pt-query-advisor | SQL Advisor | Inception | sqlautoreview |      |
| ------------ | -------- | ---------------- | ----------- | --------- | ------------- | ---- |
| 启发式建议   | ✔️        | ✔️                | ✔️           | ❌         | ✔️             | ✔️    |
| 索引建议     | ✔️        | ❌                | ❌           | ✔️         | ❌             | ✔️    |
| 查询重写     | ✔️        | ❌                | ❌           | ❌         | ❌             | ❌    |
| 执行计划展示 | ✔️        | ❌                | ❌           | ❌         | ❌             | ❌    |
| Profiling    | ✔️        | ❌                | ❌           | ❌         | ❌             | ❌    |
| Trace        | ✔️        | ❌                | ❌           | ❌         | ❌             | ❌    |
| SQL在线执行  | ❌        | ❌                | ❌           | ❌         | ✔️             | ❌    |
| 数据备份     | ❌        | ❌                | ❌           | ❌         | ✔️             | ❌    |

从上图可以看出，支持的功能丰富，其功能特点如下：

- 跨平台支持（支持 Linux, Mac 环境，Windows 环境理论上也支持，不过未全面测试）
- 目前只支持 MySQL 语法族协议的 SQL 优化
- 支持基于启发式算法的语句优化
- 支持复杂查询的多列索引优化（UPDATE, INSERT, DELETE, SELECT）
- 支持 EXPLAIN 信息丰富解读
- 支持 SQL 指纹、压缩和美化
- 支持同一张表多条 ALTER 请求合并
- 支持自定义规则的 SQL 改写

就介绍这么多吧，既然是SQL优化工具，光说是没有用的，我们还是先用起来看看效果吧。



## 安装

这里有两种安装方式，如下：

- 1. 下载二进制安装包

```shell
$ wget https://github.com/XiaoMi/soar/releases/download/0.11.0/soar.linux-amd64 -O soar
chmod a+x soar
```

这里建议直接下载最新版，要不会有`bug`。

下载好的二进制文件添加到环境变量中即可(不会的谷歌一下吧，这里就不讲了)。

测试一下：

```shell
$ echo 'select * from user' | soar.darwin-amd64(根据你自己的二进制文件名来输入)
# Query: AC4262B5AF150CB5

★ ★ ★ ☆ ☆ 75分

​```sql

SELECT
  *
FROM
  USER
​```

## 最外层 SELECT 未指定 WHERE 条件

* **Item:**  CLA.001

* **Severity:**  L4

* **Content:**  SELECT 语句没有 WHERE 子句，可能检查比预期更多的行(全表扫描)。对于 SELECT COUNT(\*) 类型的请求如果不要求精度，建议使用 SHOW TABLE STATUS 或 EXPLAIN 替代。

## 不建议使用 SELECT * 类型查询

* **Item:**  COL.001

* **Severity:**  L1

* **Content:**  当表结构变更时，使用 \* 通配符选择所有列将导致查询的含义和行为会发生更改，可能导致查询返回更多的数据。
```



- 2. 源码安装

依赖环境：

```go
1. Go 1.10+
2. git
```

高级依赖（仅面向开发人员）

- [mysql](https://dev.mysql.com/doc/refman/8.0/en/mysql.html) 客户端版本需要与容器中MySQL版本相同，避免出现由于认证原因导致无法连接问题
- [docker](https://docs.docker.com/engine/reference/commandline/cli/) MySQL Server测试容器管理
- [govendor](https://github.com/kardianos/govendor) Go包管理
- [retool](https://github.com/twitchtv/retool) 依赖外部代码质量静态检查工具二进制文件管理

生成二进制文件：

```go
go get -d github.com/XiaoMi/soar
cd ${GOPATH}/src/github.com/XiaoMi/soar && make
```

生成的二进制文件与上面一样，直接放入环境变量即可，这里我没有尝试，靠你们自己踩坑了呦～～～



## 简单使用

### 0. 前置准备

准备一个`table`，如下：

```sql
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL DEFAULT '',
  `nickname` varchar(255) DEFAULT '',
  `password` varchar(256) NOT NULL DEFAULT '',
  `salt` varchar(48) NOT NULL DEFAULT '',
  `avatar` varchar(128) DEFAULT NULL,
  `uptime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4
```

### 1. 直接输入sql语句(不运行)

```shell
$ echo "select * from users" | soar.darwin-amd64
$ # Query: 30AFCB1E1344BEBD

★ ★ ★ ☆ ☆ 75分

​```sql

SELECT
  *
FROM
  users
​```

## 最外层 SELECT 未指定 WHERE 条件

* **Item:**  CLA.001

* **Severity:**  L4

* **Content:**  SELECT 语句没有 WHERE 子句，可能检查比预期更多的行(全表扫描)。对于 SELECT COUNT(\*) 类型的请求如果不要求精度，建议使用 SHOW TABLE STATUS 或 EXPLAIN 替代。

## 不建议使用 SELECT * 类型查询

* **Item:**  COL.001

* **Severity:**  L1

* **Content:**  当表结构变更时，使用 \* 通配符选择所有列将导致查询的含义和行为会发生更改，可能导致查询返回更多的数据。
```

现在是完全根据SQL语句进行分析的，因为没有连接到`mysql`。可以看到，给出的报告也很详细，但是只是空壳子，仅凭`SQL`语句给出的分析并不是准确的，所以我们开始接下来的应用。



### 2. 连接`mysql`生成`EXPLAIN`分析报告

我们可以在配置文件中配置好`mysql`相关的配置，操作如下：

```go
vi soar.yaml
# yaml format config file
online-dsn:
    addr:     127.0.0.1:3306
    schema:   asong
    user:     root
    password: root1997
    disable:  false

test-dsn:
    addr:     127.0.0.1:3306
    schema:   asong
    user:     root
    password: root1997
    disable:  false
```

配置好了，我们来实践一下子吧：

```shell
$ echo "SELECT id,username,nickname,password,salt,avatar,uptime FROM users WHERE username = 'asong1111'" | soar.darwin-amd64 -test-dsn="root:root1997@127.0.0.1:3306/asong" -allow-online-as-test -log-output=soar.log
$ # Query: D12A420193AD1674

★ ★ ★ ★ ★ 100分

​```sql

SELECT
  id, username, nickname, PASSWORD, salt, avatar, uptime
FROM
  users
WHERE
  username  = 'asong1111'
​```

##  Explain信息

| id | select\_type | table | partitions | type | possible_keys | key | key\_len | ref | rows | filtered | scalability | Extra |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| 1  | SIMPLE | *users* | NULL | const | username | username | 258 | const | 1 | ☠️ **100.00%** | ☠️ **O(n)** | NULL |



### Explain信息解读

#### SelectType信息解读

* **SIMPLE**: 简单SELECT(不使用UNION或子查询等).

#### Type信息解读

* **const**: const用于使用常数值比较PRIMARY KEY时, 当查询的表仅有一行时, 使用system. 例:SELECT * FROM tbl WHERE col = 1.
```

这回结果中多了EXPLAIN信息分析报告。这对于刚开始入门的小伙伴们是友好的，因为我们对`Explain`解析的字段并不熟悉，有了它我们可以完美的分析`SQL`中的问题，是不是很棒。



### 3. 语法检查

`soar`工具不仅仅可以进行`sql`语句分析，还可以进行对`sql`语法进行检查，找出其中的问题，来看个例子：

```shell
$ echo "selec * from users" | soar.darwin-amd64 -only-syntax-check
At SQL 1 : line 1 column 5 near "selec * from users" (total length 18)
```

这里`select`关键字少了一个`t`，运行该指令帮助我们一下就定位了问题，当我们的`sql`语句很长时，就可以使用该指令来辅助我们检查`SQL`语句是否正确。



### 4. SQL美化

我们日常开发时，经常会看其他人写的代码，因为水平不一样，所以有些`SQL`语句会写的很乱，所以这个工具就派上用场了，我们可以把我们的`SQL`语句变得漂亮一些，更容易我们理解哦。

```shell
$ echo "SELECT id,username,nickname,password,salt,avatar,uptime FROM users WHERE username = 'asong1111'" | soar.darwin-amd64 -report-type=pretty

SELECT
  id, username, nickname, PASSWORD, salt, avatar, uptime
FROM
  users
WHERE
  username  = 'asong1111';
```

这样看起来是不是更直观了呢～～。



## 结尾

因为我也才是刚使用这个工具，更多的玩法我还没有发现，以后补充。更多玩法可以自己研究一下，github传送门：https://github.com/XiaoMi/soar。官方文档其实很粗糙，更多方法解锁还要靠自己研究，毕竟源码已经给我们了，对于学习`go`也有一定帮助，当作一个小项目慢慢优化岂不是更好呢～～。

**好啦，这一篇文章到这就结束了，我们下期见～～。希望对你们有用，又不对的地方欢迎指出，可添加我的golang交流群，我们一起学习交流。**

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**翻译了一份Machinery中文文档，会定期进行维护，有需要的小伙伴们后台回复[machinery]即可获取。**

**我是asong，一名普普通通的程序猿，让gi我一起慢慢变强吧。我自己建了一个`golang`交流群，有需要的小伙伴加我`vx`,我拉你入群。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [machinery-go异步任务队列](https://mp.weixin.qq.com/s/4QG69Qh1q7_i0lJdxKXWyg)
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
- [掌握这些Go语言特性，你的水平将提高N个档次(二)](https://mp.weixin.qq.com/s/7yyo83SzgQbEB7QWGY7k-w)





