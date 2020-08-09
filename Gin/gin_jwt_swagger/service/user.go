package service

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/dao"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/middleware"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/model"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/model/request"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/model/resp"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/util"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/util/response"
)

func Register(u *model.User)  error{
	register := dao.QueryByUsername(u.Username)
	if !register {
		global.AsongLogger.Error("already exit the user")
		return errors.New("已经注册")
	}
	u.Salt = util.GenerateSalt()
	u.Password = util.MD5V([]byte(u.Password)) + u.Salt
	u.Uptime = time.Now()
	err := dao.InsertUser(u)
	if err != nil{
		global.AsongLogger.WithFields(logrus.Fields{"err":err}).Error("insert user error")
		return err
	}
	return nil
}

func Login(u *model.User)  (*model.User,error){
	user,err :=dao.GetByUsername(u)
	if err != nil{
		return &user,err
	}
	u.Password = util.MD5V([]byte(u.Password)) + user.Salt
	if u.Password == user.Password{
		return &user,nil
	}
	return &user,errors.New("密码不对")
}

func GenerateTokenForUser(c *gin.Context,u *model.User)  {
	j := &middleware.JWT{
		SigningKey: []byte(global.AsongServer.Jwt.Signkey), // 唯一签名
	}
	claims := request.UserClaims{
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,       // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*2, // 过期时间 2h
			Issuer:    "asong",                       // 签名的发行者
		},
	}
	token ,err := j.GenerateToken(claims)
	if err != nil {
		response.FailWithMessage("获取token失败", c)
		return
	}
	res := resp.ResponseUser{Username: u.Username,Nickname: u.Nickname,Avatar: u.Avatar}
	response.OkWithData(resp.LoginResponse{User: res,Token: token,ExpiresAt: claims.ExpiresAt *1000},c)
	return
}

func ChangePassword(u *model.User,newPassword string) error {
	user,err :=dao.GetByUsername(u)
	if err != nil{
		return err
	}
	u.Password = util.MD5V([]byte(u.Password)) + user.Salt
	if u.Password != user.Password{
		return errors.New("密码不对")
	}
	newPassword = util.MD5V([]byte(newPassword)) + user.Salt
	err = dao.UpdatePassword(newPassword)
	if err != nil{
		return err
	}
	return nil
}