package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Server struct {
	System System
	Mysql  Mysql `json:"mysql" yaml:"mysql"`
}

type Mysql struct {
	Host     string `json:"host" yaml:"host"`
	Port     int `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Db       string `json:"db" yaml:"db"`
	Conn     struct {
		MaxIdle int `json:"maxidle" yaml:"maxidle"`
		Maxopen int `json:"maxopen" yaml:"maxopen"`
	}
}

type System struct {
	Port int `json:"port" yaml:"port"`
}

func (m *Server) Load(filePath string) error {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("read yaml file error %v\n", err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, m)
	if err != nil {
		fmt.Printf("unmarshal config error %v\n", err)
		return err
	}
	return nil
}