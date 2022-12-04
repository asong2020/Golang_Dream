package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Info struct {
	CreateTime time.Time `form:"create_time" binding:"required,timing" time_format:"2006-01-02"`
	UpdateTime time.Time `form:"update_time" validate:"required,timing" time_format:"2006-01-02"`
}

// 自定义验证规则断言
func timing(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

func main() {
	route := gin.Default()
	// 注册验证
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("timing", timing)
		if err != nil {
			fmt.Println("success")
		}
	}

	route.GET("/time", getTime)
	route.Run(":8080")
}

func getTime(c *gin.Context) {
	var b Info
	// 数据模型绑定查询字符串验证
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "time are valid!"})
	} else {
		_, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "param is error"})
	}
}
