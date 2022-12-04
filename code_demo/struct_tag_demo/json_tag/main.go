package main

import (
	"encoding/json"
	"fmt"
)

type Location struct {
	Longitude float32 `json:"lon,omitempty"`
	Latitude  float32 `json:"lat,omitempty"`
}


func main()  {
	l := Location{
		Longitude: 10.12,
		Latitude: 1.33,
	}
	a, err := json.Marshal(l)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(string(a))
	res := Location{}
	err = json.Unmarshal(a, &res)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
