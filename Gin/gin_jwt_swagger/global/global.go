package global

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/config"
)

var (
	AsongServer *config.Server
	AsongLogger *logrus.Logger
	AsongDb *sql.DB
	AsongRedis *redis.Client
)