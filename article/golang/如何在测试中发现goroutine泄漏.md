## 前言

> 哈喽，大家好，我是`asong`；
>
> 众所周知，`gorourtine`的设计是`Go`语言并发实现的核心组成部分，易上手，但是也会遭遇各种疑难杂症，其中`goroutine`泄漏就是重症之一，其出现往往需要排查很久，有人说可以使用`pprof`来排查，虽然其可以达到目的，但是这些性能分析工具往往是在出现问题后借助其辅助排查使用的，有没有一款可以防患于未然的工具吗？当然有，`goleak`他来了，其由 `Uber` 团队开源，可以用来检测`goroutine`泄漏，并且可以结合单元测试，可以达到防范于未然的目的，本文我们就一起来看一看`goleak`。



## goroutine泄漏

不知道你们在日常开发中是否有遇到过`goroutine`泄漏，`goroutine`泄漏其实就是`goroutine`阻塞，这些阻塞的`goroutine`会一直存活直到进程终结，他们占用的栈内存一直无法释放，从而导致系统的可用内存会越来越少，直至崩溃！简单总结了几种常见的泄漏原因：

- `Goroutine`内的逻辑进入死循坏，一直占用资源
- `Goroutine`配合`channel`/`mutex`使用时，由于使用不当导致一直被阻塞
- `Goroutine`内的逻辑长时间等待，导致`Goroutine`数量暴增

接下来我们使用`Goroutine`+`channel`的经典组合来展示`goroutine`泄漏；

```go
func GetData() {
	var ch chan struct{}
	go func() {
		<- ch
	}()
}

func main()  {
	defer func() {
		fmt.Println("goroutines: ", runtime.NumGoroutine())
	}()
	GetData()
	time.Sleep(2 * time.Second)
}
```

这个例子是`channel`忘记初始化，无论是读写操作都会造成阻塞，这个方法如果是写单测也是检查不出来问题的：

```go
func TestGetData(t *testing.T) {
	GetData()
}
```

运行结果：

```go
=== RUN   TestGetData
--- PASS: TestGetData (0.00s)
PASS
```

内置测试无法满足，接下来我们引入`goleak`来测试一下。



## goleak

**github地址**：https://github.com/uber-go/goleak

使用`goleak`主要关注两个方法即可：`VerifyNone`、`VerifyTestMain`，`VerifyNone`用于单一测试用例中测试，`VerifyTestMain`可以在`TestMain`中添加，可以减少对测试代码的入侵，举例如下：

使用`VerifyNone`:

```go
func TestGetDataWithGoleak(t *testing.T) {
	defer goleak.VerifyNone(t)
	GetData()
}
```

运行结果：

```shell
=== RUN   TestGetDataWithGoleak
    leaks.go:78: found unexpected goroutines:
        [Goroutine 35 in state chan receive (nil chan), with asong.cloud/Golang_Dream/code_demo/goroutine_oos_detector.GetData.func1 on top of the stack:
        goroutine 35 [chan receive (nil chan)]:
        asong.cloud/Golang_Dream/code_demo/goroutine_oos_detector.GetData.func1()
        	/Users/go/src/asong.cloud/Golang_Dream/code_demo/goroutine_oos_detector/main.go:12 +0x1f
        created by asong.cloud/Golang_Dream/code_demo/goroutine_oos_detector.GetData
        	/Users/go/src/asong.cloud/Golang_Dream/code_demo/goroutine_oos_detector/main.go:11 +0x3c
        ]
--- FAIL: TestGetDataWithGoleak (0.45s)

FAIL

Process finished with the exit code 1
```

通过运行结果看到具体发生`goroutine`泄漏的具体代码段；使用`VerifyNone`会对我们的测试代码有入侵，可以采用`VerifyTestMain`方法可以更快的集成到测试中：

```go
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
```

运行结果：

```go
=== RUN   TestGetData
--- PASS: TestGetData (0.00s)
PASS
goleak: Errors on successful test run: found unexpected goroutines:
[Goroutine 5 in state chan receive (nil chan), with asong.cloud/Golang_Dream/code_demo/goroutine_oos_detector.GetData.func1 on top of the stack:
goroutine 5 [chan receive (nil chan)]:
asong.cloud/Golang_Dream/code_demo/goroutine_oos_detector.GetData.func1()
	/Users/go/src/asong.cloud/Golang_Dream/code_demo/goroutine_oos_detector/main.go:12 +0x1f
created by asong.cloud/Golang_Dream/code_demo/goroutine_oos_detector.GetData
	/Users/go/src/asong.cloud/Golang_Dream/code_demo/goroutine_oos_detector/main.go:11 +0x3c
]

Process finished with the exit code 1
```

`VerifyTestMain`的运行结果与`VerifyNone`有一点不同，`VerifyTestMain`会先报告测试用例执行结果，然后报告泄漏分析，如果测试的用例中有多个`goroutine`泄漏，无法精确定位到发生泄漏的具体test，需要使用如下脚本进一步分析：

```go
# Create a test binary which will be used to run each test individually
$ go test -c -o tests

# Run each test individually, printing "." for successful tests, or the test name
# for failing tests.
$ for test in $(go test -list . | grep -E "^(Test|Example)"); do ./tests -test.run "^$test\$" &>/dev/null && echo -n "." || echo -e "\n$test failed"; done
```

这样会打印出具体哪个测试用例失败。



## goleak实现原理

从`VerifyNone`入口，我们查看源代码，其调用了`Find`方法：

```go
// Find looks for extra goroutines, and returns a descriptive error if
// any are found.
func Find(options ...Option) error {
  // 获取当前goroutine的ID
	cur := stack.Current().ID()

	opts := buildOpts(options...)
	var stacks []stack.Stack
	retry := true
	for i := 0; retry; i++ {
    // 过滤无用的goroutine
		stacks = filterStacks(stack.All(), cur, opts)

		if len(stacks) == 0 {
			return nil
		}
		retry = opts.retry(i)
	}

	return fmt.Errorf("found unexpected goroutines:\n%s", stacks)
}
```

我们在看一下`filterStacks`方法：

```go
// filterStacks will filter any stacks excluded by the given opts.
// filterStacks modifies the passed in stacks slice.
func filterStacks(stacks []stack.Stack, skipID int, opts *opts) []stack.Stack {
	filtered := stacks[:0]
	for _, stack := range stacks {
		// Always skip the running goroutine.
		if stack.ID() == skipID {
			continue
		}
		// Run any default or user-specified filters.
		if opts.filter(stack) {
			continue
		}
		filtered = append(filtered, stack)
	}
	return filtered
}
```

这里主要是过滤掉一些不参与检测的`goroutine stack`，如果没有自定义`filters`，则使用默认的`filters`：

```go
func buildOpts(options ...Option) *opts {
	opts := &opts{
		maxRetries: _defaultRetries,
		maxSleep:   100 * time.Millisecond,
	}
	opts.filters = append(opts.filters,
		isTestStack,
		isSyscallStack,
		isStdLibStack,
		isTraceStack,
	)
	for _, option := range options {
		option.apply(opts)
	}
	return opts
}
```

从这里可以看出，默认检测`20`次，每次默认间隔`100ms`；添加默认`filters`;

总结一下`goleak`的实现原理：

使用`runtime.Stack()`方法获取当前运行的所有`goroutine`的栈信息，默认定义不需要检测的过滤项，默认定义检测次数+检测间隔，不断周期进行检测，最终在多次检查后仍没有找到剩下的`goroutine`则判断没有发生`goroutine`泄漏。



## 总结

本文我们分享了一个可以在测试中发现`goroutine`泄漏的工具，但是其还是需要完备的测试用例支持，这就暴露出测试用例的重要性，朋友们好的工具可以助我们更快的发现问题，但是代码质量还是掌握在我们自己的手中，加油吧，少年们～。

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)



**参考资料**

- https://github.com/uber-go/goleak
- https://segmentfault.com/a/1190000040161853
- https://blog.schwarzeni.com/2021/04/09/goleak-%E7%A0%94%E7%A9%B6/
- https://zhuanlan.zhihu.com/p/361737398