package main

import (
	_ "asong.cloud/Chatroom/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
	_ "github.com/go-sql-driver/mysql"
)

/**
	Author: Asong
	create by on 2020/6/1
	订阅号: Golang梦工厂
	Features: 初始化数据库连接
*/
func init()  {
	//注册mysql数据库驱动
	orm.RegisterDriver("mysql",orm.DRMySQL)
	//连接数据库
	orm.RegisterDataBase("default","mysql","root:root@tcp(127.0.0.1:3306)/chatroom?charset=utf8",30)
	//create table
	orm.RunSyncdb("default",false,true)
}

const App_Version  = "2020.5.30.1.0.0"
/**
	Author: Asong
	create by on 2020/6/1
	订阅号: Golang梦工厂
 */
func main() {
	beego.Info(beego.BConfig.AppName,App_Version)
	beego.AddFuncMap("i18n",i18n.Tr)
	beego.Run()
}

