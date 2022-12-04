package main

import (
	"context"
	"time"
)

func main()  {
	m := make(map[string]string)
	m ["asong"] = "Golang梦工厂"
	ctx := context.WithValue(context.Background(), "asong", m)
	go func() {
		for {
			m1 := ctx.Value("asong")
			mm := m1.(map[string]string)
			mm["asong"] = "123213"
		}
	}()
	go func() {
		for {
			m1 := ctx.Value("asong")
			mm := m1.(map[string]string)
			mm["asong"] = "123213"
		}
	}()
	time.Sleep(10 * time.Second)
}
