package controllers

import (
	"asong.cloud/Chatroom/models"
	"container/list"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"time"
)

/**
Author: Asong
create by on 2020/6/1
订阅号: Golang梦工厂
Features: 聊天室控制层
*/

var (
	// channel 对于新加入的用户
	subscribe = make(chan Subscriber,10)
	//channel 用户退出
	unsubscribe = make(chan UnSubscriber,10)
	//send message
	publish = make(chan models.Event,10)
	//订阅者列表 也就是聊天室列表
	subscribers = list.New()
)

type UnSubscriber struct {
	Name string
	Tel string
}

//订阅者
type Subscriber struct {
	Name string
	Tel string
	Conn *websocket.Conn
}

//join 假如聊天室
// 传入 用户名  websocket连接
func Join(user ,tel string,ws *websocket.Conn)  {
	subscribe <- Subscriber{Name:user,Tel:tel,Conn:ws}
}

// 用户退出
func Leave(user UnSubscriber)  {
	unsubscribe <- user
}

func newEvent(ep models.EventType,user,tel,msg string) models.Event {
	return models.Event{ep,user,tel,int(time.Now().Unix()),msg}
}

type Subscription struct {
	Archive []models.Event // 从models archive 获取事件
	New <- chan models.Event // 新的时间进入
}

func chatRoom()  {
	for{
		select {
		//新用户处理
		case sub := <-subscribe:
			if !isUserExist(subscribers, sub.Name) {
				subscribers.PushBack(sub)
				//加入一个事件
				publish <- newEvent(models.EVENT_JOIN, sub.Name,sub.Tel, "")
				RoomList(sub.Name)
				//添加日志
				beego.Info("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				beego.Info("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
			//事件处理
		case event := <-publish:
			broadcastWebSocket(event)
			me := new(models.Message)
			me.Content = event.Content
			us , err :=models.QueryByTel(event.Tel)
			if err!= nil{
				beego.Info("该用户不存在")
			}
			me.Tel = us.Tel
			me.Send = time.Now()
			err = models.AddMessage(me)
			if err != nil{
				beego.Info("add message error: ",err)
			}
			models.NewArchive(event)
			if event.Type == models.EVENT_MESSAGE {
				beego.Info("Message from", event.User, ";Content:", event.Content)
			}
		case unsub := <- unsubscribe:
			for sub := subscribers.Front(); sub!= nil;sub = sub.Next(){
				if sub.Value.(Subscriber).Name == unsub.Name{
					subscribers.Remove(sub)
					ws := sub.Value.(Subscriber).Conn
					if ws != nil{
						ws.Close()
						beego.Info("WebSocket closed:",unsub)
					}
					publish <- newEvent(models.EVENT_LEAVE,unsub.Name,unsub.Tel,"")
					break
				}
			}
		}
	}
}

func init()  {
	go chatRoom()
}

func isUserExist(subscribers *list.List,user string)  bool{
	for sub := subscribers.Front(); sub != nil ; sub = sub.Next(){
		if sub.Value.(Subscriber).Name == user{
			return true
		}
	}
	return false
}