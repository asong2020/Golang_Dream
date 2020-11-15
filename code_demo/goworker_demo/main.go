package main

import (
	"fmt"
)

//
//func main() {
//	array := []int{}
//	array = append(array, 7,8,9)
//	fmt.Printf("len: %d cap:%d data:%+v\n", len(array), cap(array), array)
//	MyAppend(array)
//	fmt.Printf("len: %d cap:%d data:%+v\n", len(array), cap(array), array)
//	p := unsafe.Pointer(&array[2])
//	q := uintptr(p)+8
//	t := (*int)(unsafe.Pointer(q))
//	fmt.Println(*t)
//}
//
//func MyAppend(array []int) {
//	array = append(array, 10)
//}

//func main() {
//	array := []int{7,8,9}
//	fmt.Printf("main ap brfore: len: %d cap:%d data:%+v\n", len(array), cap(array), array)
//	ap(array)
//	fmt.Printf("main ap after: len: %d cap:%d data:%+v\n", len(array), cap(array), array)
//}
//
//func ap(array []int) {
//	fmt.Printf("ap brfore:  len: %d cap:%d data:%+v\n", len(array), cap(array), array)
//	array = append(array, 10)
//	fmt.Printf("ap after:   len: %d cap:%d data:%+v\n", len(array), cap(array), array)
//}

func main()  {
	var args int64= 1
	addr := &args
	fmt.Printf("原始指针的内存地址是 %p\n", addr)
	fmt.Printf("指针变量addr存放的地址 %p\n", &addr)
	modifiedNumber(addr) // args就是实际参数
	fmt.Printf("改动后的值是  %d\n",args)
}

func modifiedNumber(addr *int64)  { //这里定义的args就是形式参数
	fmt.Printf("形参地址 %p \n",&addr)
	*addr = 10
}