package main

import "fmt"

type Photos struct {
	width uint64
	height uint64
	value string
}

type OrderChangeNotificationHandler interface {
	GenerateMessage() string
	GeneratePhotos() Photos
	generateUrl() string
}


type OrderChangeNotificationHandlerImpl struct {
	url string
}

func NewOrderChangeNotificationHandlerImpl() OrderChangeNotificationHandler {
	return OrderChangeNotificationHandlerImpl{
		url: "https://base.test.com",
	}
}

func (o OrderChangeNotificationHandlerImpl) GenerateMessage() string {
	return "OrderChangeNotificationHandlerImpl GenerateMessage"
}

func (o OrderChangeNotificationHandlerImpl) GeneratePhotos() Photos {
	return Photos{
		width: 1,
		height: 1,
		value: "https://www.baidu.com",
	}
}

func (w OrderChangeNotificationHandlerImpl) generateUrl() string {
	return w.url
}

type WebOrderChangeNotificationHandler struct {
	OrderChangeNotificationHandler
	url string
}

func (w WebOrderChangeNotificationHandler) generateUrl() string {
	return w.url
}

type AppOrderChangeNotificationHandler struct {
	OrderChangeNotificationHandler
	url string
}

func (a AppOrderChangeNotificationHandler) generateUrl() string {
	return a.url
}

func check(handler OrderChangeNotificationHandler)  {
	fmt.Println(handler.GenerateMessage())
}

func main()  {
	base := NewOrderChangeNotificationHandlerImpl()
	web := WebOrderChangeNotificationHandler{
		OrderChangeNotificationHandler: base,
		url: "http://web.test.com",
	}
	fmt.Println(web.GenerateMessage())
	fmt.Println(web.generateUrl())

	check(web)
}
