## 前言

> 哈喽，大家好，我是`asong`；最近在看Go语言调度器相关的源码，发现看源码真是个技术活，所以本文就简单总结一下该如何查看`Go`源码，希望对你们有帮助。



## Go源码包括哪些？

以我个人理解，`Go`源码主要分为两部分，一部分是官方提供的标准库，一部分是`Go`语言的底层实现，`Go`语言的所有源码/标准库/编译器都在`src`目录下：https://github.com/golang/go/tree/master/src，想看什么库的源码任君选择；

观看`Go`标准库 and `Go`底层实现的源代码难易度也是不一样的，我们一般也可以先从标准库入手，挑选你感兴趣的模块，把它吃透，有了这个基础后，我们在看`Go`语言底层实现的源代码会稍微轻松一些；下面就针对我个人的一点学习心得分享一下如何查看`Go`源码；



## 查看标准库源代码

标准库的源代码看起来稍容易些，因为标准库也属于上层应用，我们可以借助IDE的帮忙，其在IDE上就可以跳转到源代码包，我们只需要不断来回跳转查看各个函数实现做好笔记即可，因为一些源代码设计的比较复杂，大家在看时最好通过画图辅助一下，个人觉得画`UML`是最有助于理解的，能更清晰的理清各个实体的关系；

有些时候只看代码是很难理解的，这时我们使用在线调试辅助我们理解，使用IDE提供的调试器或者`GDB`都可以达到目的，写一个简单的`demo`，断点一打，单步调试走起来，比如你要查看`fmt.Println`的源代码，开局一个小红点，然后就是点点点；

<img src="https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-05-14%20%E4%B8%8B%E5%8D%883.52.48.png"  />



## 查看Go语言底层实现

人都是会对未知领域充满好奇，当使用一段时间`Go`语言后，就想更深入的搞明白一些事情，例如：Go程序的启动过程是怎样的，`goroutine`是怎么调度的，`map`是怎么实现的等等一些`Go`底层的实现，这种直接依靠IDE跳转追溯代码是办不到的，这些都属于`Go`语言的内部实现，大都在`src`目录下的`runtime`包内实现，其实现了垃圾回收，并发控制， 栈管理以及其他一些 Go 语言的关键特性，在编译`Go`代码为机器代码时也会将其也编译进来，`runtime`就是`Go`程序执行时候使用的库，所以一些`Go`底层原理都在这个包内，我们需要借助一些方式才能查看到`Go`程序执行时的代码，这里分享两种方式：分析汇编代码、dlv调试；

<img src="https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-05-15%20%E4%B8%8A%E5%8D%8810.53.27.png" alt=""  />



### 分析汇编代码

前面我们已经介绍了`Go`语言实现了`runtime`库，我们想看到一些`Go`语言关键字特性对应`runtime`里的那个函数，可以查看汇编代码，`Go`语言的汇编使用的`plan9`，与`x86`汇编差别还是很大，很多朋友都不熟悉`plan9`的汇编，但是要想看懂`Go`源码还是要对`plan9`汇编有一个基本的了解的，这里推荐曹大的文章：[plan9 assembly 完全解析](https://go.xargin.com/docs/assembly/assembly/#%E5%9F%BA%E6%9C%AC%E6%8C%87%E4%BB%A4)，会一点汇编我们就可以看源代码了，比如想在我们想看`make`是怎么初始化`slice`的，这时我们可以先写一个简单的`demo`：

```go
// main.go
import "fmt"

func main() {
	s := make([]int, 10, 20)
	fmt.Println(s)
}
```

有两种方式可以查看汇编代码：

```go
1. go tool compile -S -N -l main.go
2. go build main.go && go tool objdump ./main
```

方式一是将源代码编译成`.o`文件，并输出汇编代码，方式二是反汇编，这里推荐使用方式一，执行方式一命令后，我们可以看到对应的汇编代码如下：

<img src="https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-05-15%20%E4%B8%8A%E5%8D%8811.08.50.png"  />

`s := make([]int, 10, 20)`对应的源代码就是` runtime.makeslice(SB)`，这时候我们就去`runtime`包下找`makeslice`函数，不断追踪下去就可查看源码实现了，可在`runtime/slice.go`中找到：

<img src="https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-05-15%20%E4%B8%8A%E5%8D%8811.12.06.png"  />



### 在线调试

虽然上面的方法可以帮助我们定位到源代码，但是后续的操作全靠`review`还是难于理解的，如果能在线调试跟踪代码可以更好助于我们理解，目前`Go`语言支持`GDB`、`LLDB`、`Delve`调试器，但只有`Delve`是专门为`Go`语言设计开发的调试工具，所以使用`Delve`可以轻松调试`Go`汇编程序，`Delve`的入门文章有很多，这篇就不在介绍`Delve`的详细使用方法，入门大家可以看曹大的文章：https://chai2010.cn/advanced-go-programming-book/ch3-asm/ch3-09-debug.html，本文就使用一个小例子带大家来看一看`dlv`如何调试`Go`源码，大家都知道向一个`nil`的切片追加元素，不会有任何问题，在源码中是怎么实现的呢？接下老我们使用`dlv`调试跟踪一下，先写一个小`demo`：

```go
import "fmt"

func main() {
	var s []int
	s = append(s, 1)
	fmt.Println(s)
}
```

进入命令行包目录，然后输入`dlv debug`进入调试

```shell
$ dlv debug
Type 'help' for list of commands.
(dlv)
```

因为这里我们想看到`append`的内部实现，所以在`append`那行加上断点，执行如下命令：

```shell
(dlv) break main.go:7
Breakpoint 1 set at 0x10aba57 for main.main() ./main.go:7
```

执行`continue`命令，运行到断点处：

```go
(dlv) continue
> main.main() ./main.go:7 (hits goroutine(1):1 total:1) (PC: 0x10aba57)
     2: 
     3: import "fmt"
     4: 
     5: func main() {
     6:         var s []int
=>   7:         s = append(s, 1)
     8:         fmt.Println(s)
     9: }
```

接下来我们执行`disassemble`反汇编命令查看`main`函数对应的汇编代码：

```go
(dlv) disassemble
TEXT main.main(SB) /Users/go/src/asong.cloud/Golang_Dream/code_demo/src_code/main.go
        main.go:5       0x10aba20       4c8d6424e8                      lea r12, ptr [rsp-0x18]
        main.go:5       0x10aba25       4d3b6610                        cmp r12, qword ptr [r14+0x10]
        main.go:5       0x10aba29       0f86f6000000                    jbe 0x10abb25
        main.go:5       0x10aba2f       4881ec98000000                  sub rsp, 0x98
        main.go:5       0x10aba36       4889ac2490000000                mov qword ptr [rsp+0x90], rbp
        main.go:5       0x10aba3e       488dac2490000000                lea rbp, ptr [rsp+0x90]
        main.go:6       0x10aba46       48c744246000000000              mov qword ptr [rsp+0x60], 0x0
        main.go:6       0x10aba4f       440f117c2468                    movups xmmword ptr [rsp+0x68], xmm15
        main.go:7       0x10aba55       eb00                            jmp 0x10aba57
=>      main.go:7       0x10aba57*      488d05a2740000                  lea rax, ptr [rip+0x74a2]
        main.go:7       0x10aba5e       31db                            xor ebx, ebx
        main.go:7       0x10aba60       31c9                            xor ecx, ecx
        main.go:7       0x10aba62       4889cf                          mov rdi, rcx
        main.go:7       0x10aba65       be01000000                      mov esi, 0x1
        main.go:7       0x10aba6a       e871c3f9ff                      call $runtime.growslice
        main.go:7       0x10aba6f       488d5301                        lea rdx, ptr [rbx+0x1]
        main.go:7       0x10aba73       eb00                            jmp 0x10aba75
        main.go:7       0x10aba75       48c70001000000                  mov qword ptr [rax], 0x1
        main.go:7       0x10aba7c       4889442460                      mov qword ptr [rsp+0x60], rax
        main.go:7       0x10aba81       4889542468                      mov qword ptr [rsp+0x68], rdx
        main.go:7       0x10aba86       48894c2470                      mov qword ptr [rsp+0x70], rcx
        main.go:8       0x10aba8b       440f117c2450                    movups xmmword ptr [rsp+0x50], xmm15
        main.go:8       0x10aba91       488d542450                      lea rdx, ptr [rsp+0x50]
        main.go:8       0x10aba96       4889542448                      mov qword ptr [rsp+0x48], rdx
        main.go:8       0x10aba9b       488b442460                      mov rax, qword ptr [rsp+0x60]
        main.go:8       0x10abaa0       488b5c2468                      mov rbx, qword ptr [rsp+0x68]
        main.go:8       0x10abaa5       488b4c2470                      mov rcx, qword ptr [rsp+0x70]
        main.go:8       0x10abaaa       e8f1dff5ff                      call $runtime.convTslice
        main.go:8       0x10abaaf       4889442440                      mov qword ptr [rsp+0x40], rax
        main.go:8       0x10abab4       488b542448                      mov rdx, qword ptr [rsp+0x48]
        main.go:8       0x10abab9       8402                            test byte ptr [rdx], al
        main.go:8       0x10ababb       488d35be640000                  lea rsi, ptr [rip+0x64be]
        main.go:8       0x10abac2       488932                          mov qword ptr [rdx], rsi
        main.go:8       0x10abac5       488d7a08                        lea rdi, ptr [rdx+0x8]
        main.go:8       0x10abac9       833d30540d0000                  cmp dword ptr [runtime.writeBarrier], 0x0
        main.go:8       0x10abad0       7402                            jz 0x10abad4
        main.go:8       0x10abad2       eb06                            jmp 0x10abada
        main.go:8       0x10abad4       48894208                        mov qword ptr [rdx+0x8], rax
        main.go:8       0x10abad8       eb08                            jmp 0x10abae2
        main.go:8       0x10abada       e8213ffbff                      call $runtime.gcWriteBarrier
        main.go:8       0x10abadf       90                              nop
        main.go:8       0x10abae0       eb00                            jmp 0x10abae2
        main.go:8       0x10abae2       488b442448                      mov rax, qword ptr [rsp+0x48]
        main.go:8       0x10abae7       8400                            test byte ptr [rax], al
        main.go:8       0x10abae9       eb00                            jmp 0x10abaeb
        main.go:8       0x10abaeb       4889442478                      mov qword ptr [rsp+0x78], rax
        main.go:8       0x10abaf0       48c784248000000001000000        mov qword ptr [rsp+0x80], 0x1
        main.go:8       0x10abafc       48c784248800000001000000        mov qword ptr [rsp+0x88], 0x1
        main.go:8       0x10abb08       bb01000000                      mov ebx, 0x1
        main.go:8       0x10abb0d       4889d9                          mov rcx, rbx
        main.go:8       0x10abb10       e8aba8ffff                      call $fmt.Println
        main.go:9       0x10abb15       488bac2490000000                mov rbp, qword ptr [rsp+0x90]
        main.go:9       0x10abb1d       4881c498000000                  add rsp, 0x98
        main.go:9       0x10abb24       c3                              ret
        main.go:5       0x10abb25       e8f61efbff                      call $runtime.morestack_noctxt
        .:0             0x10abb2a       e9f1feffff                      jmp $main.main
```

从以上内容我们看到调用了`runtime.growslice`方法，我们在这里加一个断点：

```shell
(dlv) break runtime.growslice
Breakpoint 2 set at 0x1047dea for runtime.growslice() /usr/local/opt/go/libexec/src/runtime/slice.go:162
```

之后我们再次执行`continue`执行到该断点处：

```shell
(dlv) continue
> runtime.growslice() /usr/local/opt/go/libexec/src/runtime/slice.go:162 (hits goroutine(1):1 total:1) (PC: 0x1047dea)
Warning: debugging optimized function
   157: // NOT to the new requested capacity.
   158: // This is for codegen convenience. The old slice's length is used immediately
   159: // to calculate where to write new values during an append.
   160: // TODO: When the old backend is gone, reconsider this decision.
   161: // The SSA backend might prefer the new length or to return only ptr/cap and save stack space.
=> 162: func growslice(et *_type, old slice, cap int) slice {
   163:         if raceenabled {
   164:                 callerpc := getcallerpc()
   165:                 racereadrangepc(old.array, uintptr(old.len*int(et.size)), callerpc, funcPC(growslice))
   166:         }
   167:         if msanenabled {
```

之后就是不断的单步调试可以看出来切片的扩容策略；到这里大家也就明白了为啥向`nil`的切片追加数据不会有问题了，因为在容量不够时会调用`growslice`函数进行扩容，具体扩容规则大家可以继续追踪，打脸网上那些瞎写的文章。

上文我们介绍调试汇编的一个基本流程，下面在介绍两个我在看源代码时经常使用的命令；

- **goroutines**命令：通过`goroutines`命令（简写grs），我们可以查看所`goroutine`，通过`goroutine (alias: gr)`命令可以查看当前的`gourtine`：

```shell
(dlv) grs
* Goroutine 1 - User: ./main.go:7 main.main (0x10aba6f) (thread 218565)
  Goroutine 2 - User: /usr/local/opt/go/libexec/src/runtime/proc.go:367 runtime.gopark (0x1035232) [force gc (idle)]
  Goroutine 3 - User: /usr/local/opt/go/libexec/src/runtime/proc.go:367 runtime.gopark (0x1035232) [GC sweep wait]
  Goroutine 4 - User: /usr/local/opt/go/libexec/src/runtime/proc.go:367 runtime.gopark (0x1035232) [GC scavenge wait]
  Goroutine 5 - User: /usr/local/opt/go/libexec/src/runtime/proc.go:367 runtime.gopark (0x1035232) [finalizer wait]
```

- `stack`命令：通过`stack`命令（简写bt），我们可查看当前函数调用栈信息：

```shell
(dlv) bt
0  0x0000000001047e15 in runtime.growslice
   at /usr/local/opt/go/libexec/src/runtime/slice.go:183
1  0x00000000010aba6f in main.main
   at ./main.go:7
2  0x0000000001034e13 in runtime.main
   at /usr/local/opt/go/libexec/src/runtime/proc.go:255
3  0x000000000105f9c1 in runtime.goexit
   at /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1581
```

- `regs`命令：通过`regs`命令可以查看全部的寄存器状态，可以通过单步执行来观察寄存器的变化：

```go
(dlv) regs
   Rip = 0x0000000001047e15
   Rsp = 0x000000c00010de68
   Rax = 0x00000000010b2f00
   Rbx = 0x0000000000000000
   Rcx = 0x0000000000000000
   Rdx = 0x0000000000000008
   Rsi = 0x0000000000000001
   Rdi = 0x0000000000000000
   Rbp = 0x000000c00010ded0
    R8 = 0x0000000000000000
    R9 = 0x0000000000000008
   R10 = 0x0000000001088c40
   R11 = 0x0000000000000246
   R12 = 0x000000c00010df60
   R13 = 0x0000000000000000
   R14 = 0x000000c0000001a0
   R15 = 0x00000000000000c8
Rflags = 0x0000000000000202     [IF IOPL=0]
    Cs = 0x000000000000002b
    Fs = 0x0000000000000000
    Gs = 0x0000000000000000
```

- `locals`命令：通过`locals`命令，可以查看当前函数所有变量值：

```shell
(dlv) locals
newcap = 1
doublecap = 0
```



## 总结

看源代码的过程是没有捷径可走的，如果说有，那就是可以先看一些大佬输出的底层原理的文章，然后参照其文章一步步入门源码阅读，最终还是要自己去克服这个困难，本文介绍了我自己查看源码的一些方式，你是否有更简便的方式呢？欢迎评论区分享出来～。

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)