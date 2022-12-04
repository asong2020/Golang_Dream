package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/model"
	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/service"
)

type UserHandler struct {
	engine  *gin.Engine
	service *service.UserService
}

func NewUserHandler(engine *gin.Engine, userService *service.UserService) *UserHandler {
	return &UserHandler{
		engine:  engine,
		service: userService,
	}
}

func (h *UserHandler) Run() {
	// Force log's color
	gin.ForceConsoleColor()
	h.engine.Use(gin.Logger())
	h.engine.Use(gin.Recovery())
	h.registerRouter()

	err := h.engine.Run()
	if err != nil {
		log.Fatalln("server start failed")
	}
}

func (h *UserHandler) registerRouter() {
	u := h.engine.Group("api/user")
	{
		u.POST("/create", h.Create)
		u.PUT("/update", h.Update)
		u.DELETE("/delete", h.Delete)
		u.GET("/info", h.MGet)
		u.POST("/search", h.Search)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	users := make([]*model.UserEs, 0)
	user := model.UserEs{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "Invalid argument"})
		return
	}
	users = append(users, &user)
	if err := h.service.BatchAdd(c, users); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *UserHandler) Update(c *gin.Context) {
	users := make([]*model.UserEs, 0)
	user := model.UserEs{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "Invalid argument"})
		return
	}
	users = append(users, &user)
	if err := h.service.BatchUpdate(c, users); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	users := make([]*model.UserEs, 0)
	user := model.UserEs{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "Invalid argument"})
		return
	}
	users = append(users, &user)
	if err := h.service.BatchDel(c, users); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *UserHandler) MGet(c *gin.Context) {
	ids := c.Query("id")
	IDS := make([]uint64, 0)
	for _, id := range strings.Split(ids, ",") {
		d, _ := strconv.Atoi(id)
		IDS = append(IDS, uint64(d))
	}

	res, err := h.service.MGet(c, IDS)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": res,
	})
}

func (h *UserHandler) Search(c *gin.Context) {
	req := model.SearchRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "Invalid argument"})
		return
	}
	res, err := h.service.Search(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": res,
	})
}
