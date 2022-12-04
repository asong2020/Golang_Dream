package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

/**
	Author: Asong
	create by on 2020/6/1
	订阅号: Golang梦工厂
	Features: 聊天记录
*/

type Message struct {
	Id int `auto`
	Tel string `orm:"size(128)"`
	Send time.Time `orm:"type(date)"`
	Content string `orm:"type(text)"`
}

func init()  {
	//注册模型 使用orm.QuerySeter进行高级查询
	orm.RegisterModel(new(Message))
}

//添加消息
func AddMessage(me *Message) error {
	o := orm.NewOrm()
	_,err:=o.Insert(me)
	return err
}