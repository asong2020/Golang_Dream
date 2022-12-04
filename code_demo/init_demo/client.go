package main

import (
	"fmt"
	"net"
	"time"
)

func main()  {
	c,err := net.Dial("tcp","127.0.0.1:8888")
	if err != nil{
		panic("1111")
	}
	for  {
		for i:=0;i<3;i++{
			time.Sleep(time.Second * 1)
			c.Write([]byte("SUB"))
		}
		n := time.Now()
		b := make([]byte,20)
		nbytes,err := c.Read(b)
		end := time.Now()
		d := end.Sub(n)
		fmt.Println(d.Seconds())
		if err != nil{
			fmt.Println(err.Error())
			c.Close()
			return
		}
		fmt.Println(nbytes)
	}
}
