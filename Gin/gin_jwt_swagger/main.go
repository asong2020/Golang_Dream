package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fvbock/endless"
	"github.com/sirupsen/logrus"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/config"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/initserver"
)

func init()  {
	initserver.LoggerInit()
	dir,err := os.Getwd()
	if err != nil{
		global.AsongLogger.WithFields(logrus.Fields{"err":err}).Error("get dir error")
	}
	err = config.CofParse(dir+"/config.yaml",&global.AsongServer)
	if err != nil{
		global.AsongLogger.WithFields(logrus.Fields{"err": err}).Error("init config file error")
	}
	global.AsongLogger.WithFields(logrus.Fields{"err": err}).Info("init config file success")
	initserver.InitMysql()
	global.AsongLogger.WithFields(logrus.Fields{"err":err}).Info("init mysql success")
	initserver.InitRedis()

}

func main()  {
	//启动项目
	run()
}

func run()  {
	address := fmt.Sprintf(":%d", global.AsongServer.System.Port)
	router := initserver.RoutersInit()
	server := endless.NewServer(address,router)
	server.ReadHeaderTimeout = 10 * time.Millisecond
	server.WriteTimeout = 10 * time.Second
	server.MaxHeaderBytes = 1 << 20
	global.AsongLogger.Info("欢迎关注公众号:Golang梦工厂")
	global.AsongLogger.Info("gin_jwt_swagger 项目已启动")
	global.AsongLogger.Info("swagger文档地址：http://localhost:8888/swagger/index.html")
	err := server.ListenAndServe()
	global.AsongLogger.WithFields(logrus.Fields{"err":err}).Info("gin start error")
}