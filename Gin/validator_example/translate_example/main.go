package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

// loca 通常取决于 http 请求头的 'Accept-Language'
func transInit(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() //chinese
		enT := en.New() //english
		uni := ut.New(enT, zhT, enT)

		var o bool
		trans, o = uni.GetTranslator(local)
		if !o {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}
		//register translate
		// 注册翻译器
		switch local {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,max=16,min=6"`
}

func main() {
	if err := transInit("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}
	router := gin.Default()

	router.POST("/user/login", login)

	err := router.Run(":8888")
	if err != nil {
		log.Println("failed")
	}
}

func login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		c.JSON(http.StatusOK, gin.H{
			"msg": errs.Translate(trans),
		})
		return
	}
	//login 操作省略
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}
