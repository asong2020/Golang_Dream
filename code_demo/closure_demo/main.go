package main

import (
	"github.com/gin-gonic/gin"
)

func VerifyHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("token")
		if header == "" {
			c.JSON(200, gin.H{
				"code":   1000,
				"msg":    "Not logged in",
			})
			return
		}
	}
}


func main()  {
	r := gin.Default()
	group := r.Group("/api/asong",VerifyHeader())
	{
		group.GET("/ping", func(context *gin.Context) {
			context.JSON(200,gin.H{
				"message": "pong",
			})
		})
	}
	r.Run()
}
