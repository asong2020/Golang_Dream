哈喽，大家好，我是`asong`。最近在逛`Go`仓库时看到了一个`commit`是关于排序算法的，即`pdqsort`排序算法，`Go`计划将在一个版本中支持该排序算法，下面我们就具体来看一看这个事情；

`commit`地址：https://github.com/golang/go/commit/72e77a7f41bbf45d466119444307fd3ae996e257

<img src="https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-05-06%20%E4%B8%8B%E5%8D%881.37.10.png" alt="截屏2022-05-06 下午1.37.10" style="zoom:67%;" />

该`commit`中介绍了`pqdsort`的测试结果：

- 在所有基准测试中，`pdqsort`未表现出比以前的其它算法慢
- 常见模式中`pdqsort`通常更快（即在排序切片中快10倍）

`pdqsort`实质为一种混合排序算法，在不同情况下切换到不同的排序机制，该实现灵感来自`C++`和`RUST`的实现，是对`C++`标准库算法`introsort`的一种改进，其理想情况下的时间复杂度为 O(n)，最坏情况下的时间复杂度为 O(n* logn)，不需要额外的空间。

`pdqsort`算法的改进在于对常见的情况做了特殊优化，其主要的思想是不断判定目前的序列情况，然后使用不同的方式和路径达到最优解；如果大家想看一下该算法的具体实现，可以查看`https://github.com/zhangyunhao116/pdqsort`中的实践，其实现就是对下面三种情况的不断循环：

- **短序列情况**：对于长度在 [0, MAX_INSERTION] 的输入，使用 insertion sort (插入排序)来进行排序后直接返回，这里的 MAX_INSERTION 我们在 Go 语言下的性能测试，选定为 24。
- **最坏情况，**如果发现改进的 quicksort 效果不佳(limit == 0)，则后续排序都使用 heap sort 来保证最坏情况时间复杂度为 O(n*logn)。
- **正常情况，**对于其他输入，使用改进的 quicksort 来排序

具体的源代码实现可以自行查看，本文就不过多分析了，下面我们来看一下`pdqsort`的demo：

```go
import (
	"fmt"

	"github.com/zhangyunhao116/pdqsort"
)

func main() {
	x := []int{3, 1, 2, 4, 5, 9, 8, 7}
	pdqsort.Slice(x)
	fmt.Printf("sort_result = %v\n", x)
	search_result := pdqsort.Search(x, 4)
	fmt.Printf("search_result = %v\n", search_result)
	is_sort := pdqsort.SliceIsSorted(x)
	fmt.Printf("is_sort = %v\n", is_sort)
}
```

运行结果：

```
sort_result = [1 2 3 4 5 7 8 9]
search_result = 3
is_sort = true
```

对于此次排序算法优化你们有什么想法？快快上手体验一下吧～。



参考链接：

- https://github.com/golang/go/commit/72e77a7f41bbf45d466119444307fd3ae996e257
- https://www.easemob.com/news/8361
- https://github.com/zhangyunhao116/pdqsort
- https://arxiv.org/pdf/2106.05123.pdf



好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)
