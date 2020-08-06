package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func CofParse(file string,in interface{})  error{
	yamlFile,err:=ioutil.ReadFile(file)
	if err != nil{
		return err
	}
	err = yaml.Unmarshal(yamlFile,in)
	if err != nil{
		return err
	}
	return nil
}
