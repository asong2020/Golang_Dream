package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/afex/hystrix-go/hystrix"
)
var filename = "./test.txt"
var f *os.File
var errfile error
var CircuitBreakerName = "api_%s_circuit_breaker"

func CircuitBreakerWrapper(ctx *gin.Context){
	name := fmt.Sprintf(CircuitBreakerName,ctx.Request.URL)
	hystrix.Do(name, func() error {
		ctx.Next()
		code := ctx.Writer.Status()
		if code != http.StatusOK{
			return errors.New(fmt.Sprintf("status code %d", code))
		}
		return nil

	}, func(err error) error {
		if err != nil{
			// 监控上报（未实现）
			_, _ = io.WriteString(f, fmt.Sprintf("circuitBreaker and err is %s\n",err.Error())) //写入文件(字符串)
			fmt.Printf("circuitBreaker and err is %s\n",err.Error())
			// 返回熔断错误
			ctx.JSON(http.StatusServiceUnavailable,gin.H{
				"msg": err.Error(),
			})
		}
		return nil
	})
}

func init()  {
	hystrix.ConfigureCommand(CircuitBreakerName,hystrix.CommandConfig{
		Timeout:                int(3*time.Second), // 执行command的超时时间为3s
		MaxConcurrentRequests:  10, // command的最大并发量
		RequestVolumeThreshold: 100, // 统计窗口10s内的请求数量，达到这个请求数量后才去判断是否要开启熔断
		SleepWindow:            int(2 * time.Second), // 当熔断器被打开后，SleepWindow的时间就是控制过多久后去尝试服务是否可用了
		ErrorPercentThreshold:  20, // 错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
	})
	if checkFileIsExist(filename) { //如果文件存在
		f, errfile = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
	} else {
		f, errfile = os.Create(filename) //创建文件
	}
}


func main()  {
	defer f.Close()
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
	r := gin.Default()
	r.GET("/api/ping/baidu", func(c *gin.Context) {
		_, err := http.Get("https://www.baidu.com")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "success"})
	}, CircuitBreakerWrapper)
	r.Run()  // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}



