package controllers

import (
	"asong.cloud/Chatroom/models"
	"asong.cloud/Chatroom/util"
	"github.com/astaxie/beego"
)

/**
Author: Asong
create by on 2020/6/1
订阅号: Golang梦工厂
Features: 用户注册
*/

type RegisterControllers struct {
	baseController
}

func (this *RegisterControllers)Get()  {
	this.TplName = "register.html"
}

func (this *RegisterControllers)Join()  {
	user := new(models.User)
	user.Name = this.GetString("name")
	user.Tel = this.GetString("tel")
	user.Email = this.GetString("email")
	user.Password = util.Get32MD5Encode(this.GetString("password"))
	beego.Info(user.Name,user.Email,user.Tel)
	if len(user.Name) == 0 || len(user.Tel) == 0 || len(user.Password) == 0 || len(user.Email) == 0{
		this.Redirect("/register/join",302)
		return
	}
	//添加数据库
	id ,err :=models.AddUser(user)
	if err != nil{
		this.Redirect("/register/join",302)
		return
	}
	beego.Info(id)
	this.Redirect("/",302)
	return
}