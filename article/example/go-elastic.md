## 前言

> 嗨，every body，我又来了。这回你们没有看错，今天带来的就是`go-elastic`的入门教程，打了好久的预告，今天终于上线了。因为笔主本人也是一个新手，所以也不敢讲太深入的东西，这一篇完全是一个入门级别的教程，适合初学者，所以本文主打通俗易懂教程，帮助那些和我一样刚入门的朋友，快速应用到开发中。所以我打算分两篇文章来讲解，第一篇主要讲一下什么是`ElasticSearch`，第二篇我们来学习一下`olivere/elastic/v7`库，应用到实际项目开发中，我会把我自己常用的轮子分享出来，还没写完，详情会发在我的第二篇博文上。



## 初识`ElasticSearch`

`ElasticSearch`是一个分布式、`RESTful`风格的搜索和数据分析引擎，在国内简称为`ES`；使用`Java`开发的，底层基于`Lucene`是一种全文检索的搜索库，直接使用使用Lucene还是比较麻烦的，Elasticsearch在Lucene的基础上开发了一个强大的搜索引擎。前面说这么多，对于新手的你，其实还是不知道他是干什么的。简单来说，他就是一个搜索引擎，可以快速存储、搜索和分析海量数据。我们常用的`github`、`Stack  Overflow`都采用的`Es`来做的。为了让你们知道他是干什么的，我们先来分析一下他的功能与适用场景。



## 适用场景

从上面的介绍，我们可以分析出`ElasticSearch`的功能：

- 分布式的搜索引擎和数据分析引擎
- 全文检索、结构化检索、数据分析
- 对海量数据进行近实时的处理

我们都知道`Elastic`的底层是开源库`Lucene`。但是，你却没法直接用`Lucene`，必须直接写代码去调用它的接口。`Elastic` 是 `Lucene` 的封装，提供了 `REST API` 的操作接口，开箱即用。我们现在来分析一下这俩的区别：

`Lucene`：是单机应用，只能在单台服务器上使用，最多只能处理单台服务器可以处理的数据量。

`Elasric`：`ES`自动可以将海量数据分散到多台服务器上去存储和检索海联数据的处理：分布式以后，就可以采用大量的服务器去存储和检索数据，自然而然就可以实现海量数据的处理了，近实时；在秒级别对数据进行搜索和分析。

国内外公司使用场景如下：

- 国外
  - 维基百科，类似百度百科，全文检索，高亮，搜索推荐
  - Stack Overflow（国外的程序异常讨论论坛）
  - GitHub（开源代码管理）
  - 电商网站，检索商品
  - 日志数据分析，logstash采集日志，ES进行复杂的数据分析（ELK技术，elasticsearch+logstash+kibana）
  - BI系统，商业智能，Business Intelligence。
- 国内
  - 站内搜索（电商，招聘，门户）
  - IT系统搜索（OA，CRM，ERP）
  - 数据分析（ES热门的一个使用场景）



## 特点

`Es`为什么这么受欢迎，他有什么特点吗？当然是有的，我们一起来看看它的优点。

- Elasticsearch不是什么新技术，主要是将全文检索、数据分析以及分布式技术，合并在了一起，才形成了独一无二的ES；lucene（全文检索），商用的数据分析软件（也是有的），分布式数据库（mycat）
- 数据库的功能面对很多领域是不够用的（事务，还有各种联机事务型的操作）；特殊的功能，比如全文检索，同义词处理，相关度排名，复杂数据分析，海量数据的近实时处理；Elasticsearch作为传统数据库的一个补充，提供了数据库所不不能提供的很多功能
- 可以作为一个大型分布式集群（数百台服务器）技术，处理PB级数据，服务大公司；也可以运行在单机上，服务小公司
- 对用户而言，是开箱即用的，非常简单，作为中小型的应用，直接3分钟部署一下ES，就可以作为生产环境的系统来使用了，数据量不大，操作不是太复杂

现在我们应该知道`ES`是什么了吧，下面我就来安装他，学习怎么使用。

## 安装

`ES`的安装，还是比较简单的。我们只需要下载压缩包，解压缩即可。根据自己的系统选择下载即可。下载地址：https://www.elastic.co/cn/downloads/elasticsearch。

因为我的是`macos`系统，所以以下操作都是基于mac的。

我下载的版本是： elasticsearch-7.9.0。下载好后，需要进行解压。

```shell
$ tar -zxvf elasticsearch-7.9.0-darwin-x86_64.tar.gz /usr/local/Cellar
```

解压好了，我们就可以进入到相应目录，启动es了。

```shel
$ cd /usr/local/elasticsearch-7.9.0/bin
$ ./elasticsearch
```

启动后，访问http://localhost:9200/ ，如果可以正常访问，就说明安装成功了。



安装成功后，我们来尝试使用一下。elasticsearch是以http Restful api的形式提供接口，我们要操作ES，只要调用http接口就行，ES的默认端口是9200, 因此上面例子可以直接通过浏览器访问ES的接口。大家都知道Http Restful api风格的请求动作，主要包括：GET、POST、PUT、DELETE四种，直接通过浏览器访问，发送的是GET请求动作，后面的三种动作，不方便用浏览器模拟，除非你自己写程序调用，但是我们平时测试，又不想写代码，所以建议使用curl命令、或者postman可视化工具发送http请求。

-   url 例子

```shell
$ curl -X GET "localhost:9200/_search"
#结果
{"took":35,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":0,"relation":"eq"},"max_score":null,"hits":[]}}%
```



## Kibana

虽然平常直接使用`Restful api`操作ES还是挺方便的，但是对我这样用惯可视化工具的人来说，还是挺难受的。这不`Kibana`出现了，解决了我的痛苦。我们可以使用Kibana工具操作ES，Kibana以Web后台的形式提供了一个可视化操作ES的系统，支持根据ES数据绘制图表，支持ES查询**语法自动补全**等高级特性。是不是很强大，我们现在就来学习怎么使用。

Kibana也是java开发的，安装启动非常简单，只要下载安装包，解压缩后启动即可。下载地址：https://www.elastic.co/cn/downloads/kibana

我下载的版本是kibana-7.9.1-darwin-x86_64.tar.gz。

```shell
$ tar -zxvf kibana-7.9.1-darwin-x86_64.tar.gz
$ cd /usr/local/Cellar/kibana-7.9.1-darwin-x86_64/bin
$ ./kibana
```

启动后，访问 [http://localhost:5601](http://localhost:5601/)，就可以进入kibana, 首次访问，因为没有数据，会显示如下窗口。

我们平时开发的时候，编写ES查询语句，可以使用Kibana提供的开发工具Console（控制台），调试ES查询有没有问题，Console支持语法补全和语法提示非常方便。

只要进入Kibana后台，点击左侧菜单的Dev Tools就可以进入Console后台。这里我就不截图了，使用还是比较简单的。



## 快速入门

好啦，前面的铺垫，就是为了接下来的使用，我们在学习一门技术时，我的个人习惯是先使用上，然后再去学习其原理，要不会是一头雾水的。就好比我在公司看其他同事的代码，我不是上来就看代码，而是先把项目运行起来，看一看实现了什么功能，哪个功能在代码中怎么实现的，有目的性的学习，才能更好的得到吸收。好啦，不废话啦，开始接下来的学习。



### 1. 存储结构

大家对mysq的存储结构应该是很清楚的，所以咱们在学习ES存储结构时，同时类比mysql，这样理解起来会更透彻。mysql的数据模型由数据库、表、字段、字段类型组成，自然ES也有自己的一套存储结构。

先看一个表格，然后我们在展开学习每一部分。

| ES存储结构    | Mysql存储结构 |
| ------------- | ------------- |
| Index（索引） | 表            |
| 文档          | 行，一行数据  |
| Field（字段） | 表字段        |
| mapping(映射) | 表结构定义    |

#### 1.1 index

`ES`中索引(index)就像`mysql`中的表一样，代表着文档数据的集合，文档就相当于`ES`中存储的一条数据，下面会详细介绍。



#### 1.2 type

type也就是文档类型，不过在`Elasticsearch7.0`以后的版本,已经废弃文档类型了。不过我们还是要知道这个概念的。在`Elasticsearch`老的版本中文档类型，代表一类文档的集合，index(索引)类似mysql的数据库、文档类型类似Mysql的表。既然新的版本文档类型没什么作用了，那么index（索引）就类似mysql的表的概念，ES没有数据库的概念了。



#### 1.3 document

`ES`是面向文档的数据库，文档是`ES`存储的最基本的存储单元，文档蕾丝`mysql`表中的一行数据。其实在`ES`中，文档指的就是一条JSON数据。`ES`中文档使用`JSON`格式存储，因此存储上要比`mysql`灵活的很多，因为`ES`支持任意格式的`json`数据。

举个例子吧：

```json
{
  "_index" : "order",
  "_type" : "_doc",
  "_id" : "1",
  "_version" : 2,
  "_seq_no" : 1,
  "_primary_term" : 1,
  "found" : true,
  "_source" : {
    "id" : 10000,
    "status" : 0,
    "total_price" : 10000,
    "create_time" : "2020-09-06 17:30:22",
    "user" : {
      "id" : 10000,
      "username" : "asong2020",
      "phone" : "888888888",
      "address" : "深圳人才区"
    }
  }
}

```

这个是`kibana`的一条数据。文档中的任何`json`字段都可以作为查询条件。并且文档的`json`格式没有严格限制，可以随意增加，减少字段，甚至每个文档的格式都不一样也可以。

**注意：**这里我特意加粗了，虽然文档格式是没有限制的，可以随便存储数据，但是，我们在实际开发中是不可以这么做的，下一篇具体实战当中，我会进行讲解。我们在实际项目开发中，一个索引只会存储格式相同的数据。

上面我们已经看到了一个文档数据，下面我们来了解一下什么是文档元数据，指的是插入JSON文档的时候，`ES`为这条数据，自动生成的系统字段。
我们常用的元数据如下：
- _index：代表当前`json`文档所属的文档名字
- _type：代表当前`json`文档所属的类型。不过在`es7.0`以后废弃了`type`用法，但是元数据还是可以看到的
- _id：文档唯一`ID`，如果我们没有为文档指定`id`，系统自动生成。
- _source：代表我们插入进入`json`数据
- _version：文档的版本号，每修改一次文档数据，字段就会加1，这个字段新版`es`也给取消了
- _seq_no：文档的编号，替代老的 version字段
- _primary_term：文档所在主分区，这个可以跟seq_no字段搭配实现乐观锁

#### 1.4 Field

文档由多个`json`字段，这个字段跟`mysql`中的表的字段是类似的。`ES`中的字段也是有类型的，常用字段类型有：

- 数值类型(long、integer、short、byte、double、float)
- Date 日期类型
- boolean布尔类型
- Text 支持全文搜索
- Keyword 不支持全文搜索，例如：phone这种数据，用一个整体进行匹配就`ok`了，也不要进行分词处理

- Geo 这里主要用于地理信息检索、多边形区域的表达。



#### 1.5 mapping

`Elasticsearch`的`mapping`类似于`mysql`中的表结构体定义，每个索引都有一个映射的规则，我们可以通过定义索引的映射规则，提前定义好文档的`json`结构和字段类型，如果没有定义索引的映射规则，`ElasticSearch`会在写入数据的时候，根据我们写入的数据字段推测出对应的字段类型，相当于自动定义索引的映射规则。

**注意：**`ES`的自动映射是很方便的，但是实际业务中，对于关键字段类型，我们都是通常预先定义好，这样可以避免`ES`自动生成的字段类型不是你想要的类型。

### 2. `ES`查询

在使用`ES`时，查询是我们经常使用的。所以我们来主要讲解一下查询。

来看一下查询的基本语法结构：

```json
GET /{索引名}/_search
{
	"from" : 0,  // 搜索结果的开始位置
  	"size" : 10, // 分页大小，也就是一次返回多少数据
  	"_source" :[ ...需要返回的字段数组... ],
	"query" : { ...query子句... },
	"aggs" : { ..aggs子句..  },
	"sort" : { ..sort子句..  }
}
```

让我们来依次解释一下每部分：

先看一下`URI`部分，{索引名}是我们要搜索的索引，可以放置多个索引，使用逗号进行分隔，比如：

```
GET /_order_demo1,_order_demo2/_search
GET /_order*/_search # 按前缀匹配索引名
```

查询结果：

```json
{
  "took" : 0,
  "timed_out" : false,
  "_shards" : {
    "total" : 1,
    "successful" : 1,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : {
      "value" : 1,
      "relation" : "eq"
    },
    "max_score" : 1.0,
    "hits" : [
      {
        "_index" : "order",
        "_type" : "_doc",
        "_id" : "1",
        "_score" : 1.0,
        "_source" : {
          "id" : 10000,
          "status" : 0,
          "total_price" : 10000,
          "create_time" : "2020-09-06 17:30:22",
          "user" : {
            "id" : 10000,
            "username" : "asong2020",
            "phone" : "888888888",
            "address" : "深圳人才区"
          }
        }
      }
    ]
  }
}
```

接下来我们一起看一看body中的查询条件：

- `ES`查询分页：通过`from`和`size`参数设置，相当于`MYSQL`的`limit`和`offset`结构
- `query`：主要编写类似SQL的Where语句，支持布尔查询（and/or）、IN、全文搜索、模糊匹配、范围查询（大于小于）
- `aggs`：主要用来编写统计分析语句，类似SQL的group by语句
- `sort`：用来设置排序条件，类似SQL的order by语句
- `source`：用于设置查询结果返回什么字段，相当于`select`语句后面指定字段



#### 2.1 几种查询语法

- 匹配单个字段

通过match实现全文索引，全文搜索是`ES`的关键特性之一，我们平时使用搜索一些文本、字符串是否包含指定的关键词，但是如果两篇文章，都包含我们的关键词，具体那篇文章内容的相关度更高？ 这个SQL的like语句是做不到的，更别说like语句的性能问题了。

ES通过分词处理、相关度计算可以解决这个问题，ES内置了一些相关度算法，例如：TF/IDF算法，大体上思想就是，如果一个关键词在一篇文章出现的频率高，并且在其他文章中出现的少，那说明这个关键词与这篇文章的相关度很高。分词就是为了提取搜索关键词，理解搜索的意图。就好像我们平常使用谷歌搜索的时候，输入的内容可能很长，但不是每个字都对搜索有帮助，所以可以通过粉刺算法，我们输入的搜索关键词，会进一步分解成多个关键词。这里具体的分词算法我就不详细讲解了，有需要的去官方文档看一看更详细的介绍吧。

我们先来看一看匹配单个字段的使用方法：

```shell
GET /{索引名}/_search
{
  "query": {
    "match": {
      "{FIELD}": "{TEXT}"
    }
  }
}
```

说明：

- {FIELD}  就是我们需要匹配的字段名
- {TEXT} 就是我们需要匹配的内容



- 精确匹配单个字段

当我们需要根据手机号、用户名来搜索一个用户信息时，这就需要使用精确匹配了。可以使用`term`实现精确匹配语法：

```shell
GET /{索引名}/_search
{
  "query": {
    "term": {
      "{FIELD}": "{VALUE}"
    }
  }
}
```

说明：

- {FIELD} - 就是我们需要匹配的字段名
- {VALUE} - 就是我们需要匹配的内容，除了TEXT类型字段以外的任意类型。



- 多值匹配

多值匹配，也就是想`mysql`中的in语句一样，一个字段包含给定数组中的任意一个值匹配。上文使用`term`实现单值精确匹配，同理`terms`就可以实现多值匹配。

```shell
GET /{索引名}/_search
{
  "query": {
    "terms": {
      "{FIELD}": [
        "{VALUE1}",
        "{VALUE2}"
      ]
    }
  }
}
```

说明：

- {FIELD} - 就是我们需要匹配的字段名
- {VALUE1}, {VALUE2} .... {VALUE N} - 就是我们需要匹配的内容，除了TEXT类型字段以外的任意类型。



- 范围查询

我们想通过范围来确实查询数据，这时应该怎么做呢？不要慌，当然有办法了，使用`range`就可以实现范围查询，相当于SQL语句的>，>=，<，<=表达式

```shell
GET /{索引名}/_search
{
  "query": {
    "range": {
      "{FIELD}": {
        "gte": 100, 
        "lte": 200
      }
    }
  }
}
```

说明：

- {FIELD} - 字段名
- gte范围参数 - 等价于>=
- lte范围参数 - 等价于 <=
- 范围参数可以只写一个，例如：仅保留 "gte": 100， 则代表 FIELD字段 >= 100

范围参数有如下：

- **gt** - 大于 （ > ）
- **gte** - 大于且等于 （ >= ）
- **lt** - 小于 （ < ）
- **lte** - 小于且等于 （ <= ）



- bool组合查询

前面的查询都是设置单个字段的查询条件，实际项目中这么应用是很少的，基本都是多个字段的查询条件，所以接下来我们就来一起学习一下组合多个字段的查询条件。

我们先来看一下`bool`查询的基本语法结构：

```shell
GET /{索引名}/_search
{
  "query": {
    "bool": { // bool查询
      "must": [], // must条件，类似SQL中的and, 代表必须匹配条件
      "must_not": [], // must_not条件，跟must相反，必须不匹配条件
      "should": [] // should条件，类似SQL中or, 代表匹配其中一个条件
    }
  }
}
```

接下来分析一下每个条件：

- must条件：类似SQL的and，代表必须匹配的条件。
- must_not条件：跟must作用刚好相反，相当于`sql`语句中的 `!=`
- should条件：类似SQL中的 or， 只要匹配其中一个条件即可

#### 2.2 排序

假设我们现在要查询订单列表，那么返回符合条件的列表肯定不会是无序的，一般都是按照时间进行排序的，所以我们就要使用到了排序语句。`ES`的默认排序是根据相关性分数排序，如果我们想根据查询结果中的指定字段排序，需要使用`sort Processors`处理。

```shell
GET /{索引名}/_search
{
  "query": {
    ...查询条件....
  },
  "sort": [
    {
      "{Field1}": { // 排序字段1
        "order": "desc" // 排序方向，asc或者desc, 升序和降序
      }
    },
    {
      "{Field2}": { // 排序字段2
        "order": "desc" // 排序方向，asc或者desc, 升序和降序
      }
    }
    ....多个排序字段.....
  ]
}
```

sort子句支持多个字段排序，类似SQL的order by。



#### 2.3 聚合查询

`ES`中的聚合查询，类似`SQL`的SUM/AVG/COUNT/GROUP BY分组查询，主要用于统计分析场景。

我们先来看一看什么是聚合查询：

ES聚合查询类似SQL的GROUP by，一般统计分析主要分为两个步骤：

- 分组
- 组内聚合

对查询的数据首先进行一轮分组，可以设置分组条件，例如：新生入学，把所有的学生按专业分班，这个分班的过程就是对学生进行了分组。

组内聚合，就是对组内的数据进行统计，例如：计算总数、求平均值等等，接上面的例子，学生都按专业分班了，那么就可以统计每个班的学生总数， 这个统计每个班学生总数的计算，就是组内聚合计算。

知道了什么是聚合，下面我们就来看其中几个重要关键字：

- 桶：桶的就是一组数据的集合，对数据分组后，得到一组组的数据，就是一个个的桶。ES中桶聚合，指的就是先对数据进行分组。
- 指标：指标指的是对文档进行统计计算方式，又叫指标聚合。桶内聚合，说的就是先对数据进行分组（分桶），然后对每一个桶内的数据进行指标聚合。说白了就是，前面将数据经过一轮桶聚合，把数据分成一个个的桶之后，我们根据上面计算指标对桶内的数据进行统计。常用的指标有：SUM、COUNT、MAX等统计函数。

了解了真正的概念，我们就可以学习聚合查询的语法了：

```shell
{
  "aggregations" : {
    "<aggregation_name>" : {
        "<aggregation_type>" : {
            <aggregation_body>
        }
        [,"aggregations" : { [<sub_aggregation>]+ } ]? // 嵌套聚合查询，支持多层嵌套
    }
    [,"<aggregation_name_2>" : { ... } ]* // 多个聚合查询，每个聚合查询取不同的名字
  }
}
```

说明：

- **aggregations** - 代表聚合查询语句，可以简写为aggs
- **<aggregation_name>** - 代表一个聚合计算的名字，可以随意命名，因为ES支持一次进行多次统计分析查询，后面需要通过这个名字在查询结果中找到我们想要的计算结果。
- **<aggregation_type>** - 聚合类型，代表我们想要怎么统计数据，主要有两大类聚合类型，桶聚合和指标聚合，这两类聚合又包括多种聚合类型，例如：指标聚合：sum、avg， 桶聚合：terms、Date histogram等等。
- **<aggregation_body>** - 聚合类型的参数，选择不同的聚合类型，有不同的参数。
- **aggregation_name_2** - 代表其他聚合计算的名字，意思就是可以一次进行多种类型的统计。

光看这个查询语法，大家可能是懵逼的，所以我们来举个例子，更好的理解一下：

假设现在`order`索引中，存储了每一笔外卖订单，里面包含了店铺名字这个字段，那我们想要统计每个店铺的订单数量，就需要用到聚合查询。

```
GET /order/_search
{
    "size" : 0, // 设置size=0的意思就是，仅返回聚合查询结果，不返回普通query查询结果。
    "aggs" : { // 简写
        "count_store" : { // 聚合查询名字
            "terms" : { // 聚合类型为，terms，terms是桶聚合的一种，类似SQL的group by的作用，根据字段分组，相同字段值的文档分为一组。
              "field" : "store_name" // terms聚合类型的参数，这里需要设置分组的字段为store_name，根据store_name分组
            }
        }
    }
}

```

这里我们没有明确指定指标聚合函数，默认使用的是Value Count聚合指标统计文档总数。

接下来我们就来介绍一下各个指标聚合函数:

- Value Count：值聚合，主要用于统计文档总数，类似SQL的count函数。

```she
GET /sales/_search?size=0
{
  "aggs": {
    "types_count": { // 聚合查询的名字，随便取个名字
      "value_count": { // 聚合类型为：value_count
        "field": "type" // 计算type这个字段值的总数
      }
    }
  }
}
```



- cardinality

基数聚合，也是用于统计文档的总数，跟Value Count的区别是，基数聚合会去重，不会统计重复的值，类似SQL的count(DISTINCT 字段)用法。

```shell
POST /sales/_search?size=0
{
    "aggs" : {
        "type_count" : { // 聚合查询的名字，随便取一个
            "cardinality" : { // 聚合查询类型为：cardinality
                "field" : "type" // 根据type这个字段统计文档总数
            }
        }
    }
}
```



- avg

求平均值

```shell
POST /exams/_search?size=0
{
  "aggs": {
    "avg_grade": { // 聚合查询名字，随便取一个名字
      "avg": { // 聚合查询类型为: avg
        "field": "grade" // 统计grade字段值的平均值
      }
    }
  }
}
```



- Sum

求和计算

```shell
POST /sales/_search?size=0
{
  "aggs": {
    "hat_prices": { // 聚合查询名字，随便取一个名字
      "sum": { // 聚合类型为：sum
        "field": "price" // 计算price字段值的总和
      }
    }
  }
}
```



- max

求最大值

```shell
POST /sales/_search?size=0
{
  "aggs": {
    "max_price": { // 聚合查询名字,随便取一个名字
      "max": { // 聚合类型为：max
        "field": "price" // 求price字段的最大值
      }
    }
  }
}
```

- min

求最小值

```shel
POST /sales/_search?size=0
{
  "aggs": {
    "min_price": { // 聚合查询名字，随便取一个
      "min": { // 聚合类型为: min
        "field": "price" // 求price字段值的最小值
      }
    }
  }
}
```





## 总结

好啦，这一篇到这里就结束了，一些基本概念以及基础的使用方法都介绍了一遍，因为`ES`知识点就是比较多的，这里只是介绍了一下入门级别的使用方法，适合新手，所以就没有讲更深入的东西，下一篇我将从代码入手，来讲解`ES`在代码中的使用，会比这一篇有意思很多，因为讲概念嘛，很枯燥的，下篇，我们敬请期待呦！！！

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**我是asong，一名普普通通的程序猿，让我一起慢慢变强吧。我自己建了一个`golang`交流群，有需要的小伙伴加我`vx`,我拉你入群。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [面试官：go中for-range使用过吗？这几个问题你能解释一下原因吗](https://mp.weixin.qq.com/s/G7z80u83LTgLyfHgzgrd9g)

- [学会wire依赖注入、cron定时任务其实就这么简单！](https://mp.weixin.qq.com/s/qmbCmwZGmqKIZDlNs_a3Vw)

- [听说你还不会jwt和swagger-饭我都不吃了带着实践项目我就来了](https://mp.weixin.qq.com/s/z-PGZE84STccvfkf8ehTgA)
- [掌握这些Go语言特性，你的水平将提高N个档次(二)](https://mp.weixin.qq.com/s/7yyo83SzgQbEB7QWGY7k-w)
- [go实现多人聊天室，在这里你想聊什么都可以的啦！！！](https://mp.weixin.qq.com/s/H7F85CncQNdnPsjvGiemtg)
- [grpc实践-学会grpc就是这么简单](https://mp.weixin.qq.com/s/mOkihZEO7uwEAnnRKGdkLA)
- [go标准库rpc实践](https://mp.weixin.qq.com/s/d0xKVe_Cq1WsUGZxIlU8mw)
- [2020最新Gin框架中文文档 asong又捡起来了英语，用心翻译](https://mp.weixin.qq.com/s/vx8A6EEO2mgEMteUZNzkDg)
- [基于gin的几种热加载方式](https://mp.weixin.qq.com/s/CZvjXp3dimU-2hZlvsLfsw)
- [boss: 这小子还不会使用validator库进行数据校验，开了～～～](https://mp.weixin.qq.com/s?__biz=MzIzMDU0MTA3Nw==&mid=2247483829&idx=1&sn=d7cf4f46ea038a68e74a4bf00bbf64a9&scene=19&token=1606435091&lang=zh_CN#wechat_redirect)

