package common

import (
	"github.com/gin-gonic/gin"
)

func NewRouterClient()  *gin.Engine{
	engine := gin.Default()
	return engine
}
