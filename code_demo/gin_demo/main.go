package main

import (
	"flag"
	"github.com/gin-gonic/gin"
)

var Port string

func init()  {
	flag.StringVar(&Port, "port", "8080", "Input Your Port")
}

func main() {
	flag.Parse()
	r := gin.Default()
	r.Use()
	r1 := r.Group("/api")
	{
		r1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	r.Run("localhost:" + Port)
}