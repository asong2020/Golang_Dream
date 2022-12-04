package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main()  {
	http.HandleFunc("/", SayHello) // 设置访问的路由

	log.Fatalln(http.ListenAndServe(":8080",nil))
}

func SayHello(writer http.ResponseWriter, request *http.Request)  {
	fmt.Println(&request)

	go func() {
		for range time.Tick(time.Second) {
			select {
			case <- request.Context().Done():
				fmt.Println("request is outgoing")
				return
			default:
				fmt.Println("Current request is in progress")
			}
		}
	}()

	time.Sleep(2 * time.Second)
	writer.Write([]byte("Hi"))
}