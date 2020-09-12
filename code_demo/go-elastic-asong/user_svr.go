package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/config"
	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/wire"
)

type UserSvr struct {
	conf config.ServerConfig
}


func (u *UserSvr)init()  {
	file,err := ioutil.ReadFile("./config.yaml")
	if err != nil{
		fmt.Println("read yaml file failed")
	}

	err = yaml.UnmarshalStrict(file,&u.conf)
	if err != nil{
		fmt.Println("yaml unmarshal failed")
	}
}

func (u *UserSvr)Run()  {
	handler := wire.InitializeHandler(&u.conf)
	handler.Run()
}