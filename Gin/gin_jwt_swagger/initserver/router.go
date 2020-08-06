package initserver

import (
	"github.com/gin-gonic/gin"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/handler"
)

func RoutersInit() *gin.Engine {
	var Router = gin.Default()
	RouterGroup := Router.Group("")

	handler.RouterBaseInit(RouterGroup)
	handler.RouteUserInit(RouterGroup)

	global.AsongLogger.Info("router register success")
	return Router
}
