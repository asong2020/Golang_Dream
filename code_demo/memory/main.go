package main

import (
	"fmt"
	"unsafe"
)

// 64位平台，对齐参数是8
type User struct {
	A int32 // 4
	B []int32 // 24
	C string // 16
	D bool // 1
}



func main()  {
	var u User
	fmt.Println("u size is ",unsafe.Sizeof(u))
}
