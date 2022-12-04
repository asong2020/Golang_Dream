## 前言

> 哈喽，大家好，我是`asong`。
>
> 我们都知道在`Unix`中万物都被称为文件，文件处理是一个非常常见的问题，所以本文就总结了`Go`语言操作文件的常见方式，整体思路如下：

**Go语言版本：1.18**

本文所有代码已经上传`github`：https://github.com/asong2020/Golang_Dream/tree/master/code_demo/file_operate_demo



## 操作文件包括哪些

操作一个文件离不开这几个动作：

- 创建文件
- 打开文件
- 读取文件
- 写入文件
- 关闭文件
- 打包/解包
- 压缩/解压缩
- 改变文件权限
- 删除文件
- 移动文件
- 重命名文件
- 清空文件

所以本文就针对这些操作总结了一些示例方法供大家参考；

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-07-10%20%E4%B8%8B%E5%8D%887.18.55.png)




## Go语言操作文件可使用的库

Go语言官方库：`os`、`io/ioutil`、`bufio`涵盖了文件操作的所有场景，

`os`提供了对文件`IO`直接调用的方法，`bufio`提供缓冲区操作文件的方法，`io/ioutil`也提供对文件`IO`直接调用的方法，不过Go语言在`Go1.16`版本已经弃用了`io/ioutil`库，这个[`io/ioutil`](https://go.dev/pkg/io/ioutil/)包是一个定义不明确且难以理解的东西集合。该软件包提供的所有功能都已移至其他软件包，所以`io/ioutil`中操作文件的方法都在`io`库有相同含义的方法，大家以后在使用到ioutil中的方法是可以通过注释在其他包找到对应的方法。



## 文件的基础操作

这里我把 创建文件、打开文件、关闭文件、改变文件权限这些归为对文件的基本操作，对文件的基本操作直接使用`os`库中的方法即可，因为我们需要进行`IO`操作，来看下面的例子：

```go
import (
	"log"
	"os"
)
func main() {
	// 创建文件
	f, err := os.Create("asong.txt")
	if err != nil{
		log.Fatalf("create file failed err=%s\n", err)
	}
	// 获取文件信息
	fileInfo, err := f.Stat()
	if err != nil{
		log.Fatalf("get file info failed err=%s\n", err)
	}

	log.Printf("File Name is %s\n", fileInfo.Name())
	log.Printf("File Permissions is %s\n", fileInfo.Mode())
	log.Printf("File ModTime is %s\n", fileInfo.ModTime())

	// 改变文件权限
	err = f.Chmod(0777)
	if err != nil{
		log.Fatalf("chmod file failed err=%s\n", err)
	}

	// 改变拥有者
	err = f.Chown(os.Getuid(), os.Getgid())
	if err != nil{
		log.Fatalf("chown file failed err=%s\n", err)
	}

	// 再次获取文件信息 验证改变是否正确
	fileInfo, err = f.Stat()
	if err != nil{
		log.Fatalf("get file info second failed err=%s\n", err)
	}
	log.Printf("File change Permissions is %s\n", fileInfo.Mode())

	// 关闭文件
	err = f.Close()
	if err != nil{
		log.Fatalf("close file failed err=%s\n", err)
	}
	
	// 删除文件
	err = os.Remove("asong.txt")
	if err != nil{
		log.Fatalf("remove file failed err=%s\n", err)
	}
}
```



## 写文件

### 快写文件

`os`/`ioutil`包都提供了`WriteFile`方法可以快速处理创建/打开文件/写数据/关闭文件，使用示例如下：

```go
func writeAll(filename string) error {
	err := os.WriteFile("asong.txt", []byte("Hi asong\n"), 0666)
	if err != nil {
		return err
	}
	return nil
}
```



### 按行写文件

`os`、`buffo`写数据都没有提供按行写入的方法，所以我们可以在调用`os.WriteString`、`bufio.WriteString`方法是在数据中加入换行符即可，来看示例：

```go
import (
	"bufio"
	"log"
	"os"
)
// 直接操作IO
func writeLine(filename string) error {
	data := []string{
		"asong",
		"test",
		"123",
	}
	f, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil{
		return err
	}

	for _, line := range data{
		_,err := f.WriteString(line + "\n")
		if err != nil{
			return err
		}
	}
	f.Close()
	return nil
}
// 使用缓存区写入
func writeLine2(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	// 为这个文件创建buffered writer
	bufferedWriter := bufio.NewWriter(file)
	
	for i:=0; i < 2; i++{
		// 写字符串到buffer
		bytesWritten, err := bufferedWriter.WriteString(
			"asong真帅\n",
		)
		if err != nil {
			return err
		}
		log.Printf("Bytes written: %d\n", bytesWritten)
	}
	// 写内存buffer到硬盘
	err = bufferedWriter.Flush()
	if err != nil{
		return err
	}

	file.Close()
	return nil
}
```



### 偏移量写入

某些场景我们想根据给定的偏移量写入数据，可以使用`os`中的`writeAt`方法，例子如下：

```go
import "os"

func writeAt(filename string) error {
	data := []byte{
		0x41, // A
		0x73, // s
		0x20, // space
		0x20, // space
		0x67, // g
	}
	f, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil{
		return err
	}
	_, err = f.Write(data)
	if err != nil{
		return err
	}

	replaceSplace := []byte{
		0x6F, // o
		0x6E, // n
	}
	_, err = f.WriteAt(replaceSplace, 2)
	if err != nil{
		return err
	}
	f.Close()
	return nil
}
```



### 缓存区写入

`os`库中的方法对文件都是直接的`IO`操作，频繁的`IO`操作会增加`CPU`的中断频率，所以我们可以使用内存缓存区来减少`IO`操作，在写字节到硬盘前使用内存缓存，当内存缓存区的容量到达一定数值时在写内存数据buffer到硬盘，`bufio`就是这样示一个库，来个例子我们看一下怎么使用：

```go
import (
	"bufio"
	"log"
	"os"
)

func writeBuffer(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	// 为这个文件创建buffered writer
	bufferedWriter := bufio.NewWriter(file)

	// 写字符串到buffer
	bytesWritten, err := bufferedWriter.WriteString(
		"asong真帅\n",
	)
	if err != nil {
		return err
	}
	log.Printf("Bytes written: %d\n", bytesWritten)

	// 检查缓存中的字节数
	unflushedBufferSize := bufferedWriter.Buffered()
	log.Printf("Bytes buffered: %d\n", unflushedBufferSize)

	// 还有多少字节可用（未使用的缓存大小）
	bytesAvailable := bufferedWriter.Available()
	if err != nil {
		return err
	}
	log.Printf("Available buffer: %d\n", bytesAvailable)
	// 写内存buffer到硬盘
	err = bufferedWriter.Flush()
	if err != nil{
		return err
	}

	file.Close()
	return nil
}
```



## 读文件

### 读取全文件

有两种方式我们可以读取全文件：

- `os`、`io/ioutil`中提供了`readFile`方法可以快速读取全文
- `io/ioutil`中提供了`ReadAll`方法在打开文件句柄后可以读取全文；

```go
import (
	"io/ioutil"
	"log"
	"os"
)

func readAll(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	log.Printf("read %s content is %s", filename, data)
	return nil
}

func ReadAll2(filename string) error {
	file, err := os.Open("asong.txt")
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(file)
	log.Printf("read %s content is %s\n", filename, content)

	file.Close()
	return nil
}
```



### 逐行读取

`os`库中提供了`Read`方法是按照字节长度读取，如果我们想要按行读取文件需要配合`bufio`一起使用，`bufio`中提供了三种方法`ReadLine`、`ReadBytes("\n")`、`ReadString("\n")`可以按行读取数据，下面我使用`ReadBytes("\n")`来写个例子：

```go
func readLine(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	bufferedReader := bufio.NewReader(file)
	for {
		// ReadLine is a low-level line-reading primitive. Most callers should use
		// ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
		lineBytes, err := bufferedReader.ReadBytes('\n')
		bufferedReader.ReadLine()
		line := strings.TrimSpace(string(lineBytes))
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		log.Printf("readline %s every line data is %s\n", filename, line)
	}
	file.Close()
	return nil
}
```



### 按块读取文件

有些场景我们想按照字节长度读取文件，这时我们可以如下方法：

- `os`库的`Read`方法
- `os`库配合`bufio.NewReader`调用`Read`方法
- `os`库配合`io`库的`ReadFull`、`ReadAtLeast`方法

```go
// use bufio.NewReader
func readByte(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	// 创建 Reader
	r := bufio.NewReader(file)

	// 每次读取 2 个字节
	buf := make([]byte, 2)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}
		log.Printf("writeByte %s every read 2 byte is %s\n", filename, string(buf[:n]))
	}
	file.Close()
	return nil
}

// use os
func readByte2(filename string) error{
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	// 每次读取 2 个字节
	buf := make([]byte, 2)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}
		log.Printf("writeByte %s every read 2 byte is %s\n", filename, string(buf[:n]))
	}
	file.Close()
	return nil
}


// use os and io.ReadAtLeast
func readByte3(filename string) error{
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	// 每次读取 2 个字节
	buf := make([]byte, 2)
	for {
		n, err := io.ReadAtLeast(file, buf, 0)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}
		log.Printf("writeByte %s every read 2 byte is %s\n", filename, string(buf[:n]))
	}
	file.Close()
	return nil
}
```



### 分隔符读取

`bufio`包中提供了`Scanner`扫描器模块，它的主要作用是把数据流分割成一个个标记并除去它们之间的空格，他支持我们定制`Split`函数做为分隔函数，分隔符可以不是一个简单的字节或者字符，我们可以自定义分隔函数，在分隔函数实现分隔规则以及指针移动多少，返回什么数据，如果没有定制`Split`函数，那么就会使用默认`ScanLines`作为分隔函数，也就是使用换行作为分隔符，`bufio`中还提供了默认方法`ScanRunes`、`ScanWrods`，下面我们用`SacnWrods`方法写个例子，获取用空格分隔的文本：

```go
func readScanner(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	// 可以定制Split函数做分隔函数
	// ScanWords 是scanner自带的分隔函数用来找空格分隔的文本字
	scanner.Split(bufio.ScanWords)
	for {
		success := scanner.Scan()
		if success == false {
			// 出现错误或者EOF是返回Error
			err = scanner.Err()
			if err == nil {
				log.Println("Scan completed and reached EOF")
				break
			} else {
				return err
			}
		}
		// 得到数据，Bytes() 或者 Text()
		log.Printf("readScanner get data is %s", scanner.Text())
	}
	file.Close()
	return nil
}
```



## 打包/解包

Go语言的`archive`包中提供了`tar`、`zip`两种打包/解包方法，这里以`zip`的打包/解包为例子：

`zip`解包示例：

```go
import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
)

func main()  {
	// Open a zip archive for reading.
	r, err := zip.OpenReader("asong.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		fmt.Printf("Contents of %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.CopyN(os.Stdout, rc, 68)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
	}
}
```

`zip`打包示例：

```go
func writerZip()  {
	// Create archive
	zipPath := "out.zip"
	zipFile, err := os.Create(zipPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new zip archive.
	w := zip.NewWriter(zipFile)
	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"asong.txt", "This archive contains some text files."},
		{"todo.txt", "Get animal handling licence.\nWrite more examples."},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatal(err)
		}
	}
	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}
```



## 总结

本文归根结底是介绍`os`、`io`、`bufio`这些包如何操作文件，因为`Go`语言操作提供了太多了方法，借着本文全都介绍出来，在使用的时候可以很方便的当作文档查询，如果你问用什么方法操作文件是最优的方法，这个我也没法回答你，需要根据具体场景分析的，如果这些方法你都知道了，在写一个benchmark对比一下就可以了，实践才是检验真理的唯一标准。

本文所有代码已经上传`github`：https://github.com/asong2020/Golang_Dream/tree/master/code_demo/file_operate_demo

好啦，本文到这里就结束了，我是**asong**，我们下期见。

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**


![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)
