package main

import (
	"fmt"
)

func main()  {
	var buf []byte
	fmt.Println(cap(buf))
	var bufMake []byte = make([]byte,5)
	fmt.Println(cap(bufMake))
}
