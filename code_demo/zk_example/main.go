package main

import (
	"asong.cloud/Golang_Dream/code_demo/zk_example/zksvr"
)

func main()  {

	zkList := []string{"127.0.0.1:2181"}
	conn := zksvr.GetConnect(zkList)
	defer conn.Close()
	zksvr.AddNode(conn)
	println(conn.Server())
}
