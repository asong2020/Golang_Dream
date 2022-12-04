## 前言

> 哈喽，大家好，我是`asong`。今天与大家聊一个比较冷门的高频面试题，关于切片的，`Go`语言中的切片原生支持并发吗？怎么样，心里有答案了嘛，带着你的思考我们一起来看一看这个知识点。



## 实践检验真理

实践是检验真理的唯一标准，所以当我们遇到一个不确定的问题，直接写demo来验证，因为切片的特点，我们可以分多种情况来验证：

1. 不指定索引，动态扩容并发向切片添加数据

```go
func concurrentAppendSliceNotForceIndex() {
	sl := make([]int, 0)
	wg := sync.WaitGroup{}
	for index := 0; index < 100; index++{
		k := index
		wg.Add(1)
		go func(num int) {
			sl = append(sl, num)
			wg.Done()
		}(k)
	}
	wg.Wait()
	fmt.Printf("final len(sl)=%d cap(sl)=%d\n", len(sl), cap(sl))
}
```

通过打印数据发现每次的结果都不一致，先不急出结论，我们在写其他的demo测试一下；

2. 指定索引，指定容量并发向切片添加数据

```go
func concurrentAppendSliceForceIndex() {
	sl := make([]int, 100)
	wg := sync.WaitGroup{}
	for index := 0; index < 100; index++{
		k := index
		wg.Add(1)
		go func(num int) {
			sl[num] = num
			wg.Done()
		}(k)
	}
	wg.Wait()
	fmt.Printf("final len(sl)=%d cap(sl)=%d\n", len(sl), cap(sl))
}
```

通过结果我们可以发现符合我们的预期，长度和容量都是100，所以说slice支持并发吗？





## slice支持并发吗？

我们都知道切片是对数组的抽象，其底层就是数组，在并发下写数据到相同的索引位会被覆盖，并且切片也有自动扩容的功能，当切片要进行扩容时，就要替换底层的数组，在切换底层数组时，多个`goroutine`是同时运行的，哪个`goroutine`先运行是不确定的，不论哪个`goroutine`先写入内存，肯定就有一次写入会覆盖之前的写入，所以在动态扩容时并发写入数组是不安全的；

所以当别人问你`slice`支持并发时，你就可以这样回答它：

> 当指定索引使用切片时，切片是支持并发读写索引区的数据的，但是索引区的数据在并发时会被覆盖的；当不指定索引切片时，并且切片动态扩容时，并发场景下扩容会被覆盖，所以切片是不支持并发的～。

`github`上著名的`iris`框架也曾遇到过切片动态扩容导致`webscoket`连接数减少的`bug`，最终采用`sync.map`解决了该问题，感兴趣的可以看一下这个`issue`:https://github.com/kataras/iris/pull/1023#event-1777396646；



## 总结

针对上述问题，我们可以多种方法来解决切片并发安全的问题：

1. 加互斥锁
2. 使用`channel`串行化操作
3. 使用`sync.map`代替切片

切片的问题还是比较容易解决，针对不同的场景可以选择不同的方案进行优化，你学会了吗？

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)

