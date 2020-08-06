package initserver

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
)
func InitRedis()  {
	client := redis.NewClient(&redis.Options{
		Addr: global.AsongServer.Redis.Addr,
		Password: global.AsongServer.Redis.Password,
		DB: global.AsongServer.Redis.Db,
	})
	pong , err := client.Ping().Result()
	if err != nil{
		global.AsongLogger.WithFields(logrus.Fields{"err":err}).Error("init redis err")
	}
	global.AsongLogger.WithFields(logrus.Fields{"err":err}).Info("redis connect ping response:",pong)
	global.AsongRedis = client
}
