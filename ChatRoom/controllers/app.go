package controllers

import (
	"asong.cloud/Chatroom/models"
	"asong.cloud/Chatroom/util"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"strings"
)

/**
Author: Asong
create by on 2020/6/1
订阅号: Golang梦工厂
Features: 首页控制层
*/

var langTypes []string
//初始化语言类型
func init()  {
	//从配置文件中获取语言支持类型  通过| 进行区分
	langTypes = strings.Split(beego.AppConfig.String("lang_types"),"|")
	//加载本地化文件
	for _,lang:=range langTypes{
		beego.Trace("Loading language" + lang)
		if err := i18n.SetMessage(lang,"conf/"+"locale_"+lang+".ini");err !=nil{
			beego.Error("Fail to set message file:",err)
			return
		}
	}
}
//其他路由器的基础路由
type baseController struct{
	beego.Controller
	i18n.Locale //用于处理数据和渲染模板时使用的i18n
}
//在 init之后执行 然后执行请求函数
// 用于语言选项的检查和设置
func (this *baseController)Prepare()  {
	this.Lang = ""
	//获取语言类型从 'Accept-Language'

	al := this.Ctx.Request.Header.Get("Accept-Language")
	beego.Info("Get Accept-Language: "+al)
	if len(al) > 4 {
		al = al[:5] // Only compare first 5 letters.
		if i18n.IsExist(al) {
			this.Lang = al
		}
	}

	// 没有则是English
	if len(al) == 0{
		this.Lang = "en-US"
	}
	//设置模板级语言选项
	this.Data["Lang"] = this.Lang
}

// AppController 处理 首页屏幕 允许用户选择技术和用户名
type AppController struct {
	baseController
}

func (this *AppController)Get()  {
	this.TplName = "welcome.html"
}

//join 方法处理 POST 请求
func (this *AppController)Join()  {
	tel := this.GetString("telephone")
	password := this.GetString("password")
	tech := this.GetString("tech")
	if len(tel)== 0 || len(password) == 0{
		this.Redirect("/",302)
		return
	}
	//数据库查询
	us,err := models.QueryByTel(tel)
	if err!=nil {
		beego.Info("不存在用户, 请进行注册")
		this.Redirect("/",302)
		return
	}else if us.Password != util.Get32MD5Encode(password) {
		beego.Info("密码错误,请输入正确的密码")
		this.Redirect("/",302)
		return
	}
	switch tech {
	case "websocket":
		beego.Info("WebSocket jump")
		this.Redirect("/ws?name="+us.Name+"&tel="+us.Tel,302)
	default:
		this.Redirect("/",302)

	}
	return
}