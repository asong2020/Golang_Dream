package initserver

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/docs"

	"github.com/swaggo/gin-swagger/swaggerFiles"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/handler"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/middleware"
)

func RoutersInit() *gin.Engine {
	var Router = gin.Default()
	// 跨域
	Router.Use(middleware.Cors())
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	RouterGroup := Router.Group("")

	handler.RouterBaseInit(RouterGroup)
	handler.RouteUserInit(RouterGroup)

	global.AsongLogger.Info("router register success")
	return Router
}
