package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main()  {
	go func() {
		for{
				client := &http.Client{
					Timeout: 15 * time.Second,
					Transport: &http.Transport{
						MaxIdleConnsPerHost: 100,
						DisableKeepAlives:   false,
					},
				}
				resp,err := client.Get("http://localhost:8080/ping")
				if err != nil{
					return
				}
				if resp.StatusCode != 200{
					fmt.Println(resp.StatusCode)
				}
		}
	}()
	r := gin.Default()
	r.GET("/ping", func(context *gin.Context) {
		client := &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
				DisableKeepAlives:   false,
			},
		}
		_,err := client.Get("http://www.baidu.com")
		if err != nil{
			context.JSON(200,gin.H{
				"err": err,
			})
			return
		}

		context.JSON(200,gin.H{
			"message": "ok",
		})
	})
	r.Run()
}
