package handler

import (
	"github.com/gin-gonic/gin"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/middleware"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/model"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/model/request"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/service"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/util/response"
)

func RouteUserInit(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user").Use(middleware.Auth())
	{
		UserRouter.PUT("setPassword", setPassword)
	}
}

func RouterBaseInit(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.POST("register", register)
		BaseRouter.POST("login", login)
	}
}

// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body request.LoginRequest true "用户登录接口"
// @Success 200 {string} string "{"success":true,"data": { "user": { "username": "asong", "nickname": "", "avatar": "" }, "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6ImFzb25nIiwiZXhwIjoxNTk2OTAyMzEyLCJpc3MiOiJhc29uZyIsIm5iZiI6MTU5Njg5NDExMn0.uUS1TreZusX-hL3nKOSNYZIeZ_0BGrxWjKI6xdpdO40", "expiresAt": 1596902312000 },,"msg":"操作成功"}"
// @Router /base/login [post]
func login(c *gin.Context) {
	var req request.LoginRequest
	_ = c.ShouldBindJSON(&req)
	if req.Username == "" || req.Password == "" {
		response.FailWithMessage("参数错误", c)
		return
	}
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
	}
	u, err := service.Login(user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service.GenerateTokenForUser(c, u)
}

// @Tags User
// @Summary 用户修改密码
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.ChangePassword true "用户修改密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/setPassword [PUT]
func setPassword(c *gin.Context) {
	var req request.ChangePassword
	_ = c.ShouldBindJSON(&req)
	if req.Username == "" || req.Password == "" || req.NewPassword == "" || req.Password == req.NewPassword {
		response.FailWithMessage("参数错误", c)
		return
	}
	if value, exists := c.Get("claims"); exists {
		if v, ok := value.(*request.UserClaims); ok {
			if v.Username != req.Username {
				response.FailWithMessage("请先登录", c)
				return
			}
		}
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password,
	}
	err := service.ChangePassword(user, req.NewPassword)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)

}

// @Tags Base
// @Summary 用户注册账号
// @Produce  application/json
// @Param data body request.RegisterRequest true "用户注册接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"注册成功"}"
// @Router /base/register [POST]
func register(c *gin.Context) {
	var req request.RegisterRequest

	_ = c.ShouldBindJSON(&req)
	global.AsongLogger.Info(req)
	if req.Username == "" || req.Password == "" || req.Nickname == "" {
		response.FailWithMessage("参数错误", c)
		return
	}
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
	}
	err := service.Register(user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("注册成功", c)

}
func sdasd() {

}
