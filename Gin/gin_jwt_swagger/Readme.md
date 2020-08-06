# gin jwt swagger example
`author: asong`
`time: 2020-08-02`
`公众号：Golang梦工厂`
![qrCode](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)


|   版本  | 更新日期  |  更新内容 |
|  ----  | ----  | ---- |
| v1.0.0 | 2020.08.02 12:00| 增加版本说明、项目介绍、项目架构 |
| v1.0.1  | 2020.08.02 1:00| 增加mysql数据库设计 |

## 版本说明
- golang: 1.14
- mysql: 5.7
- redis: 4.0.14

## 项目介绍

这个项目主要是为了方便大家学习gin框架、jwt、swagger（后续会添加其他模块，在此基础上添加），因此将其结合，完成用户登陆、修改操作,使用jwt进行鉴权，
swagger生成接口文档。

## 使用说明

## 项目架构
### 系统架构
![项目架构](./static/images/system.png)


### 目录结构


## 数据库设计

### mysql

#### 用户表结构设计(users)

```mysql
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL,
  `username` varchar(64) NOT NULL,
  `nickname` varchar(255) DEFAULT NULL,
  `password` varchar(64) NOT NULL,
  `salt` varchar(64) NOT NULL,
  `avatar` varchar(128) NOT NULL,
  `uptime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

```
|字段|类型|KEY|可否为空|注释|
|----|----|----|----|----|
|id|bigint(20)|PRI|not||
|username|varchar(64)|UNI|not|用户名|
|nickname|varchar(255)| |yes|昵称|
|password|varchar(64)| |not|密钥|
|salt|varchar(16)| |not|属性|
|avatar|varchar(128)| |yes|头像地址|
|uptime|datetime| |yes|更新信息时间|