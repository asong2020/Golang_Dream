package main

import (
	"fmt"
	"unsafe"
)

type demo1 struct {
	a struct{}
	b int32
}

type demo2 struct {
	a int32
	b struct{}
}

func main()  {
	fmt.Println(unsafe.Sizeof(demo2{}))
}