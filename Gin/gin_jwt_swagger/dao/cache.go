package dao

import (
	"time"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
)

func GetSting(token string) (error, string) {
	value, err := global.AsongRedis.Get(token).Result()
	return err, value
}

func SetString(token string, username string) error {
	err := global.AsongRedis.Set(token, username, time.Duration(global.AsongServer.Redis.Cache.Tokenexpired)).Err()
	return err
}
