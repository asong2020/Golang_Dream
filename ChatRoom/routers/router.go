package routers

import (
	"asong.cloud/Chatroom/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.AppController{})
    beego.Router("/join",&controllers.AppController{},"post:Join")

	beego.Router("/ws",&controllers.WebSocketControllers{})
    beego.Router("ws/join",&controllers.WebSocketControllers{},"get:Join")

    beego.Router("/register",&controllers.RegisterControllers{})
    beego.Router("/register/join",&controllers.RegisterControllers{},"post:Join")
}
