package dao

import (
	sql2 "database/sql"

	"github.com/sirupsen/logrus"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/model"
)

func InsertUser(user *model.User) error {
	sql := `insert into users(username,nickname,password,salt,avatar,uptime) values(?,?,?,?,?,?)`
	result, err := global.AsongDb.Exec(sql, user.Username, user.Nickname, user.Password, user.Salt, user.Avatar, user.Uptime)
	if err != nil {
		global.AsongLogger.WithFields(logrus.Fields{"err": err}).Error("insert user err and the username is ", user.Username)
	}
	_, err = result.LastInsertId()
	if err != nil {
		global.AsongLogger.WithFields(logrus.Fields{"err": err}).Error("get insert id err")
	}
	return err
}

func QueryByUsername(username string) bool {
	sql := `select username from users where username = ?`
	var user string
	err := global.AsongDb.QueryRow(sql, username).Scan(&user)
	if err == sql2.ErrNoRows {
		global.AsongLogger.WithFields(logrus.Fields{"err": err}).Error("query user by username")
		return true
	}
	return false
}

func GetByUsername(user *model.User) (model.User, error) {
	var u model.User
	sql := `select username,password,salt,nickname,avatar from users where username = ?`
	err := global.AsongDb.QueryRow(sql, user.Username).Scan(&u.Username, &u.Password, &u.Salt, &u.Nickname, &u.Avatar)
	if err == sql2.ErrNoRows {
		global.AsongLogger.Error("please register")
		return u, err
	}
	return u, err
}

func UpdatePassword(password string) error {
	sql := `update users set password = ?`
	_, err := global.AsongDb.Exec(sql, password)
	if err != nil {
		global.AsongLogger.WithFields(logrus.Fields{"err": err}).Error("update password error")
		return err
	}
	return nil
}
