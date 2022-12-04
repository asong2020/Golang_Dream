## 前言
> 哈喽，大家好，我是`asong`。上篇文章：[动手实现一个localcache - 设计篇](https://mp.weixin.qq.com/s/ZtSA3J8HK4QarhrJwBQtXw) 介绍了设计一个本地缓存要思考的点，有读者朋友反馈可以借鉴bigcache的存储设计，可以减少GC压力，这个是我之前没有考虑到的，这种开源的优秀设计值得我们学习，所以在动手之前我阅读了几个优质的本地缓存库，总结了一下各个开源库的优秀设计，本文我们就一起来看一下。

## 高效的并发访问

本地缓存的简单实现可以使用`map[string]interface{}` + `sync.RWMutex`的组合，使用`sync.RWMutex`对读进行了优化，但是当并发量上来以后，还是变成了串行读，等待锁的`goroutine`就会`block`住。为了解决这个问题我们可以进行分桶，每个桶使用一把锁，减少竞争。分桶也可以理解为分片，每一个缓存对象都根据他的`key`做`hash(key)`，然后在进行分片：`hash(key)%N`，N就是要分片的数量；理想情况下，每个请求都平均落在各自分片上，基本无锁竞争。

分片的实现主要考虑两个点：

- `hash`算法的选择，哈希算法的选择要具有如下几个特点：
    - 哈希结果离散率高，也就是随机性高
    - 避免产生多余的内存分配，避免垃圾回收造成的压力
    - 哈希算法运算效率高

- 分片的数量选择，分片并不是越多越好，根据经验，我们的分片数可以选择`N`的`2`次幂，分片时为了提高效率还可以使用位运算代替取余。

开源的本地缓存库中 `bigcache`、`go-cache`、`freecache`都实现了分片功能，`bigcache`的`hash`选择的是`fnv64a`算法、`go-cache`的`hash`选择的是djb2算法、`freechache`选择的是`xxhash`算法。这三种算法都是非加密哈希算法，具体选哪个算法更好呢，需要综合考虑上面那三点，先对比一下运行效率，相同的字符串情况下，对比`benchmark`：

```go
func BenchmarkFnv64a(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		fnv64aSum64("test")
	}
	b.StopTimer()
}

func BenchmarkXxxHash(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		hashFunc([]byte("test"))
	}
	b.StopTimer()
}


func BenchmarkDjb2(b *testing.B) {
	b.ResetTimer()
	max := big.NewInt(0).SetUint64(uint64(math.MaxUint32))
	rnd, err := rand.Int(rand.Reader, max)
	var seed uint32
	if err != nil {
		b.Logf("occur err %s", err.Error())
		seed = insecurerand.Uint32()
	}else {
		seed = uint32(rnd.Uint64())
	}
	for i:=0; i < b.N; i++{
		djb33(seed,"test")
	}
	b.StopTimer()
}

```

运行结果：

```go
goos: darwin
goarch: amd64
pkg: github.com/go-localcache
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkFnv64a-16      360577890                3.387 ns/op           0 B/op          0 allocs/op
BenchmarkXxxHash-16     331682492                3.613 ns/op           0 B/op          0 allocs/op
BenchmarkDjb2-16        334889512                3.530 ns/op           0 B/op          0 allocs/op
```

通过对比结果我们可以观察出来`Fnv64a`算法的运行效率还是很高，接下来我们在对比一下随机性，先随机生成`100000`个字符串，都不相同；

```go
func init() {
	insecurerand.Seed(time.Now().UnixNano())
	for i := 0; i < 100000; i++{
		randString[i] = RandStringRunes(insecurerand.Intn(10))
	}
}
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[insecurerand.Intn(len(letterRunes))]
	}
	return string(b)
}
```

然后我们跑单元测试统计冲突数：

```go
func TestFnv64a(t *testing.T) {
	m := make(map[uint64]struct{})
	conflictCount :=0
	for i := 0; i < len(randString); i++ {
		res := fnv64aSum64(randString[i])
		if _,ok := m[res]; ok{
			conflictCount++
		}else {
			m[res] = struct{}{}
		}
	}
	fmt.Printf("Fnv64a conflict count is %d", conflictCount)
}

func TestXxxHash(t *testing.T) {
	m := make(map[uint64]struct{})
	conflictCount :=0
	for i:=0; i < len(randString); i++{
		res := hashFunc([]byte(randString[i]))
		if _,ok := m[res]; ok{
			conflictCount++
		}else {
			m[res] = struct{}{}
		}
	}
	fmt.Printf("Xxxhash conflict count is %d", conflictCount)
}


func TestDjb2(t *testing.T) {
	max := big.NewInt(0).SetUint64(uint64(math.MaxUint32))
	rnd, err := rand.Int(rand.Reader, max)
	conflictCount := 0
	m := make(map[uint32]struct{})
	var seed uint32
	if err != nil {
		t.Logf("occur err %s", err.Error())
		seed = insecurerand.Uint32()
	}else {
		seed = uint32(rnd.Uint64())
	}
	for i:=0; i < len(randString); i++{
		res := djb33(seed,randString[i])
		if _,ok := m[res]; ok{
			conflictCount++
		}else {
			m[res] = struct{}{}
		}
	}
	fmt.Printf("Djb2 conflict count is %d", conflictCount)
}

```

运行结果：

```go
Fnv64a conflict count is 27651--- PASS: TestFnv64a (0.01s)
Xxxhash conflict count is 27692--- PASS: TestXxxHash (0.01s)
Djb2 conflict count is 39621--- PASS: TestDjb2 (0.01s)
```

综合对比下，使用`fnv64a`算法会更好一些。

## 减少GC

`Go`语言是带垃圾回收器的，`GC`的过程也是很耗时的，所以要真的要做到高性能，如何避免`GC`也是一个重要的思考点。`freecacne`、`bigcache`都号称避免高额`GC`的库，`bigcache`做到避免高额`GC`的设计是基于`Go`语言垃圾回收时对`map`的特殊处理；在`Go1.5`以后，如果map对象中的key和value不包含指针，那么垃圾回收器就会无视他，针对这个点们的`key`、`value`都不使用指针，就可以避免`gc`。`bigcache`使用哈希值作为`key`，然后把缓存数据序列化后放到一个预先分配好的字节数组中，使用`offset`作为`value`，使用预先分配好的切片只会给GC增加了一个额外对象，由于字节切片除了自身对象并不包含其他指针数据，所以GC对于整个对象的标记时间是O(1)的。具体原理还是需要看源码来加深理解，推荐看原作者的文章：https://dev.to/douglasmakey/how-bigcache-avoids-expensive-gc-cycles-and-speeds-up-concurrent-access-in-go-12bb；作者在BigCache的基础上自己写了一个简单版本的cache，然后通过代码来说明上面原理，更通俗易懂。

`freecache`中的做法是自己实现了一个`ringbuffer`结构，通过减少指针的数量以零GC开销实现map，`key`、`value`都保存在`ringbuffer`中，使用索引查找对象。`freecache`与传统的哈希表实现不一样，实现上有一个`slot`的概念，画了一个总结性的图，就不细看源码了：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-12-19%20%E4%B8%8B%E5%8D%884.50.11.png)


### 推荐文章

- https://colobu.com/2019/11/18/how-is-the-bigcache-is-fast/

- https://dev.to/douglasmakey/how-bigcache-avoids-expensive-gc-cycles-and-speeds-up-concurrent-access-in-go-12bb

- https://studygolang.com/articles/27222

- https://blog.csdn.net/chizhenlian/article/details/108435024

## 总结

一个高效的本地缓存中，并发访问、减少`GC`这两个点是最重要的，在动手之前，看了这几个库中的优雅设计，直接推翻了我之前写好的代码，真是没有十全十美的设计，无论怎么设计都会在一些点上有牺牲，这是无法避免的，软件开发的道路上仍然道阻且长。自己实现的代码还在缝缝补补当中，后面完善了后发出来，全靠大家帮忙`CR`了。

**好啦，本文到这里就结束了，我是`asong`，我们下期见。**

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)