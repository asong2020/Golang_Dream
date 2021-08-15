package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r0 := r.Group("/api")
	{
		r.GET("/pong", func(context *gin.Context) {
			context.JSON(200,gin.H{
				"message": "pong",
			})
		})
		r0.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}