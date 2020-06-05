# chatroom 

## 前言

`chatroom`是一个完整的多人聊天室项目，采用beego框架进行开发，使用mysql数据对用户信息与聊天信息进行存储，可以作为入门beego框架的项目demo。

## 项目介绍

`chatroom`是一个完整的多人聊天室项目，基于beego + mysql实现。功能具有用户注册、技术选择（默认选择websocket）、用户登陆、聊天信息发送、聊天信息保存。密码加密采用了MD5算法，使用Go官方自带库创建。该项目前后端没有分离，已经全部上传，可自行下载学习。

## 项目演示

目前该项目还没部署到服务器，不过dockerfile文件已经写好，可以自行学习部署到自己服务器，这里就不做项目演示了。

### 项目运行

将项目克隆到本地，在该目录下打开终端，输入并执行以下命令：

```go
go mod init ChatRoom
go get github.com/astaxie/beego
go get github.com/beego/bee
go get github.com/go-sql-driver/mysql
```
因为要使用bee工具运行该项目，下载好后需要到你的GOPATH/bin目录下查看有没有bee可执行文件，默认下载bee包会在bin目录下生成执行文件。如果没有可以关注我的另一篇博客自己重新编译即可（https://blog.csdn.net/qq_39397165/article/details/106406773）

安装bee工具成功后，我们需要创建一个数据库，数据库名可以自己定，使用如下命令进行创建：

```sql
create database chatroom
```
这里不需要创建数据库表，因为beego orm会自动根据结构体中的结构自动生成相关字段。

以上都成功后，可以进行项目运行。

```go
bee run
```
该项目即可成功运行，在浏览器打开localhost:8080即可。

### 组织架构

``` lua
chatroom
├── conf -- 配置文件
├── controllers -- 控制层相关代码
├── models -- 业务逻辑层相关代码
├── routers -- 路由控制
├── static -- 静态文件
├── utils -- 各个插件代码
└── views -- 前端代码
```



## 公众号

**关注公众号**，第一时间观看优质文章，第一时间获取资料。

![公众号图片](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)