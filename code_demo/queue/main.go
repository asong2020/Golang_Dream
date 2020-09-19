package main

import (
	"fmt"
	"time"

	"Golang_Dream/code_demo/queue/mq"
)

var (
	topic = "Golang梦工厂"
)

func main()  {
	//OnceTopic()
	//ManyTopic()
}

// 一个topic 测试
func OnceTopic()  {
	m := mq.NewClient()
	m.SetConditions(10)
	ch,err :=m.Subscribe(topic)
	if err != nil{
		fmt.Println("subscribe failed")
		return
	}
	go OncePub(m)
	OnceSub(ch,m)
	defer m.Close()
}

// 定时推送
func OncePub(c *mq.Client)  {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for  {
		select {
		case <- t.C:
			err := c.Publish(topic,"asong真帅")
			if err != nil{
				fmt.Println("pub message failed")
			}
		default:

		}
	}
}

// 接受订阅消息
func OnceSub(m <-chan interface{},c *mq.Client)  {
	for  {
		val := c.GetPayLoad(m)
		fmt.Printf("get message is %s\n",val)
	}
}

//多个topic测试
func ManyTopic()  {
	m := mq.NewClient()
	defer m.Close()
	m.SetConditions(10)
	top := ""
	for i:=0;i<10;i++{
		top = fmt.Sprintf("Golang梦工厂_%02d",i)
		go Sub(m,top)
	}
	ManyPub(m)
}

func ManyPub(c *mq.Client)  {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for  {
		select {
		case <- t.C:
			for i:= 0;i<10;i++{
				//多个topic 推送不同的消息
				top := fmt.Sprintf("Golang梦工厂_%02d",i)
				payload := fmt.Sprintf("asong真帅_%02d",i)
				err := c.Publish(top,payload)
				if err != nil{
					fmt.Println("pub message failed")
				}
			}
		default:

		}
	}
}

func Sub(c *mq.Client,top string)  {
	ch,err := c.Subscribe(top)
	if err != nil{
		fmt.Printf("sub top:%s failed\n",top)
	}
	for  {
		val := c.GetPayLoad(ch)
		if val != nil{
			fmt.Printf("%s get message is %s\n",top,val)
		}
	}
}