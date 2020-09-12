## 前言

> 哈喽，everybody，这是`go-elastic`学习系列教程第二篇文章。[上一篇](https://mp.weixin.qq.com/s/mV2hnfctQuRLRKpPPT9XRw)我们学习了`ElasticSearch`基础，如果还不懂基础的，可以先看一看上一篇文章，[传送门](https://mp.weixin.qq.com/s/mV2hnfctQuRLRKpPPT9XRw)。这一篇我们开始实战，写了一个小`demo`，带你们轻松入门`ElasticSearch`实战开发，再也不用担心`es`部分的需求开发了。代码已上传[github](https://github.com/asong2020/Golang_Dream/tree/master/code_demo/go-elastic-asong),可自行下载学习。如果能给一个小星星就好啦。好啦，废话不多说，直接开始吧。
>
> github地址：https://github.com/asong2020/Golang_Dream/tree/master/code_demo/go-elastic-asong



## 背景

在开始之前，我先来介绍一下我这个样例的功能：

- 添加用户信息
- 更新用户信息
- 删除用户信息
- 根据电话查询指定用户
- 根据昵称、身份、籍贯查询相关用户（查找相似昵称的用户列表、身份相同的用户列表、同城的用户列表）



## 1. 创建客户端

在进行开发之前，我们需要下载一个`Es`依赖库。

```shell
$  go get -u github.com/olivere/elastic/v7
```

下载好了依赖库，下面我们开始编写代码，首先我们需要创建一个`client`，用于操作`ES`，先看代码，然后在进行讲解：

```go
func NewEsClient(conf *config.ServerConfig) *elastic.Client {
	url := fmt.Sprintf("http://%s:%d", conf.Elastic.Host, conf.Elastic.Port)
	client, err := elastic.NewClient(
		//elastic 服务地址
		elastic.SetURL(url),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		log.Fatalln("Failed to create elastic client")
	}
	return client
}
```

这里创建`client`是使用的`NewClient`这个方法进行实现的，在创建时，可以提供`ES`连接参数。上面列举的不全，下面给大家介绍一下。

- `elastic.SetURL(url)`用来设置`ES`服务地址，如果是本地，就是`127.0.0.1:9200`。支持多个地址，用逗号分隔即可。
- `elastic.SetBasicAuth("user", "secret")`这个是基于http base auth 验证机制的账号密码。
- `elastic.SetGzip(true)`启动`gzip`压缩
- `elastic.SetHealthcheckInterval(10*time.Second)`用来设置监控检查时间间隔
- `elastic.SetMaxRetries(5)`设置请求失败最大重试次数，v7版本以后已被弃用
- `elastic.SetSniff(false)`允许指定弹性是否应该定期检查集群（默认为true）
- `elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags))` 设置错误日志输出
- `elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags))` 设置info日志输出

这面这些参数根据自己的使用进行选择。



## 2. 创建index及mapping

上一步，我们创建了`client`，接下来我们就要创建对应的索引以及`mapping`。根据开始介绍的功能，我们来设计我们的`mapping`结构：

```go
mappingTpl = `{
	"mappings":{
		"properties":{
			"id": 				{ "type": "long" },
			"username": 		{ "type": "keyword" },
			"nickname":			{ "type": "text" },
			"phone":			{ "type": "keyword" },
			"age":				{ "type": "long" },
			"ancestral":		{ "type": "text" },
			"identity":         { "type": "text" },
			"update_time":		{ "type": "long" },
			"create_time":		{ "type": "long" }
			}
		}
	}`
```

索引设计为：`index =asong_golang_dream `。

设计好了index及`mapping`后，我们开始编写代码进行创建：

```go
func NewUserES(client *elastic.Client) *UserES {
	index := fmt.Sprintf("%s_%s", author, project)
	userEs := &UserES{
		client:  client,
		index:   index,
		mapping: mappingTpl,
	}

	userEs.init()

	return userEs
}

func (es *UserES) init() {
	ctx := context.Background()

	exists, err := es.client.IndexExists(es.index).Do(ctx)
	if err != nil {
		fmt.Printf("userEs init exist failed err is %s\n", err)
		return
	}

	if !exists {
		_, err := es.client.CreateIndex(es.index).Body(es.mapping).Do(ctx)
		if err != nil {
			fmt.Printf("userEs init failed err is %s\n", err)
			return
		}
	}
}
```

这里我们首先判断`es`中是否已经存在要创建的索引，不存在，调用`CreateIndex`进行创建。



## 3. 批量添加

完成一切准备工作，我们接下来就该进行数据的增删改查了。目前该索引下是没有数据，我们先来学习批量添加，添加一些数据，方便后面的使用。

这里批量添加使用的是`bulk`API，`bulk`API允许在单个步骤中进行多次`create`、`index`，`update`、`delete`请求。如果你需要索引一个数据流比如日志事件，它可以排队和索引数百或数千批次。`bulk` 与其他的请求体格式稍有不同，如下所示：

```go
{ action: { metadata }}\n
{ request body        }\n
{ action: { metadata }}\n
{ request body        }\n
...
```

这种格式类似一个有效的单行 JSON 文档 *流* ，它通过换行符(`\n`)连接到一起。注意两个要点：

- 每行一定要以换行符(`\n`)结尾， *包括最后一行* 。这些换行符被用作一个标记，可以有效分隔行。
- 这些行不能包含未转义的换行符，因为他们将会对解析造成干扰。这意味着这个 JSON *不* 能使用 pretty 参数打印。

`action/metadata` 行指定 *哪一个文档* 做 *什么操作* 。

`action` 必须是以下选项之一:

**`create`**：如果文档不存在，那么就创建它。

**`index`**：创建一个新文档或者替换一个现有的文档。

**`update`**：部分更新一个文档

**`delete`**：删除一个文档

这里我使用的是`index`，代码实现如下：

```go
func (es *UserES) BatchAdd(ctx context.Context, user []*model.UserEs) error {
	var err error
	for i := 0; i < esRetryLimit; i++ {
		if err = es.batchAdd(ctx, user); err != nil {
			fmt.Println("batch add failed ", err)
			continue
		}
		return err
	}
	return err
}

func (es *UserES) batchAdd(ctx context.Context, user []*model.UserEs) error {
	req := es.client.Bulk().Index(es.index)
	for _, u := range user {
		u.UpdateTime = uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
		u.CreateTime = uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
		doc := elastic.NewBulkIndexRequest().Id(strconv.FormatUint(u.ID, 10)).Doc(u)
		req.Add(doc)
	}
	if req.NumberOfActions() < 0 {
		return nil
	}
	if _, err := req.Do(ctx); err != nil {
		return err
	}
	return nil
}
```

写好了代码，接下来我们就来测试一下，这个程序使用的`gin`框架，`API`：`http://localhost:8080/api/user/create`，运行代码，发送一个请求，测试一下：

```shell
$ curl --location --request POST 'http://localhost:8080/api/user/create' \
--header 'Content-Type: application/json' \
--data-raw '{
"id": 6,
"username": "asong6",
"nickname": "Golang梦工厂",
"phone": "17897875432",
"age": 20,
"ancestral": "吉林省深圳市",
"identity": "工人"
}'
```

返回结果：

```json
{
    "code": 0,
    "msg": "success"
}
```

**注意：这里有一个点需要说一下，这里我加了一个`for`循环是为了做重试机制的，重试机会为3次，超过则返回。**

为了确保我们插入成功，可以验证一下，发送如下请求：

```shell
$ curl --location --request GET 'http://localhost:9200/asong_golang_dream/_search'
```



## 4. 批量更新

上面介绍了`bulk`API，批量更新依然也是采用的这个方法，`action`选项为`update`。实现代码如下：

```go
func (es *UserES) BatchUpdate(ctx context.Context, user []*model.UserEs) error {
	var err error
	for i := 0; i < esRetryLimit; i++ {
		if err = es.batchUpdate(ctx, user); err != nil {
			continue
		}
		return err
	}
	return err
}

func (es *UserES) batchUpdate(ctx context.Context, user []*model.UserEs) error {
	req := es.client.Bulk().Index(es.index)
	for _, u := range user {
		u.UpdateTime = uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
		doc := elastic.NewBulkUpdateRequest().Id(strconv.FormatUint(u.ID, 10)).Doc(u)
		req.Add(doc)
	}

	if req.NumberOfActions() < 0 {
		return nil
	}
	if _, err := req.Do(ctx); err != nil {
		return err
	}
	return nil
}
```

验证一下：

```shell
$ curl --location --request PUT 'http://localhost:8080/api/user/update' \
--header 'Content-Type: application/json' \
--data-raw '{
"id": 1,
"username": "asong",
"nickname": "Golang梦工厂",
"phone": "17888889999",
"age": 21,
"ancestral": "吉林省",
"identity": "工人"
}'
```

结果：

```json
{
    "code": 0,
    "msg": "success"
}
```



## 5. 批量删除

批量删除也是采用的`bulk`API，即`action`选项为`delete`。代码实现如下：

```go
func (es *UserES) BatchDel(ctx context.Context, user []*model.UserEs) error {
	var err error
	for i := 0; i < esRetryLimit; i++ {
		if err = es.batchDel(ctx, user); err != nil {
			continue
		}
		return err
	}
	return err
}

func (es *UserES) batchDel(ctx context.Context, user []*model.UserEs) error {
	req := es.client.Bulk().Index(es.index)
	for _, u := range user {
		doc := elastic.NewBulkDeleteRequest().Id(strconv.FormatUint(u.ID, 10))
		req.Add(doc)
	}

	if req.NumberOfActions() < 0 {
		return nil
	}

	if _, err := req.Do(ctx); err != nil {
		return err
	}
	return nil
}
```

测试一下：

```shell
curl --location --request DELETE 'http://localhost:8080/api/user/delete' \
--header 'Content-Type: application/json' \
--data-raw '{
"id": 1,
"username": "asong",
"nickname": "Golang梦工厂",
"phone": "17888889999",
"age": 21,
"ancestral": "吉林省",
"identity": "工人"
}'
```



## 6. 查询

有了数据，我们根据条件查询我们想要的数据了。这里我使用的是`bool`组合查询，这个查询语法，我在之前的文章也讲解过，不懂得可以先看一下这一篇文章：

https://mp.weixin.qq.com/s/mV2hnfctQuRLRKpPPT9XRw。

我们先看代码吧：

```go
func (r *SearchRequest) ToFilter() *EsSearch {
	var search EsSearch
	if len(r.Nickname) != 0 {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewMatchQuery("nickname", r.Nickname))
	}
	if len(r.Phone) != 0 {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewTermsQuery("phone", r.Phone))
	}
	if len(r.Ancestral) != 0 {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewMatchQuery("ancestral", r.Ancestral))
	}
	if len(r.Identity) != 0 {
		search.ShouldQuery = append(search.ShouldQuery, elastic.NewMatchQuery("identity", r.Identity))
	}

	if search.Sorters == nil {
		search.Sorters = append(search.Sorters, elastic.NewFieldSort("create_time").Desc())
	}

	search.From = (r.Num - 1) * r.Size
	search.Size = r.Size
	return &search
}

func (es *UserES) Search(ctx context.Context, filter *model.EsSearch) ([]*model.UserEs, error) {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(filter.MustQuery...)
	boolQuery.MustNot(filter.MustNotQuery...)
	boolQuery.Should(filter.ShouldQuery...)
	boolQuery.Filter(filter.Filters...)

	// 当should不为空时，保证至少匹配should中的一项
	if len(filter.MustQuery) == 0 && len(filter.MustNotQuery) == 0 && len(filter.ShouldQuery) > 0 {
		boolQuery.MinimumShouldMatch("1")
	}

	service := es.client.Search().Index(es.index).Query(boolQuery).SortBy(filter.Sorters...).From(filter.From).Size(filter.Size)
	resp, err := service.Do(ctx)
	if err != nil {
		return nil, err
	}

	if resp.TotalHits() == 0 {
		return nil, nil
	}
	userES := make([]*model.UserEs, 0)
	for _, e := range resp.Each(reflect.TypeOf(&model.UserEs{})) {
		us := e.(*model.UserEs)
		userES = append(userES, us)
	}
	return userES, nil
}
```

我们查询之前进行了条件绑定，这个条件通过`API`进行设定的，根据条件绑定不同`query`。`phone`是具有唯一性的，所以我们可以采用精确查询，也就是使用`NewTermsQuery`进行绑定。`Nickname`、`Identity`、`Ancestral`这些都属于模糊查询，所以我们可以使用匹配查询，用`NewMatchQuery`进行绑定·。查询的数据我们在根据创建时间进行排序。时间由近到远进行排序。

代码量不是很多，看一篇就能懂了，我接下来测试一下：

```shell
$ curl --location --request POST 'http://localhost:8080/api/user/search' \
--header 'Content-Type: application/json' \
--data-raw '{
    "nickname": "",
    "phone": "",
    "identity": "",
    "ancestral": "吉林省",
    "num": 1,
    "size":10
}'
```

这里进行说明一下，使用`json`来选择不同的条件，需要那个条件就填写`json`就好了。这个测试的查询条件就是查找出籍贯是`吉林省`的用户列表，通过`num`、`size`限制查询数据量，即第一页，数据量为10。

验证结果：

```json
{
    "code": 0,
    "data": [
        {
            "id": 6,
            "username": "asong6",
            "nickname": "Golang梦工厂",
            "phone": "17897875432",
            "age": 20,
            "ancestral": "吉林省吉林市",
            "identity": "工人",
            "update_time": 1599905564941,
            "create_time": 1599905564941
        },
        {
            "id": 2,
            "username": "asong2",
            "nickname": "Golang梦工厂",
            "phone": "17897873456",
            "age": 20,
            "ancestral": "吉林省吉林市",
            "identity": "学生",
            "update_time": 1599905468869,
            "create_time": 1599905468869
        },
        {
            "id": 1,
            "username": "asong1",
            "nickname": "Golang梦工厂",
            "phone": "17897870987",
            "age": 20,
            "ancestral": "吉林省吉林市",
            "identity": "工人",
            "update_time": 1599900090160,
            "create_time": 1599900090160
        }
    ],
    "msg": "success"
}
```

目前我的数据量没有那么大，所以只有三条数据，你们可以自己测试一下，添加更多的数据进行测试。



## 6. 批量查询

在一些场景中，我们需要通过多个`ID`批量查询文档。`es`中提供了一个`multiGet`进行批量查询，不过我这里实现的不是用这个方法。因为用更好的方法可以使用。`multiGet`批量查询的实现是跟`redis`的`pipeline`是一个道理的，缓存所有请求，然后统一进行请求，所以这里只是减少了IO的使用。所以我们可以使用更好的方法，使用`search`查询，它提供了根据`id`查询的方法，这个方法是一次请求，完成所有的查询，更高效，所以推荐大家使用这个方法进行批量查询。

代码实现如下：

```go
// 根据id 批量获取
func (es *UserES) MGet(ctx context.Context, IDS []uint64) ([]*model.UserEs, error) {
	userES := make([]*model.UserEs, 0, len(IDS))
	idStr := make([]string, 0, len(IDS))
	for _, id := range IDS {
		idStr = append(idStr, strconv.FormatUint(id, 10))
	}
	resp, err := es.client.Search(es.index).Query(
		elastic.NewIdsQuery().Ids(idStr...)).Size(len(IDS)).Do(ctx)

	if err != nil {
		return nil, err
	}

	if resp.TotalHits() == 0 {
		return nil, nil
	}
	for _, e := range resp.Each(reflect.TypeOf(&model.UserEs{})) {
		us := e.(*model.UserEs)
		userES = append(userES, us)
	}
	return userES, nil
}
```

好啦，写好了代码我们进行验证一下吧。

```shell
$ curl --location --request GET 'http://localhost:8080/api/user/info?id=1,2,3'
```

验证结果:

```json
{
    "code": 0,
    "data": [
        {
            "id": 1,
            "username": "asong",
            "nickname": "Golang梦工厂",
            "phone": "88889999",
            "age": 18,
            "ancestral": "广东省深圳市",
            "identity": "工人"
        },
        {
            "id": 2,
            "username": "asong1",
            "nickname": "Golang梦工厂",
            "phone": "888809090",
            "age": 20,
            "ancestral": "吉林省吉林市",
            "identity": "学生"
        },
        {
            "id": 3,
            "username": "asong2",
            "nickname": "Golang梦工厂",
            "phone": "88343409090",
            "age": 21,
            "ancestral": "吉林省吉林市",
            "identity": "学生"
        }
    ],
    "msg": "success"
}
```



## 总结

这一篇到这里就结束了。本文通过一个代码样例，学习使用go进行`eslatic`开发，本文没有将所有方法都讲全，只是将我们日常使用的一些方法整理出来，供大家入门使用，也可以修改一下使用到项目中呦，以为我在项目中也是这么使用的。如果上面的代码段没有看懂，可以到我的github上下载源代码进行学习，运行整个项目，通过`api`进行测试。如果觉得有用，给个小星星呗！！！

github地址：https://github.com/asong2020/Golang_Dream/tree/master/code_demo/go-elastic-asong

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**我是asong，一名普普通通的程序猿，让我一起慢慢变强吧。我自己建了一个`golang`交流群，有需要的小伙伴加我`vx`,我拉你入群。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [go-ElasticSearch入门看这一篇就够了(一)](https://mp.weixin.qq.com/s/mV2hnfctQuRLRKpPPT9XRw)

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

