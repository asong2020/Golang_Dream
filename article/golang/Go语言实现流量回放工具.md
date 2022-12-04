## 前言

> 哈喽，大家好，我是`asong`。
>
> 今天给大家推荐一款使用Go语言编写的流量回放工具 -- **goreplay**；工作中你一定遇到过需要在服务器上抓包的场景，有了这个工具就可以助你一臂之力，**goreplay**的功能十分强大，支持流量的放大、缩小，并且集成了`ElasticSearch`，将流量存入ES进行实时分析；
>
> 废话不多，我们接下来来看一看这个工具；



## goreplay介绍与安装

项目地址：https://github.com/buger/goreplay

**goreplay**是一个开源网络监控工具，可以实时记录TCP/HTTP流量，支持把流量记录到文件或者`elasticSearch`实时分析，也支持流量的放大、缩小，还支持频率限制；**goreplay**不是代理，无需任何代码入侵，只需要在服务相同的机器上运行`goreplay`守护程序，其会在后台侦听网络接口上的流量，`goreplay`的设计遵循 Unix 设计哲学：**一切都是由管道组成的，各种输入将数据复用为输出；**可以看一下官网画的架构图：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-09-04%20%E4%B8%8B%E5%8D%888.16.00.png)

`goreplay`的安装也比较简单，只需要在https://github.com/buger/goreplay/releases 下载对应操作系统的二进制文件即可，我的电脑是`mac`的：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-09-04%20%E4%B8%8B%E5%8D%888.19.53.png)

解压缩后就是一个二进制文件`gor`，将其添加到您的环境变量中，方便我们后续的操作；



## 使用示例

### 实时流量转发

首先我们要准备一个`Web`服务，最简单的就是用`Gin `快速实现一个`helloworld`，替大家实现好了：https://github.com/asong2020/Golang_Dream/tree/master/code_demo/gin_demo；

```go
import (
	"flag"
	"github.com/gin-gonic/gin"
)

var Port string

func init()  {
	flag.StringVar(&Port, "port", "8081", "Input Your Port")
}

func main() {
	flag.Parse()
	r := gin.Default()
	r.Use()
	r1 := r.Group("/api")
	{
		r1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	r.Run("localhost:" + Port)
}
```

因为资源有限，这里我用一台电脑起两个进程来模拟流量转发，分别启动两个web服务分别监控端口号`8081`、`8082`：

```bash
$ go run . --port="8081"
$ go run . --port="8082"
```

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-09-04%20%E4%B8%8B%E5%8D%888.48.07.png)

服务弄好了，现在我们来开启`gor`守护进程进行流量监听与转发，将`8081`端口的流量转发到`8082`端口上：

```bash
$ sudo gor --input-raw :8081 --output-http="http://127.0.0.1:8082"
```

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-09-04%20%E4%B8%8B%E5%8D%888.52.29.png)

现在我们请求`8081`端口：

```bash
$ curl --location --request GET 'http://127.0.0.1:8081/api/ping'
```

可以看到`8082`端口同样被请求了：

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-09-04%20%E4%B8%8B%E5%8D%888.53.51.png)



### 流量放大、缩小

`goreplay`支持将捕获的流量存储到文件中，实际工作中我们可以使用捕获的流量做压力测试，首先我们需要将捕获的流量保存到本地文件，然后利用该文件进行流量回放；

还是上面的`Web`程序，我们将端口`8081`的流量保存到本地文件：

```go
$ sudo gor --input-raw :8081 --output-file ./requests.gor
```

我们对`8081`端口执行了5次请求：

<img src="https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%88%AA%E5%B1%8F2022-09-04%20%E4%B8%8B%E5%8D%889.16.45.png" style="zoom:67%;" />

然后我们对`8082`端口进行流量缩小测试，缩小一倍：

```go
gor --input-file "requests_0.gor" --output-http="http://127.0.0.1:8082|50%"
```

调整百分比就是进行流量放大、缩小，这里我们缩小了一倍，可以看到只有2次请求到了`8082`端口；我们可以调整流量回放的速度，比如我们调整流量以10倍速度进行重播：

```bash
$ gor --input-file "requests_0.gor|1000%" --output-http="http://127.0.0.1:8082|50%" # 1000%就是放大10倍
```

### 流量写入到`ElastichSearch`

`goreplay`可以将捕获的流量导出到`Es`中，只需要执行如下命令：

```bash
$ gor --input-raw :8000 --output-http http://staging.cm  --output-http-elasticsearch localhost:9200/gor
```

我们不需要提前创建索引结构，他将自动创建，具体结构如下：

```go
type ESRequestResponse struct {
	ReqURL               string `json:"Req_URL"`
	ReqMethod            string `json:"Req_Method"`
	ReqUserAgent         string `json:"Req_User-Agent"`
	ReqAcceptLanguage    string `json:"Req_Accept-Language,omitempty"`
	ReqAccept            string `json:"Req_Accept,omitempty"`
	ReqAcceptEncoding    string `json:"Req_Accept-Encoding,omitempty"`
	ReqIfModifiedSince   string `json:"Req_If-Modified-Since,omitempty"`
	ReqConnection        string `json:"Req_Connection,omitempty"`
	ReqCookies           string `json:"Req_Cookies,omitempty"`
	RespStatus           string `json:"Resp_Status"`
	RespStatusCode       string `json:"Resp_Status-Code"`
	RespProto            string `json:"Resp_Proto,omitempty"`
	RespContentLength    string `json:"Resp_Content-Length,omitempty"`
	RespContentType      string `json:"Resp_Content-Type,omitempty"`
	RespTransferEncoding string `json:"Resp_Transfer-Encoding,omitempty"`
	RespContentEncoding  string `json:"Resp_Content-Encoding,omitempty"`
	RespExpires          string `json:"Resp_Expires,omitempty"`
	RespCacheControl     string `json:"Resp_Cache-Control,omitempty"`
	RespVary             string `json:"Resp_Vary,omitempty"`
	RespSetCookie        string `json:"Resp_Set-Cookie,omitempty"`
	Rtt                  int64  `json:"RTT"`
	Timestamp            time.Time
}
```

`goreplay`提供了太多的功能，就不一一介绍了，可以通过执行`help`命令查看其他高级用法，每个命令都提供了例子，入手很快；

```bash
$ gor -h
Gor is a simple http traffic replication tool written in Go. Its main goal is to replay traffic from production servers to staging and dev environments.
Project page: https://github.com/buger/gor
Author: <Leonid Bugaev> leonsbox@gmail.com
Current Version: v1.3.0

  -copy-buffer-size value
    	Set the buffer size for an individual request (default 5MB)
  -cpuprofile string
    	write cpu profile to file
  -exit-after duration
    	exit after specified duration
  -http-allow-header value
    	A regexp to match a specific header against. Requests with non-matching headers will be dropped:
    		 gor --input-raw :8080 --output-http staging.com --http-allow-header api-version:^v1
  -http-allow-method value
    	Whitelist of HTTP methods to replay. Anything else will be dropped:
    		gor --input-raw :8080 --output-http staging.com --http-allow-method GET --http-allow-method OPTIONS
  -http-allow-url value
    	A regexp to match requests against. Filter get matched against full url with domain. Anything else will be dropped:
    		 gor --input-raw :8080 --output-http staging.com --http-allow-url ^www.
  -http-basic-auth-filter value
    	A regexp to match the decoded basic auth string against. Requests with non-matching headers will be dropped:
    		 gor --input-raw :8080 --output-http staging.com --http-basic-auth-filter "^customer[0-9].*"
  -http-disallow-header value
    	A regexp to match a specific header against. Requests with matching headers will be dropped:
    		 gor --input-raw :8080 --output-http staging.com --http-disallow-header "User-Agent: Replayed by Gor"
    		 ..........省略
```



## `goreplay`基本实现原理

`goreplay`底层也是调用`Libpcap`，`Libpcap`即数据包捕获函数库，`tcpdump`也是基于这个库实现的，`Libpcap`是`C`语言写的，`Go`语言不能直接调用`C`语言，需要使用`CGo`，所以`goreplay`可以直接使用谷歌的包github.com/google/gopacket，提供了更方便的操作接口，基于`goreplay`封装了`input`、`output`，在启动的时候通过命令行参数解析指定的`input`、`output`，`input`读取数据写入到`output`中，默认是一个`input`复制多份，写多个`output`，多个`input`之前是并行的，但是单个`intput`到多个`output`是串行的，所以`input-file`会有性能瓶颈，压测的时候需要开多个进程同时跑来达到压测需求；

`goreplay`的源码有点多，就不在这里分析了，大家感兴趣哪一部分可以从`gor.go`的`main`函数入手，看自己感兴趣的部分就可以了；



## 总结

`goreplay`提供的玩法非常丰富，合理的改造可以做成回归工具帮助我们确保服务的稳定性，别放过这个自我展现的机会～。

**好啦，本文到这里就结束了，我是asong，我们下期见。**

**创建了读者交流群，欢迎各位大佬们踊跃入群，一起学习交流。入群方式：关注公众号获取。更多学习资料请到公众号领取。**


![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/扫码_搜索联合传播样式-白色版.png)
