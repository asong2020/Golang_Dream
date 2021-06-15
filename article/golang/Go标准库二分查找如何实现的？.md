## 前言

> 哈喽，大家好，我是`asong`。今天与大家分享一下`Go`标准库`sort.Search`是如何实现二分查找的，为什么突然想到分享这个函数呢。起因是这周在刷`leetcode`的一道题时，需要使用到二分查找，最开始是自己实现了一版，后面就想`Go`语言标准库这么全，应该会有这个封装，果然被我找到了，顺便看了一下他是怎么实现的，感觉挺有意思的，就赶紧来分享一下。
>
> 小声逼逼一下，最近有读者反馈说最近发文不够频繁，在这里同步一下我最近干什么：
>
> 1. 工作太忙了
> 2. 正在写一个本地缓存，后面会分享出来
> 3. 正在筹备我的第一个开源项目，敬请期待





## 什么是二分查找

以下来自百度百科：

> 二分查找也称折半查找（Binary Search），它是一种效率较高的查找方法。但是，折半查找要求线性表必须采用[顺序存储结构](https://baike.baidu.com/item/顺序存储结构/1347176)，而且表中元素按关键字有序排列。

总结一下，使用二分查找必须要符合两个要求：

- 必须采用顺序存储结构
- 必须按关键字大小有序排列

我们来举个例子，就很容易理解了，就拿我和我女朋友的聊天内容做个例子吧：

> 我家宝宝每次买一件新衣服，就会习惯性的问我一句，小松子，来猜一猜我这次花了多少钱？
>
> 小松子：50？
>
> 臭宝：你埋汰我呢？欠打！
>
> 小松子：500？
>
> 臭宝：你当我是富婆呢？能不能有点智商！
>
> 小松子：250？
>
> 臭宝：我感觉你在骂我，可我还没有证据。少啦，少啦！！！
>
> 小松子：哎呀，好难猜呀，290？
>
> 臭宝：啊！你这臭男人，就是在骂我！气死啦，气死啦！多啦，多啦！
>
> 小松子：难道是260？
>
> 臭宝：哎呦，挺厉害呀，竟然猜对了！
>
> 后面对话内容就省略啦.....

这里我只需要`5`次就成功猜出来了，这就是二分查找的思想，每一次猜测，我们都选取一段整数范围的中位数，根据条件帮我们逐步缩小范围，每一次都以让剩下的选择范围缩小一半，效率提高。

- 二分查找的时间复杂度

二分查找每次把搜索区域减少一半，时间复杂度为O(logn)。（n代表集合中元素的个数）。

空间复杂度为O(1)。



## 自己实现一个二分查找

二分算法的实现还是比较简单的，可以分两种方式实现：递归和非递归方式实现，示例代码如下：

非递归方式实现：

> ```go
> // 二分查找非递归实现
> func binarySearch(target int64, nums []int64) int {
>    left := 0
>    right := len(nums)
>    for left <= right {
>       mid := left + (right - left) / 2
>       if target == nums[mid] {
>          return mid
>       }
>       if target > nums[mid] {
>          left = mid + 1
>          continue
>       }
>       if target < nums[mid] {
>          right = mid - 1
>          continue
>       }
>    }
>    return -1
> }
> ```

总体思路很简单，每次我们都选取一段整数范围的中位数，然后根据条件让剩下的选择范围缩小一半。这里有一个要注意的点就是我们在获取中位数时使用的写法是` mid := left + (right - left) / 2`，而不是`mid = （left +right）/2`的写法，这是因为后者有可能会造成位数的溢出，也就会导致结果出问题。

递归方式实现：

```go
func Search(nums []int64, target int64) int {
	return binarySearchRecursive(target, nums, 0, len(nums))
}

func binarySearchRecursive(target int64, nums []int64, left, right int) int {
	if left > right {
		return -1
	}

	mid := left + (right - left) / 2

	if target == nums[mid] {
		return mid
	}
	if nums[mid] < target {
		return binarySearchRecursive(target, nums, mid+1, right)
	}
	if nums[mid] > target {

		return binarySearchRecursive(target, nums, left, mid-1)
	}
	return -1

}
```



## Go标准库是如何实现二分查找的？

我们先看一下标准库中的代码实现：

```go
func Search(n int, f func(int) bool) int {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if !f(h) {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}
```

初一看，与我们上面的实现一点也不同呢，我们来分析一下。

入参`n`就是代表要查找序列的长度，入参`f`就是我们自定义的条件。这段代码很短，大概思路就是：

- 定义好这段序列的开始、结尾的位置
- 使用位移操作获取中位数，这样能更好的避免溢出
- 然后根据我们传入的条件判断是否符合条件，逐渐缩小范围

这段代码与我实现的不同在于，它并不是在用户传入的比较函数`f`返回`true`就结束查找，而是继续在当前`[i, j)`区间的前半段查找，并且，当`f`为`false`时，也不比较当前元素与要查找的元素的大小关系，而是直接在后半段查找。所以`for`循环退出的唯一条件就是`i>=j`，如果我们这样使用，就会出现问题：

```go
func main() {
	nums := []int64{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(sort.Search(len(nums), func(i int) bool{
		return nums[i] == 1
	}))
}
```

运行结果竟然是`7`，而不是`1`，如果我们把条件改成`return nums[i] >=1`，运行结果就对了。这是因为我们传入的条件并不是让用户确认目标条件，这里的思想是让我们逐步缩小范围，通过这个条件，我们每次都可以缩小范围，说的有点饶，就上面的代码举个例子。现在是一个升序数组，我们要找的数值是`1`，我们传入的条件是`return nums[i]>=1`，第一进入函数`Search`，我们获取中的中位数`h`是`3`，当前元素是大于目标数值的，所以我们只能在前半段查找，就是这样不断缩小范围，找到我们最终的那个数值，如果当前序列中没有我们要找的目标数值，那么就会返回我们可以插入的位置，也就是最后一位元素的坐标+1的位置。

这个逻辑说实话，我也是第一次接触，仔细思考了一下，这种实现还是有一些优点的：

- 使用**移位操作**，避免因为`i+j`太大而造成的溢出
- 如果我们查找序列中有多个元素相等时，且我们要找的元素就是这个时，我们总会找到下标最小的那个元素
- 如果我们没找到要找的目标元素时，返回的下标是我们可插入的位置，我们在进行数据插入时，依然可以保证数据的有序

**注意：使用`sort.Search`时，入参条件是根据要查找的序列是升序序列还是降序序列来决定的，如果是升序序列，则传入的条件应该是`>=目标元素值`，如果是降序序列，则传入的条件应该是`<=目标元素值`**



### 解析`int(uint(i+j) >> 1)`这段代码

这里我想单独解析一下这段代码，因为很少见，所以可以当作一个知识点记一下。这里使用到的是移位操作，通过向右移动一位，正好可以得到`/2`的结果。具体什么原因呢，我画了一个图，手工画的，看完你就懂了：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2021-06-14%20%E4%B8%8B%E5%8D%882.17.31.png)

懂了吧，兄弟们！移位实现要比乘除发的效率高很多，我们在平常开发中可以使用这种方式来提升效率。

这里还有一个点就是使用`uint`数据类型，因为`uint`的数据范围是`2^32`即`0`到`4294967295`。使用`uint`可以避免因为`i+j`太大而造成的溢出。



## 总结

好啦，今天的文章到这里就结束了，最后我想说的是，没事大家可以专研一下`Go`标准库中的一些方法是怎样实现的，好的思想我们要借鉴过来。如果我们在面试中使用这种方式写出二分查找，那得到的`offer`的几率不就又增加了嘛～。

**素质三连（分享、点赞、在看）都是笔者持续创作更多优质内容的动力！我是`asong`，我们下期见。**

**创建了一个Golang学习交流群，欢迎各位大佬们踊跃入群，我们一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E6%A0%87%E5%87%86%E8%89%B2%E7%89%88.png)

推荐往期文章：

- [Go语言如何实现可重入锁？](https://mp.weixin.qq.com/s/wBp4k7pJLNeSzyLVhGHLEA)
- [Go语言中new和make你使用哪个来分配内存？](https://mp.weixin.qq.com/s/XJ9O9O4KS3LbZL0jYnJHPg)
- [源码剖析panic与recover，看不懂你打我好了！](https://mp.weixin.qq.com/s/mzSCWI8C_ByIPbb07XYFTQ)
- [空结构体引发的大型打脸现场](https://mp.weixin.qq.com/s/dNeCIwmPei2jEWGF6AuWQw)
- [Leaf—Segment分布式ID生成系统（Golang实现版本）](https://mp.weixin.qq.com/s/wURQFRt2ISz66icW7jbHFw)
- [面试官：两个nil比较结果是什么？](https://mp.weixin.qq.com/s/Dt46eoEXXXZc2ymr67_LVQ)
- [面试官：你能用Go写段代码判断当前系统的存储方式吗?](https://mp.weixin.qq.com/s/ffEsTpO-tyNZFR5navAbdA)
- [如何平滑切换线上Elasticsearch索引](https://mp.weixin.qq.com/s/8VQxK_Xh-bkVoOdMZs4Ujw)