package models

import "container/list"

/**
Author: Asong
create by on 2020/6/1
订阅号: Golang梦工厂
Features: 实践模型层
*/

type EventType int

//事件类型
const  (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
	EVENT_LIST
)

type Event struct {
	Type EventType // JOIN LEAVE MESSAGE
	User string
	Tel string // 电话
	Timestamp int//UNIX timestamp
	Content string //内容
}
// 存放EVENT 最大数量
const archiveSize = 20
// 创建一个链表存放EVENT
var archive = list.New()

//存放新的EVENT到链表
func NewArchive(event Event)  {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

//获取事件
func GetEvents(lastReceived int) []Event {
	events := make([]Event,0,archive.Len())
	for event := archive.Front();event != nil; event = event.Next(){
		e := event.Value.(Event)
		if e.Timestamp > int(lastReceived){
			events = append(events,e)
		}
	}
	return events
}