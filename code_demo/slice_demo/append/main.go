package main

import (
	"fmt"
)

func appendSlice(s []string)  {
	s = append(s, "快关注！！")
	fmt.Println("out slice: ", s)
}

func modifySlice(s []string)  {
	s[0] = "song"
	s[1] = "Golang"
	fmt.Println("out slice: ", s)
}

func main()  {
	s := []string{"asong", "Golang梦工厂"}
	appendSlice(s)
	//modifySlice(s)
	fmt.Println("inner slice: ", s)
}
