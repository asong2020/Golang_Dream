package controllers

import (
	"asong.cloud/Chatroom/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"net/http"
)

/**
Author: Asong
create by on 2020/6/1
订阅号: Golang梦工厂
Features: WebSocket
*/

type WebSocketControllers struct {
	baseController
}

func (this *WebSocketControllers)Get()  {
	uname := this.GetString("name")
	tel := this.GetString("tel")
	//如果用户名为空回到首页
	if len(uname) == 0 || len(tel) == 0{
		this.Redirect("/",302)
		return
	}

	this.TplName = "websocket.html"
	this.Data["IsWebSocket"] = true
	this.Data["UserName"] = uname
	this.Data["UserTel"] = tel
}

func (this *WebSocketControllers)Join()  {
	uname := this.GetString("name")
	tel := this.GetString("tel")
	if len(uname)==0 || len(tel) == 0{
		this.Redirect("/",302)
		return
	}
	ws ,err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _,ok := err.(websocket.HandshakeError); ok{
		http.Error(this.Ctx.ResponseWriter,"Not a websocket handshake",400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection",err)
		return
	}

	//join 加入聊天室
	Join(uname,tel,ws)
	unsub := UnSubscriber{
		Name: uname,
		Tel:  tel,
	}
	//用户退出
	defer Leave(unsub)
	//循环读取消息
	for{
		_,message,err := ws.ReadMessage()
		if err != nil{
			return
		}
		publish <- newEvent(models.EVENT_MESSAGE,uname,tel,string(message))
	}
}

func broadcastWebSocket(event models.Event)  {
	data ,err := json.Marshal(event)
	if err != nil{
		beego.Error("Fail to marshal event",err)
		return
	}
	//对聊天室列表中的每个用户进行消息推送
	for sub := subscribers.Front(); sub != nil ; sub = sub.Next(){
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage,data) != nil{
				unsub := UnSubscriber{
					Name: sub.Value.(Subscriber).Name,
					Tel:  sub.Value.(Subscriber).Tel,
				}
				//user 断开连接
				unsubscribe <- unsub
			}
		}
	}
}

func RoomList(name string)  {
	var unameList string = ""
	var ws *websocket.Conn
	for sub := subscribers.Front(); sub != nil ; sub = sub.Next(){
		unameList += sub.Value.(Subscriber).Name + ";"
		if sub.Value.(Subscriber).Name == name{
			ws = sub.Value.(Subscriber).Conn
		}
	}
	event := models.Event{
		Type:      models.EVENT_LIST,
		User:      unameList,
		Tel:       "",
		Timestamp: 0,
		Content:   "",
	}
	data ,err := json.Marshal(event)
	if err !=nil{
		beego.Error("Fail to marshal event",err)
		return
	}
	if ws != nil{
		ws.WriteMessage(websocket.TextMessage,data)
	}
}