package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"asong.cloud/Golang_Dream/wire_cron_example/config"
	"asong.cloud/Golang_Dream/wire_cron_example/wire"
)

func init()  {
	file,err := ioutil.ReadFile("./conf.yaml")
	if err != nil{
		log.Fatalln("read yaml file error")
	}
	conf := config.Server{}
	err = yaml.Unmarshal(file,&conf)
	if err != nil{
		log.Fatalln("Unmarshal file error")
	}
	cron := wire.InitializeCron(&conf.Mysql)
	err = cron.Start()
	if err != nil{
		log.Fatalln(err.Error())
	}
}

func main()  {
	select {

	}
}
