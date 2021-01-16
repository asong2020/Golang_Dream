## 前言

> 上周六马上就下班了，正兴高采烈的想着下班吃什么呢！突然QA找到我，说我们的DB与es无法同步数据了，真是令人头皮发秃，好不容易休一天，啊啊啊，难受呀，没办法，还是赶紧找`bug`吧。下面我就把我这次的`bug`原因分享给大家，避免踩坑～。





## bug原因之`bulk`隐藏错误信息

第一时间，我去看了一下错误日志，竟然没有错误日志，很是神奇，既然这样，那我们就`DEBUG`一下吧，`DEBUG`之前我先贴一段代码：

```go
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

就是上面这段代码，使用`es`的`bulk`批量操作，经过`DEBUG`仍然没有发现任何问题，卧槽！！！没有头绪了，那就看一看`es`源码吧，里面是不是有什么隐藏的点没有注意到。还真被我找到了，我们先看一下`req.Do(ctx)`的实现：

```go
// Do sends the batched requests to Elasticsearch. Note that, when successful,
// you can reuse the BulkService for the next batch as the list of bulk
// requests is cleared on success.
func (s *BulkService) Do(ctx context.Context) (*BulkResponse, error) {
	/**
	...... 省略部分代码
  **/
	// Get response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method:      "POST",
		Path:        path,
		Params:      params,
		Body:        body,
		ContentType: "application/x-ndjson",
		Retrier:     s.retrier,
		Headers:     s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return results
	ret := new(BulkResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}

	// Reset so the request can be reused
	s.Reset()

	return ret, nil
}
```

我只把重要部分代码贴出来，看这一段就好了，我来解释一下：

- 首先构建`Http`请求
- 发送`Http`请求并分析，并解析`response`
- 重置`request`可以重复使用

这里的重点就是`ret := new(BulkResponse)`，`new`了一个`BulkResponse`结构，他的结构如下：

```go
type BulkResponse struct {
	Took   int                            `json:"took,omitempty"`
	Errors bool                           `json:"errors,omitempty"`
	Items  []map[string]*BulkResponseItem `json:"items,omitempty"`
}
// BulkResponseItem is the result of a single bulk request.
type BulkResponseItem struct {
	Index         string        `json:"_index,omitempty"`
	Type          string        `json:"_type,omitempty"`
	Id            string        `json:"_id,omitempty"`
	Version       int64         `json:"_version,omitempty"`
	Result        string        `json:"result,omitempty"`
	Shards        *ShardsInfo   `json:"_shards,omitempty"`
	SeqNo         int64         `json:"_seq_no,omitempty"`
	PrimaryTerm   int64         `json:"_primary_term,omitempty"`
	Status        int           `json:"status,omitempty"`
	ForcedRefresh bool          `json:"forced_refresh,omitempty"`
	Error         *ErrorDetails `json:"error,omitempty"`
	GetResult     *GetResult    `json:"get,omitempty"`
}
```

先来解释一个每个字段的意思：

- `took`：总共耗费了多长时间，单位是毫秒
- `Errors`：如果其中任何子请求失败，该 `errors` 标志被设置为 `true` ，并且在相应的请求报告出错误明细（看下面的Items解释）
- `Items`：这个里就是存储每一个子请求的`response`，这里的`Error`存储的是详细的错误信息

现在我想大家应该知道为什么我们的代码没有报`err`信息了，**`bulk`的每个请求都是独立的执行，因此某个子请求的失败不会对其他子请求的成功与否造成影响，所以其中某一条出现错误我们需要从`BulkResponse`解出来**。现在我们把代码改正确：

```go
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
	res, err := req.Do(ctx)
	if err != nil {
		return err
	}
	// 任何子请求失败，该 `errors` 标志被设置为 `true` ，并且在相应的请求报告出错误明细
	// 所以如果没有出错，说明全部成功了，直接返回即可
	if !res.Errors {
		return nil
	}
	for _, it := range res.Failed() {
		if it.Error == nil {
			continue
		}
		return &elastic.Error{
			Status:  it.Status,
			Details: it.Error,
		}
	}
	return nil
}
```

这里再解释一下`res.Failed`方法，这里会把`items`中`bulk response`带错误的返回，所以在这里面找错误信息就可以了。

至此，这个`bug`原因终于被我找到了，接下来可以看下一个`bug`了，我们先简单总结一下：

**`bulk` API 允许在单个步骤中进行多次 `create` 、 `index` 、 `update` 或 `delete` 请求，每个子请求都是独立执行，因此某个子请求的失败不会对其他子请求的成功与否造成影响。`bulk`的response结构中`Erros`字段，如果其中任何子请求失败，该 `errors` 标志被设置为 `true` ，并且在相应的请求报告出错误明细，`items`字段是一个数组，，这个数组的内容是以请求的顺序列出来的每个请求的结果。所以在使用`bulk`时一定要从`response`中判断是否有`err`。**



## bug原因之数值范围越界

这里完全是自己使用不当造成，但还是想说一说`es`的映射数字类型范围的问题：

数字类型有如下分类:

| 类型         | 说明                                                         |
| ------------ | ------------------------------------------------------------ |
| byte         | 有符号的8位整数, 范围: [-128 ~ 127]                          |
| short        | 有符号的16位整数, 范围: [-32768 ~ 32767]                     |
| integer      | 有符号的32位整数, 范围: [−231−231 ~ 231231-1]                |
| long         | 有符号的64位整数, 范围: [−263−263 ~ 263263-1]                |
| float        | 32位单精度浮点数                                             |
| double       | 64位双精度浮点数                                             |
| half_float   | 16位半精度IEEE 754浮点类型                                   |
| scaled_float | 缩放类型的的浮点数, 比如price字段只需精确到分, 57.34缩放因子为100, 存储结果为5734 |

这里都是有符号类型的，无符号在`es7.10.1`版本才开始支持，有兴趣的同学[戳这里](https://www.elastic.co/guide/en/elasticsearch/reference/7.x/release-notes-7.10.0.html)。

这里把这些数字类型及范围列出来就是方便说我的`bug`原因，这里直接解释一下：

我在DB设置字段的类型是`tinyint unsigned`，`tinyint`是一个字节存储，无符号的话范围是`0-255`，而我在`es`中映射类型选择的是`byte`，范围是`-128~127`，当DB中数值超过这个范围是，在进行同步时就会出现这个问题，这里需要大家注意一下数值范围的问题，不要像我一样，因为这个还排查了好久的`bug`，有些空间没必要省，反正也占不了多少空间。





## 总结

这篇文章就是简单总结一下我在工作中遇到的问题，发表出来就是给大家提个醒，有人踩过的坑，就不要在踩了，浪费时间！！！！

**好啦，这篇文章就到这里啦，素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！**

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**翻译了一份Machinery中文文档，会定期进行维护，有需要的小伙伴们后台回复[machinery]即可获取。**

**我是asong，一名普普通通的程序猿，让gi我一起慢慢变强吧。欢迎各位的关注，我们下期见~~~**

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

