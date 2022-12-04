package main

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

func AsyncAdd(run func() error)  {
	//TODO: 扔进异步协程池
	go run()
}

func GetInstance(ctx context.Context,id uint64) (string, error) {
	data,err := GetFromRedis(ctx,id)
	if err != nil && err != redis.Nil{
		return "", err
	}
	// 没有找到数据
	if err == redis.Nil {
		data,err = GetFromDB(ctx,id)
		if err != nil{
			return "", err
		}
		AsyncAdd(func() error{
			ctxAsync,cancel := context.WithTimeout(context.Background(),3 * time.Second)
			defer cancel()
			return UpdateCache(ctxAsync,id,data)
		})
	}
	return data,nil
}

func GetFromRedis(ctx context.Context,id uint64) (string,error) {
	// TODO: 从redis获取信息
	return "",nil
}

func GetFromDB(ctx context.Context,id uint64) (string,error) {
	// TODO: 从DB中获取信息
	return "",nil
}

func UpdateCache(ctx context.Context,id interface{},data string) error {
	// TODO：更新缓存信息
	return nil
}

func main()  {
	ctx,cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	_,err := GetInstance(ctx,2021)
	if err != nil{
		return
	}
}
