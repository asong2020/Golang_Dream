package main

import (
	"fmt"
)

func init()  {
	fmt.Println("package A exec c.go test 顺序")
}

func init()  {
	fmt.Println("package A exec c.go 111111")
}

func init()  {
	fmt.Println("package A exec c.go 22222")
}


func A_c()  string{
	return ""
}