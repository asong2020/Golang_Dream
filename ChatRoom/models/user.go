package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

/**
	Author: Asong
	create by on 2020/6/1
	订阅号: Golang梦工厂
	Features: User数据表
*/

type User struct {
	Id int `auto`
	Name string `orm:"size(100)"`
	Tel string `orm:"size(128)"`
	Password string `orm:"size(256)"`
	Email string `orm:"size(256)"`
}

func init()  {
	//注册模型 使用orm.QuerySeter进行高级查询
	orm.RegisterModel(new(User))
}


//添加用户表
func AddUser(us *User)  (int64,error){
	o := orm.NewOrm()
	id,err := o.Insert(us)
	if err != nil{
		return id,err
	}
	beego.Info("INSERT TO USER ID:", id)
	return id,err
}
//更新用户表
func UpdateUser(us *User)  error{
	o := orm.NewOrm()
	num,err:= o.Update(us)
	if err != nil{
		return err
	}
	beego.Info("UPDATE NUM:",num)
	return err
}
//读取用户表
func ReadUser(us *User)  (*User,error){
	o := orm.NewOrm()
	err := o.Read(us)
	return us,err
}
//删除用户表
func DeleteUser(us *User) error {
	o := orm.NewOrm()
	num, err := o.Delete(us)
	beego.Info("DELETE NUM:",num)
	return err
}
//根据电话号查询信息
func QueryByTel(tel string)  (*User,error){
	us := new(User)
	o := orm.NewOrm()
	err:=o.Raw("SELECT * FROM user WHERE tel = ?",tel).QueryRow(us)
	return us,err
}