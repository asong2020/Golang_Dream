 ## 前言

哈喽，大家好，我是`asong`；

通常我们在业务项目中会借助使用静态代码检查工具来保证代码质量，通过静态代码检查工具我们可以提前发现一些问题，比如变量未定义、类型不匹配、变量作用域问题、数组下标越界、内存泄露等问题，工具会按照自己的规则进行问题的严重等级划分，给出不同的标识和提示，静态代码检查助我们尽早的发现问题，`Go`语言中常用的静态代码检查工具有`golang-lint`、`golint`，这些工具中已经制定好了一些规则，虽然已经可以满足大多数场景，但是有些时候我们会遇到针对特殊场景来做一些定制化规则的需求，所以本文我们一起来学习一下如何自定义linter需求；



## Go语言中的静态检查是如何实现？

众所周知`Go`语言是一门编译型语言，编译型语言离不开词法分析、语法分析、语义分析、优化、编译链接几个阶段，学过编译原理的朋友对下面这个图应该很熟悉：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-05-29%20%E4%B8%8B%E5%8D%881.46.27.png)

编译器将高级语言翻译成机器语言，会先对源代码做词法分析，词法分析是将字符序列转换为**Token**序列的过程，**Token**一般分为这几类：关键字、标识符、字面量（包含数字、字符串）、特殊符号（如加号、等号），生成`Token`序列后，需要进行语法分析，进一步处理后，生成一棵以 表达式为结点的 语法树，这个语法树就是我们常说的`AST`，在生成语法树的过程就可以检测一些形式上的错误，比如括号缺少，语法分析完成后，就需要进行语义分析，在这里检查编译期所有能检查静态语义，后面的过程就是中间代码生成、目标代码生成与优化、链接，这里就不详细描述了，这里主要是想引出抽象语法树（AST），**我们的静态代码检查工具就是通过分析抽象语法树（AST）根据定制的规则来做的**；那么抽象语法树长什么样子呢？我们可以使用标准库提供的`go/ast`、`go/parser`、`go/token`包来打印出`AST`，也可以使用可视化工具：http://goast.yuroyoro.net/ 查看`AST`，具体`AST`长什么样我们可以看下文的例子；



## 制定linter规则

假设我们现在要在我们团队制定这样一个代码规范，所有函数的第一个参数类型必须是`Context`，不符合该规范的我们要给出警告；好了，现在规则已经订好了，现在我们就来想办法实现它；先来一个有问题的示例：

```go
// example.go
package main

func add(a, b int) int {
	return a + b
}
```

对应`AST`如下：

```go
*ast.FuncDecl {
     8  .  .  .  Name: *ast.Ident {
     9  .  .  .  .  NamePos: 3:6
    10  .  .  .  .  Name: "add" 
    11  .  .  .  .  Obj: *ast.Object {
    12  .  .  .  .  .  Kind: func
    13  .  .  .  .  .  Name: "add" // 函数名
    14  .  .  .  .  .  Decl: *(obj @ 7)
    15  .  .  .  .  }
    16  .  .  .  }
    17  .  .  .  Type: *ast.FuncType {
    18  .  .  .  .  Func: 3:1
    19  .  .  .  .  Params: *ast.FieldList {
    20  .  .  .  .  .  Opening: 3:9
    21  .  .  .  .  .  List: []*ast.Field (len = 1) {
    22  .  .  .  .  .  .  0: *ast.Field {
    23  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 2) {
    24  .  .  .  .  .  .  .  .  0: *ast.Ident {
    25  .  .  .  .  .  .  .  .  .  NamePos: 3:10
    26  .  .  .  .  .  .  .  .  .  Name: "a"
    27  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    28  .  .  .  .  .  .  .  .  .  .  Kind: var
    29  .  .  .  .  .  .  .  .  .  .  Name: "a"
    30  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 22)
    31  .  .  .  .  .  .  .  .  .  }
    32  .  .  .  .  .  .  .  .  }
    33  .  .  .  .  .  .  .  .  1: *ast.Ident {
    34  .  .  .  .  .  .  .  .  .  NamePos: 3:13
    35  .  .  .  .  .  .  .  .  .  Name: "b"
    36  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    37  .  .  .  .  .  .  .  .  .  .  Kind: var
    38  .  .  .  .  .  .  .  .  .  .  Name: "b"
    39  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 22)
    40  .  .  .  .  .  .  .  .  .  }
    41  .  .  .  .  .  .  .  .  }
    42  .  .  .  .  .  .  .  }
    43  .  .  .  .  .  .  .  Type: *ast.Ident {
    44  .  .  .  .  .  .  .  .  NamePos: 3:15
    45  .  .  .  .  .  .  .  .  Name: "int" // 参数名
    46  .  .  .  .  .  .  .  }
    47  .  .  .  .  .  .  }
    48  .  .  .  .  .  }
    49  .  .  .  .  .  Closing: 3:18
    50  .  .  .  .  }
    51  .  .  .  .  Results: *ast.FieldList {
    52  .  .  .  .  .  Opening: -
    53  .  .  .  .  .  List: []*ast.Field (len = 1) {
    54  .  .  .  .  .  .  0: *ast.Field {
    55  .  .  .  .  .  .  .  Type: *ast.Ident {
    56  .  .  .  .  .  .  .  .  NamePos: 3:20
    57  .  .  .  .  .  .  .  .  Name: "int"
    58  .  .  .  .  .  .  .  }
    59  .  .  .  .  .  .  }
    60  .  .  .  .  .  }
    61  .  .  .  .  .  Closing: -
    62  .  .  .  .  }
    63  .  .  .  }
```





## 方式一：标准库实现custom linter

通过上面的`AST`结构我们可以找到函数参数类型具体在哪个结构上，因为我们可以根据这个结构写出解析代码如下：

```go
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	v := visitor{fset: token.NewFileSet()}
	for _, filePath := range os.Args[1:] {
		if filePath == "--" { // to be able to run this like "go run main.go -- input.go"
			continue
		}

		f, err := parser.ParseFile(v.fset, filePath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filePath, err)
		}
		ast.Walk(&v, f)
	}
}

type visitor struct {
	fset *token.FileSet
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}

	params := funcDecl.Type.Params.List // get params
	// list is equal of zero that don't need to checker.
	if len(params) == 0 {
		return v
	}

	firstParamType, ok := params[0].Type.(*ast.SelectorExpr)
	if ok && firstParamType.Sel.Name == "Context" {
		return v
	}

	fmt.Printf("%s: %s function first params should be Context\n",
		v.fset.Position(node.Pos()), funcDecl.Name.Name)
	return v
}
```

然后执行命令如下：

```shell
$ go run ./main.go -- ./example.go
./example.go:3:1: add function first params should be Context
```

通过输出我们可以看到，函数`add()`第一个参数必须是**Context**；这就是一个简单实现，因为`AST`的结构实在是有点复杂，就不再这里详细介绍每个结构体了，可以看曹大之前写的一篇文章：[golang 和 ast](https://xargin.com/ast/)



## 方式二：go/analysis 

看过上面代码的朋友肯定有点抓狂了，有很多实体存在，要开发一个`linter`，我们需要搞懂好多实体，好在`go/analysis`进行了封装，`go/analysis`为`linter` 提供了统一的接口，它简化了与IDE,metalinters，代码Review等工具的集成。如，任何`go/analysis`linter都可以高效的被`go vet`执行，下面我们通过代码方式了来介绍`go/analysis`的优势；

新建一个项目代码结构如下：

```lua
.
├── firstparamcontext
│   └── firstparamcontext.go
├── go.mod
├── go.sum
└── testfirstparamcontext
    ├── example.go
    └── main.go

```

添加检查模块代码，在`firstparamcontext.go`添加如下代码：

```go
package firstparamcontext

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "firstparamcontext",
	Doc:  "Checks that functions first param type is Context",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}

		params := funcDecl.Type.Params.List // get params
		// list is equal of zero that don't need to checker.
		if len(params) == 0 {
			return true
		}

		firstParamType, ok := params[0].Type.(*ast.SelectorExpr)
		if ok && firstParamType.Sel.Name == "Context" {
			return true
		}

		pass.Reportf(node.Pos(), "''%s' function first params should be Context\n",
			funcDecl.Name.Name)
		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}
```

然后添加分析器：

```go
package main

import (
	"asong.cloud/Golang_Dream/code_demo/custom_linter/firstparamcontext"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(firstparamcontext.Analyzer)
}
```

命令行执行如下：

```shell
$ go run ./main.go -- ./example.go 
/Users/go/src/asong.cloud/Golang_Dream/code_demo/custom_linter/testfirstparamcontext/example.go:3:1: ''add' function first params should be Context
```

如果我们想添加更多的规则，使用`golang.org/x/tools/go/analysis/multichecker`追加即可。



## 集成到golang-cli

我们可以把`golang-cli`的代码下载到本地，然后在`pkg/golinters `下添加`firstparamcontext.go`，代码如下：

```go
import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/fisrtparamcontext"
)


func NewfirstparamcontextCheck() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"firstparamcontext",
		"Checks that functions first param type is Context",
		[]*analysis.Analyzer{firstparamcontext.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
```

然后重新`make`一个`golang-cli`可执行文件，加到我们的项目中就可以了；



## 总结

`golang-cli`仓库中`pkg/golinters`目录下存放了很多静态检查代码，学会一个知识点的最快办法就是抄代码，先学会怎么使用的，慢慢在把它变成我们自己的；本文没有对`AST`标准库做过多的介绍，因为这部分文字描述比较难以理解，最好的办法还是自己去看官方文档、加上实践才能更快的理解。

本文所有代码已经上传：https://github.com/asong2020/Golang_Dream/tree/master/code_demo/custom_linter

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)