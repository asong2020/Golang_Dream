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

之前我们写了一个`Leaf` - 号段获取ID，我们借助这个服务我们来





