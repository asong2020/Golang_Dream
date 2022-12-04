package main

import (
	"reflect"
	"unsafe"
)





func stringtoslicebytetmp(s string) []byte {
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	ret := reflect.SliceHeader{Data: str.Data, Len: str.Len, Cap: str.Len}
	return *(*[]byte)(unsafe.Pointer(&ret))
}

func main()  {
	a := "hello"
	b := stringtoslicebytetmp(a)
	b[0] = 'H'
}
