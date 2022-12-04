## 前言

> 哈喽，大家好，我是`asong`。好久没有更新了，最近因工作需要忙着写`python`，`Go`语言我都有些生疏了，但是我不能放弃`Go`语言，该学习还是要学习的，今天与大家聊一聊`Go`语言的函数调用惯例，调用惯例是调用方和被调用方对于参数和返回值传递的约定，Go语言的调用惯例在1.17版本进行了优化，本文我们就看一下两个版本的调用惯例是什么样的吧～。



## 1.17版本前栈传递

在`Go1.17`版本之前，`Go`语言函数调用是通过栈来传递的，我们使用`Go1.12`版本写个例子来看一下：

```go
package main

func Test(a, b int) (int, int) {
	return a + b, a - b
}

func main() {
	Test(10, 20)
}
```

执行`go tool compile -S -N -l main.go`可以看到其汇编指令，我们分两部分来看，先看`main`函数部分：

```go
"".main STEXT size=68 args=0x0 locals=0x28
        0x0000 00000 (main.go:7)        TEXT    "".main(SB), ABIInternal, $40-0
        0x0000 00000 (main.go:7)        MOVQ    (TLS), CX
        0x0009 00009 (main.go:7)        CMPQ    SP, 16(CX)
        0x000d 00013 (main.go:7)        JLS     61
        0x000f 00015 (main.go:7)        SUBQ    $40, SP // 分配40字节栈空间
        0x0013 00019 (main.go:7)        MOVQ    BP, 32(SP) // 基址指针存储到栈上
        0x0018 00024 (main.go:7)        LEAQ    32(SP), BP
        0x001d 00029 (main.go:7)        FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x001d 00029 (main.go:7)        FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x001d 00029 (main.go:7)        FUNCDATA        $3, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x001d 00029 (main.go:8)        PCDATA  $2, $0
        0x001d 00029 (main.go:8)        PCDATA  $0, $0
        0x001d 00029 (main.go:8)        MOVQ    $10, (SP) // 第一个参数压栈
        0x0025 00037 (main.go:8)        MOVQ    $20, 8(SP) // 第二个参数压栈
        0x002e 00046 (main.go:8)        CALL    "".Test(SB) // 调用函数Test 
        0x0033 00051 (main.go:9)        MOVQ    32(SP), BP // Test函数返回后恢复栈基址指针
        0x0038 00056 (main.go:9)        ADDQ    $40, SP // 销毁40字节栈内存
        0x003c 00060 (main.go:9)        RET // 返回
        0x003d 00061 (main.go:9)        NOP
        0x003d 00061 (main.go:7)        PCDATA  $0, $-1
        0x003d 00061 (main.go:7)        PCDATA  $2, $-1
        0x003d 00061 (main.go:7)        CALL    runtime.morestack_noctxt(SB)
        0x0042 00066 (main.go:7)        JMP     0
        0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 2e 48  eH..%....H;a.v.H
        0x0010 83 ec 28 48 89 6c 24 20 48 8d 6c 24 20 48 c7 04  ..(H.l$ H.l$ H..
        0x0020 24 0a 00 00 00 48 c7 44 24 08 14 00 00 00 e8 00  $....H.D$.......
        0x0030 00 00 00 48 8b 6c 24 20 48 83 c4 28 c3 e8 00 00  ...H.l$ H..(....
        0x0040 00 00 eb bc                                      ....
        rel 5+4 t=16 TLS+0
        rel 47+4 t=8 "".Test+0
        rel 62+4 t=8 runtime.morestack_noctxt+0
```

通过上面的汇编指令我们可以分析出，参数`10`、`20`按照从右向左进行压栈，所以第一个参数在栈顶的位置`SP~SP+8`，第二个参数存储在`SP+8 ~ SP+16`，参数准备完毕后就去调用`TEST`函数，对应的汇编指令：`CALL    "".Test(SB)`，对应的汇编指令如下：

```go
"".Test STEXT nosplit size=49 args=0x20 locals=0x0
        0x0000 00000 (main.go:3)        TEXT    "".Test(SB), NOSPLIT|ABIInternal, $0-32
        0x0000 00000 (main.go:3)        FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (main.go:3)        FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (main.go:3)        FUNCDATA        $3, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (main.go:3)        PCDATA  $2, $0
        0x0000 00000 (main.go:3)        PCDATA  $0, $0
        0x0000 00000 (main.go:3)        MOVQ    $0, "".~r2+24(SP)// SP+16 ~ SP+24 存储第一个返回值
        0x0009 00009 (main.go:3)        MOVQ    $0, "".~r3+32(SP)
// SP+24 ~ SP+32 存储第二个返回值
        0x0012 00018 (main.go:4)        MOVQ    "".a+8(SP), AX // 第一个参数放入AX寄存器 AX = 10
        0x0017 00023 (main.go:4)        ADDQ    "".b+16(SP), AX // 第二个参数放入AX寄存器做加法 AX = AX + 20 = 30
        0x001c 00028 (main.go:4)        MOVQ    AX, "".~r2+24(SP)
// AX寄存器中的值在存回栈中：24(SP)
        0x0021 00033 (main.go:4)        MOVQ    "".a+8(SP), AX
// 第一个参数放入AX寄存器 AX = 10
        0x0026 00038 (main.go:4)        SUBQ    "".b+16(SP), AX
// 第二个参数放入AX寄存器做减法 AX = AX - 20 = -10
        0x002b 00043 (main.go:4)        MOVQ    AX, "".~r3+32(SP)
// AX寄存器中的值在存回栈中：32(SP)
        0x0030 00048 (main.go:4)        RET // 函数返回

```

通过以上的汇编指令我们可以得出结论：`Go`语言使用栈传递参数和接收返回值，多个返回值也是通过多分配一些内存来完成的。

这种基于栈传递参数和接收返回值的设计大大降低了实现的复杂度，但是牺牲了函数调用的性能，像`C`语言采用同时使用栈和寄存器传递参数，在性能上是优于`Go`语言的，下面我们就来看一看`Go1.17`引入的寄存器传参。



## 为什么寄存器传参性能优于栈传参

我们都知道`CPU`是一台计算机的运算核心和控制核心，其主要功能是解释计算机指令以及处理计算机软件中的数据，`CPU`的大致内部结构如下：

![图片来自于网络](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-03-20%20%E4%B8%8B%E5%8D%883.07.57.png)

主要由运算器和控制器组成，运算器负责完成算术运算和逻辑运算，寄存器临时保存将要被运算器处理的数据和处理后的结果，回到主题上，寄存器是`CPU`内部组件，而存储一般在外部，`CPU`操作寄存器与读取内存的速度差距是数量级别的，当要进行数据计算时，如果数据处于内存中，`CPU`需要先将数据从内存拷贝到寄存器进行计算，所以对于栈传递参数与接收返回值这种调用规约，每次计算都需要从内存拷贝到寄存器，计算完毕在拷贝回内存，如果使用寄存器传参的话，参数就已经按顺序放在特定寄存器了，这样就减少了内存和寄存器之间的数据拷贝，从而改善了性能，提供程序运行效率。

既然寄存器传参性能高于栈传递参数，为什么所有语言不都使用寄存器传递参数呢？因为不同架构上的寄存器差异不同，所以要支持寄存器传参就要在编译器上进行支持，这要就使编译器变得更加复杂且不易维护，并且寄存器的数量也是有限的，还要考虑超过寄存器数量的参数应该如何传递。



## 1.17基于寄存器传递

`Go`语言在`1.17`版本设计了一套基于寄存器传参的调用规约，目前也只支持`x86`平台，我们也是通过一个简单的例子看一下：

```go
func Test(a, b, c, d int) (int,int,int,int) {
	return a, b, c, d
}

func main()  {
	Test(1, 2, 3 ,4)
}
```

执行`go tool compile -S -N -l main.go`可以看到其汇编指令，我们分两部分来看，先看`main`函数部分：

```go
"".main STEXT size=62 args=0x0 locals=0x28 funcid=0x0
        0x0000 00000 (main.go:7)        TEXT    "".main(SB), ABIInternal, $40-0
        0x0000 00000 (main.go:7)        CMPQ    SP, 16(R14)
        0x0004 00004 (main.go:7)        PCDATA  $0, $-2
        0x0004 00004 (main.go:7)        JLS     55
        0x0006 00006 (main.go:7)        PCDATA  $0, $-1
        0x0006 00006 (main.go:7)        SUBQ    $40, SP// 分配40字节栈空间，基址指针存储到栈上
        0x000a 00010 (main.go:7)        MOVQ    BP, 32(SP)// 基址指针存储到栈上
        0x000f 00015 (main.go:7)        LEAQ    32(SP), BP
        0x0014 00020 (main.go:7)        FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0014 00020 (main.go:7)        FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0014 00020 (main.go:8)        MOVL    $1, AX // 参数1使用AX寄存器传递
        0x0019 00025 (main.go:8)        MOVL    $2, BX // 参数2使用BX寄存器传递
        0x001e 00030 (main.go:8)        MOVL    $3, CX // 参数3使用CX寄存器传递
        0x0023 00035 (main.go:8)        MOVL    $4, DI // 参数4使用DI寄存器传递
        0x0028 00040 (main.go:8)        PCDATA  $1, $0
        0x0028 00040 (main.go:8)        CALL    "".Test(SB) // 调用Test函数
        0x002d 00045 (main.go:9)        MOVQ    32(SP), BP // Test函数返回后恢复栈基址指针
        0x0032 00050 (main.go:9)        ADDQ    $40, SP // 销毁40字节栈内存
        0x0036 00054 (main.go:9)        RET // 返回

```

通过上面的汇编指令我们可以分析出，现在参数已经不是从右向左进行压栈了，参数直接在寄存器上了，参数准备完毕后就去调用`TEST`函数，对应的汇编指令：`CALL    "".Test(SB)`，对应的汇编指令如下：

```go
"".Test STEXT nosplit size=133 args=0x20 locals=0x28 funcid=0x0
        0x0000 00000 (main.go:3)        TEXT    "".Test(SB), NOSPLIT|ABIInternal, $40-32
        0x0000 00000 (main.go:3)        SUBQ    $40, SP
        0x0004 00004 (main.go:3)        MOVQ    BP, 32(SP)
        0x0009 00009 (main.go:3)        LEAQ    32(SP), BP
        0x000e 00014 (main.go:3)        FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x000e 00014 (main.go:3)        FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x000e 00014 (main.go:3)        FUNCDATA        $5, "".Test.arginfo1(SB)
0x000e 00014 (main.go:3)        MOVQ    AX, "".a+48(SP) // 从寄存器AX获取参数 1 放入栈中 48(SP)
0x0013 00019 (main.go:3)        MOVQ    BX, "".b+56(SP) // 从寄存器BX获取参数 2 放入栈中 56(SP)
0x0018 00024 (main.go:3)        MOVQ    CX, "".c+64(SP) // 从寄存器CX获取参数 3 放入栈中 64(SP)
0x001d 00029 (main.go:3)        MOVQ    DI, "".d+72(SP) // 从寄存器DI获取参数 4 放入栈中 72(SP)
        0x0022 00034 (main.go:3)        MOVQ    $0, "".~r4+24(SP)
        0x002b 00043 (main.go:3)        MOVQ    $0, "".~r5+16(SP)
        0x0034 00052 (main.go:3)        MOVQ    $0, "".~r6+8(SP)
        0x003d 00061 (main.go:3)        MOVQ    $0, "".~r7(SP)
        0x0045 00069 (main.go:4)        MOVQ    "".a+48(SP), DX // 以下操作是返回值放到寄存器中返回
        0x004a 00074 (main.go:4)        MOVQ    DX, "".~r4+24(SP)
        0x004f 00079 (main.go:4)        MOVQ    "".b+56(SP), DX
        0x0054 00084 (main.go:4)        MOVQ    DX, "".~r5+16(SP)
        0x0059 00089 (main.go:4)        MOVQ    "".c+64(SP), DX
        0x005e 00094 (main.go:4)        MOVQ    DX, "".~r6+8(SP)
        0x0063 00099 (main.go:4)        MOVQ    "".d+72(SP), DI
        0x0068 00104 (main.go:4)        MOVQ    DI, "".~r7(SP)
        0x006c 00108 (main.go:4)        MOVQ    "".~r4+24(SP), AX
        0x0071 00113 (main.go:4)        MOVQ    "".~r5+16(SP), BX
        0x0076 00118 (main.go:4)        MOVQ    "".~r6+8(SP), CX
        0x007b 00123 (main.go:4)        MOVQ    32(SP), BP
        0x0080 00128 (main.go:4)        ADDQ    $40, SP
        0x0084 00132 (main.go:4)        RET

```

传参和返回都采用了寄存器进行传递，并且返回值和输入都使用了完全相同的寄存器序列，并且使用的顺序也是一致的。

因为这个优化，在一些函数调用嵌套层次较深的场景下，内存有一定概率会降低，有机会做压测可以试一试～。



## 总结

熟练掌握并理解函数的调用过程是我们深入学习`Go`语言的重要一课，看完本文希望你已经熟练掌握了函数的调用惯例～。

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)

