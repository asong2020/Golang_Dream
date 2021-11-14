package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string `asong:"Username"`
	Age  uint16 `asong:"age"`
	Password string `asong:"min=6,max=10"`
}

func getTag(u User) {
	t := reflect.TypeOf(u)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("asong")
		fmt.Println("get tag is ", tag)
	}
}



func main()  {
	u := User{
		Name: "asong",
		Age: 5,
		Password: "123456",
	}
	getTag(u)
}