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

现在是完全根据SQL语句进行分析的，因为没有连接到`mysql`。可以看到，给出的报告也很详细，但是只是空壳子，仅凭`SQL`语句给出的分析并不是准确的，所以我开始接下来的应用。



### 2. 连接`mysql`生成`EXPLAIN`分析报告

