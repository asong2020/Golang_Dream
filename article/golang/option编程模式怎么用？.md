## 前言

> 哈喽，大家好，我是`asong`。`option`编程模式大家一定熟知，但是其写法不唯一，主要是形成了两个版本的`option`设计，本文就探讨一下其中的优缺点。



## `option`编程模式的引出

在我们日常开发中，经常在初始化一个对象时需要进行属性配置，比如我们现在要写一个本地缓存库，设计本地缓存结构如下：

```go
type cache struct {
	// hashFunc represents used hash func
	HashFunc HashFunc
	// bucketCount represents the number of segments within a cache instance. value must be a power of two.
	BucketCount uint64
	// bucketMask is bitwise AND applied to the hashVal to find the segment id.
	bucketMask uint64
	// segment is shard
	segments []*segment
	// segment lock
	locks    []sync.RWMutex
	// close cache
	close chan struct{}
}
```

在这个对象中，字段`hashFunc`、`BucketCount`是对外暴露的，但是都不是必填的，可以有默认值，针对这样的配置，因为`Go`语言不支持重载函数，我们就需要多种不同的创建不同配置的缓存对象的方法：

```go
func NewDefaultCache() (*cache,error){}
func NewCache(hashFunc HashFunc, count uint64) (*cache,error) {}
func NewCacheWithHashFunc(hashFunc HashFunc) (*cache,error) {}
func NewCacheWithBucketCount(count uint64) (*cache,error) {}
```

这种方式就要我们提供多种创建方式，以后如果我们要添加配置，就要不断新增创建方法以及在当前方法中添加参数，也会导致`NewCache`方法会越来越长，为了解决这个问题，我们就可以使用配置对象方案：

```go
type Config struct {
  HashFunc HashFunc
  BucketCount uint64
}
```

我们把非必填的选项移动`config`结构体内，创建缓存的对象的方法就可以只提供一个，变成这样：

```go
func DefaultConfig() *Config {}
func NewCache(config *Config) (*cache,error) {}
```

这样虽然可以解决上述的问题，但是也会造成我们在`NewCache`方法内做更多的判空操作，`config`并不是一个必须项，随着参数增多，`NewCache`的逻辑代码也会越来越长，这就引出了`option`编程模式，接下来我们就看一下`option`编程模式的两种实现。



## option编程模式一

使用闭包的方式实现，具体实现：

```go
type Opt func(options *cache)

func NewCache(opts ...Opt) {
	c := &cache{
		close: make(chan struct{}),
	}
	for _, each := range opts {
		each(c)
	}
}

func NewCache(opts ...Opt) (*cache,error){
	c := &cache{
		hashFunc: NewDefaultHashFunc(),
		bucketCount: defaultBucketCount,
		close: make(chan struct{}),
	}
	for _, each := range opts {
		each(c)
	}
  ......
}

func SetShardCount(count uint64) Opt {
	return func(opt *cache) {
		opt.bucketCount = count
	}
}

func main() {
	NewCache(SetShardCount(256))
}

```

这里我们先定义一个类型`Opt`，这就是我们`option`的`func`型态，其参数为`*cache`，这样创建缓存对象的方法是一个可变参数，可以给多个`options`，我们在初始化方法里面先进行默认赋值，然后再通过`for loop`将每一个`options`对缓存参数的配置进行替换，这种实现方式就将默认值或零值封装在`NewCache`中了，新增参数我们也不需要改逻辑代码了。但是这种实现方式需要将缓存对象中的`field`暴露出去，这样就增加了一些风险，其次`client`端也需要了解`Option`的参数是什么意思，才能知道要怎样设置值，为了减少`client`端的理解度，我们可以自己提前封装好`option`函数，例如上面的`SetShardCount`，`client`端直接调用并填值就可以了。



## option编程模式二

这种`option`编程模式是`uber`推荐的，是在第一版本上面的延伸，将所有`options`的值进行封装，并设计一个`Option interface`，我们先看例子：

```go
type options struct {
	hashFunc HashFunc
	bucketCount uint64
}

type Option interface {
	apply(*options)
}

type Bucket struct {
	count uint64
}

func (b Bucket) apply(opts *options) {
	opts.bucketCount = b.count
}

func WithBucketCount(count uint64) Option {
	return Bucket{
		count: count,
	}
}

type Hash struct {
	hashFunc HashFunc
}

func (h Hash) apply(opts *options)  {
	opts.hashFunc = h.hashFunc
}

func WithHashFunc(hashFunc HashFunc) Option {
	return Hash{hashFunc: hashFunc}
}

func NewCache(opts ...Option) (*cache,error){
	o := &options{
		hashFunc: NewDefaultHashFunc(),
		bucketCount: defaultBucketCount,
	}
	for _, each := range opts {
		each.apply(o)
	}
  .....
}

func main() {
	NewCache(WithBucketCount(128))
}
```

这种方式我们使用`Option`接口，该接口保存一个未导出的方法，在未导出的`options`结构上记录选项，这种模式为`client`端提供了更多的灵活性，针对每一个`option`可以做更细的`custom function`设计，更加清晰且不暴露`cache`的结构，也提高了单元测试的覆盖性，缺点是当`cache`结构发生变化时，也要同时维护`option`的结构，维护复杂性升高了。



## 总结

这两种实现方式都很常见，其都有自己的优缺点，采用闭包的实现方式，我们不需要为维护`option`，维护者的编码也大大减少了，但是这种方式需要`export`对象中的`field`，是有安全风险的，其次是`client`端需要了解对象结构中参数的意义，才能写出`option`参数，不过这个可以通过自定义`option`方法来解决；采用接口的实现方式更加灵活，每一个`option`都可以做精细化设计，不需要`export`对象中的`field`，并且很容易进行调试和测试，缺点是需要维护两套结构，当对象结构发生变更时，`option`结构也要变更，增加了代码维护复杂性。

实际应用中，我们可以自由变化，不能直接定义哪一种实现就是好的，凡事都有两面性，适合才是最好的。

你有什么想法，评论区讨论起来～

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**


![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)

