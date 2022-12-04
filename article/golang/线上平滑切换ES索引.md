## 前言

> 哈喽，大家好，我是`asong`，今天与大家聊一聊如何平滑切换线上的`ES`索引。使用过`ES`的朋友们都知道，修改索引真的是一件费时又费力的工作，所以我们应该在创建索引的时候就尽量设计好索引能够满足需求，当然这几乎是不可能的，毕竟存在着万恶的产品经理，所以掌握"平滑切换线上的`ES`索引"就很必要，接下来我们就来看一看如何实现！





## 前置条件

能够平滑切换线上的`ES`索引需要有两个先决条件，只有满足了这两个条件才能去执行接下来的平滑切换操作，否则一切操作都是白费。



### 前置条件之使用别名访问索引

重建索引的问题是必须更新应用中的索引名称，索引别名就是用来解决这个问题的。索引别名就像一个快捷方式或软连接，可以指向一个或多个索引，也可以给任何一个需要索引名的API来使用。*别名* 带给我们极大的灵活性，允许我们做下面这些：

- 在运行的集群中可以无缝的从一个索引切换到另一个索引
- 给多个索引分组
- 给索引的一个子集创建 `视图`

索引与索引别名的关系，我们画个图来说一下：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-04-02%20%E4%B8%8B%E5%8D%889.12.01.png)

上图中`user_index`就是索引别名，`user_index_v1`、`user_index_v2`、`user_index_v3`分别是三个索引，这里索引别名`user_index`与`user_index_v1`进行了关联，所以我们搜索的时候使用索引别名，也就是去索引`user_index_v1`上查询。假设现在我们不想使用索引`user_index_v1`了，想使用索引`user_index_v2`，那么直接使用`_aliases`操作执行原子操作(后面介绍具体使用)，将索引别名`user_index`与索引`user_index_v2`进行关联，现在使用索引别名`user_index`搜索的就是索引`user_index_v2`的数据了。



### 前置条件之足够空间

既然要重建`ES`索引，就一定保证你有足够的空间存储数据，可以使用如下指令查看`ES`每个节点的可用磁盘空间：

```go
curl http://localhost:9200/_cat/allocation\?v
```

获得结果如下：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-04-02%20%E4%B8%8B%E5%8D%889.32.20.png)

## 如何平滑切换

因为大家使用的`ES`场景不同，所以平滑切换的步骤会稍有偏差，但是都离不开这几个步骤：

1. 创建新索引
2. 同步数据/数据迁移到新索引
3. 切换索引

先介绍一下数据迁移和切换索引使用什么指令操作：

- 数据迁移

使用`ES`中提供的`reindex api`就可以将数据`copy`到新索引中，比如：

```curl
curl --location --request POST 'http://localhost:9200/_reindex' \
--header 'Content-Type: application/json' \
--data-raw '{
  "conflicts": "proceed",
  "source": {
    "index": "user_index_v1"
  },
  "dest": {
    "index": "user_index_v2",
    "op_type": "create",
    "version_type": "external"
  }
}'
```

介绍一下上面几个字段的意义：

- `"source":{"index": "user_index_v1"}`：这里代表我们要迁移数据的源索引；
- `"dest":{"index": "user_index_v2"}`：这里代表我们要迁移的目标索引；
- `"conflicts": "proceed"`：默认情况下，版本冲突会导致`_reindex`操作终止，可以设置这个字段使该请求遇到冲突时不会终止，而是统计冲突数量；
- `"version_type": "extrenal"`：这个字段介绍起来比较复杂，且听我细细道来。`_reindex`指令会生成源索引的快照，它的目标索引必须是一个不同的索引[新索引]，以便避免版本冲突。如果不设置`version_type`字段，默认为`internal`，`ES`会直接将文档转存储到目标索引中(`dest index`)，直接覆盖任何具有相同类型和`id`的`document`，不会产生版本冲突。如果把`version_type`设置为`extertral`，那么`ES`会从源索引(`source index`)中读取`version`字段，当遇到具有相同类型和`id`的`document`时，只会保留`new version`，即最新的`version`对应的数据。此时可能会有冲突产生，比如当把`op_tpye`设置为`create`，对于产生的冲突现象，返回体中的 `failures` 会携带冲突的数据信息【类似详细的日志可以查看】。
- `op_type`：`op_type` 参数控制着写入数据的冲突处理方式，如果把 `op_type` 设置为 `create`【默认值】，在 `_reindex API` 中，表示写入时只在 `dest index` 中添加不存在的 `doucment`，如果相同的 `document` 已经存在，则会报 `version confilct` 的错误，那么索引操作就会失败。【这种方式与使用 `_create API` 时效果一致】。

更多`_redinx api`使用方法可以移步官方文档学习：https://www.elastic.co/guide/en/elasticsearch/reference/5.6/docs-reindex.html

上面只是举一个简单的例子，具体要在数据迁移中使用哪些参数需要根据场景而定。

什么时候可以选择数据迁移：

1. 当我们新创建的索引只改变了`mapping`结构时，例如：删除字段，更新字段的类型，这种场景就可以直接使用`_reindex`进行数据迁移；
2. 新创建的索引中添加了新字段，但是新的字段都是由老的字段计算得到的，这种情况，也可以使用`_reindex`进行数据迁移，`api`中使用`script`参数，编写你的脚本即可。

注意：使用`_redindex`接口时要注意一个问题，接口会在`reindex`结束后返回，接口超时控制只有`30s`，如果`reindex`时间过长，建议加上`wait_for_completion=false`参数，这样`redindex`就变成异步任务，返回的是`taskID`，查看进度可以通过 `_tasks API` 进行查看。

- 切换索引

`ES`中两种方式管理别名：`_alias`用于单个操作，`_aliases`用于执行多个原子级操作。

因为我们这里要做的是切换索引，主要分为两个步骤：

1. 移除当前索引与索引别名的关联
2. 将新建的索引与索引别名进行关联

所以我们可以选择`_alisases`执行原子操作：

```curl
curl --location --request POST 'http://localhost:9200/_aliases' \
--header 'Content-Type: application/json' \
--data-raw '{
    "actions": [
        {"remove": {"index": "user_index_v2", "alias": "user_index"}},
        { "add": {"index": "user_index_v1",  "alias": "user_index"}}
    ]
}'
```



### 举例子

假设我们有一个`user_index_v1`，他的`mapping`结构如下；

```json
{
    "mappings":{
        "properties":{
            "id":{
                "type":"byte"
            },
            "Name":{
                "type":"text"
            },
            "Age":{
                "type":"byte"
            }
        }
    }
}
```

现在这个`v1`索引中，我们的`id`字段使用的`byte`类型，显然范围是比较小的，随着数据量增多，`id`数值的增大，该字段已经不能满足存储需求了，所以需要把它换成`long`类型，因此可以创建`v2`索引：

```json
{
    "mappings":{
        "properties":{
            "id":{
                "type":"long"
            },
            "Name":{
                "type":"text"
            },
            "Age":{
                "type":"byte"
            }
        }
    }
}
```

现在我们就来考虑一下，如何平滑的进行索引切换。这里假设我们`ES`中数据同步采用的`消息队列`推送完成的，所以在切换索引时要考虑数据损失的问题。

这里我们可以列举几种方案如下：

- 方案一：直接创建`v2`索引，使用`_aliases`切换索引，进行数据迁移，优点是直接切换别名和索引的关联，简单方便，缺点是出现问题回退到旧索引，会有数据损失，直接切换到`v2`索引会导致服务在数据没有迁移完之前不可用。
- 方案二：创建`v2`索引，添加`v2`索引与别名的关联，进行数据迁移，`_alias`操作解除别名和`v2`索引的关联。优点是不会造成服务不可用，缺点是在解除别名和`v1`关联之前，一个别名关联两个索引，单索引操作无法执行，只能搜索，搜索也会出现数据重复，并且也会造成数据损失。
- 方案三：创建`v2`索引，添加`v2`索引与别名的关联，修改代码写入操作使用`v2`索引，搜索操作使用别名索引，进行数据迁移，解除`v1`索引与别名的关联，优点是搜索和写入操作分开了，缺点是回退需要修改代码，并且会出现数据损失，如果`v2`索引不可用了，不能立刻回退索引。
- 方案四：创建`v2`索引，进行数据迁移，然后切换索引；优点是同步数据到v2期间搜索功能正常使用，回退无数据损失；缺点是会造成数据丢失。
- 方案五：创建`v2`索引，添加两个别名索引`read`和`write`，添加别名`read`和`v1`索引、`v2`索引的关联，添加别名`write`和`v2`索引的关联，进行数据迁移，解除别名`read`和`v1`索引的关联；优点是搜索和写入分开了，更新索引时只需要创建新索引，数据同步完成后，解除别名`read`和旧索引关联即可；缺点是数据迁移完成之前，搜索结果会出现重复，回退到旧索引，会有数据损失。

这里总共列举了5种方案，我也不推荐具体使用那个方案比较好，各有利弊，大家可以根据自己的业务场景来进行选择。

这里以选择方案四为例子，给出我的脚本数据，作为样例；

- 创建`user_index_v2`索引：

```bash
#!/bin/bash

url=$1
index=$2


echo `curl --location --request GET ${url}/${index}`

echo `curl --location --request PUT ${url}/${index} \
--header 'Content-Type: application/json' \
--data-raw '{
    "mappings":{
        "properties":{
            "id":{
                "type":"long"
            },
            "Name":{
                "type":"text"
            },
            "Age":{
                "type":"byte"
            }
        }
    }
}'`

echo `curl --location --request GET ${url}/${index}`
```

运行指令：`./create_index.sh http://localhost:9200 user_index_v2`

- 进行数据迁移(数据量比较大时建议分批`and`异步处理)

```sh
#!/bin/bash

url=$1



echo `curl --location --request POST ${url}/'_reindex?wait_for_completion=false' \
--header 'Content-Type: application/json' \
--data-raw '{
  "conflicts": "proceed",
  "source": {
    "index": "user_index_v1"
  },
  "dest": {
    "index": "user_index_v2",
    "op_type": "create",
    "version_type": "external"
  }
}'`

```

运行指令：`./reindex.sh http://localhost:9200` 

- 切换索引

```sh
#!/bin/bash

url=$1
aliasIndex=$2
oldIndex=$3
newIndex=$4

echo `curl --location --request GET ${url}/${aliasIndex}`

echo `curl --location --request POST ${url}/_aliases --header 'Content-Type: application/json' --data-raw '{"actions": [{"remove": {"index": "'$oldIndex'", "alias": "'$aliasIndex'"}},{ "add": {"index": "'$newIndex'",  "alias": "'$aliasIndex'"}}]}'`

echo `curl --location --request GET ${url}/${aliasIndex}`
```

`运行指令：./aliases.sh http://localhost:9200 user_index user_index_v1 user_index_v2`



## 总结

本文例举了几种平滑切换`ES`索引的方案，可以看出修改索引真不是一件容易的事情，要考虑的事情比较多，所以最好在第一次创建索引的时候就多考虑一下以后的使用场景，确定好字段和类型，这样就可以避免重建`ES`索引。当然随着产品的需求变更，重建`ES`索引也是不可避免的，上面几种仅供大家参考，根据自己的场景去选择就好啦。

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

