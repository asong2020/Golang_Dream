package handler

import (
	"github.com/gin-gonic/gin"
)

func RouteUserInit(Router *gin.RouterGroup)  {
	UserRouter := Router.Group("user")
	{
		UserRouter.POST("setPassword",setPassword)
	}
}

func RouterBaseInit(Router *gin.RouterGroup)  {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.POST("login",login)
	}
}

func login(c *gin.Context)  {

}

func setPassword(c *gin.Context)  {
	
}