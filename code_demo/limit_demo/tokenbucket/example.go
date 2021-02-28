package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var rateLimit *rate.Limiter

func tokenBucket() gin.HandlerFunc {
	return func(c *gin.Context) {
		if rateLimit.Allow() {
			c.String(http.StatusOK, "rate limit,Drop")
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	limit := rate.Every(100 * time.Millisecond)
	rateLimit = rate.NewLimiter(limit, 10)
	r := gin.Default()
	r.GET("/ping", tokenBucket(), func(c *gin.Context) {
		c.JSON(200, true)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
