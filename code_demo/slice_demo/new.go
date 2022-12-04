package main

import (
	"fmt"
	"reflect"
)

type Errno struct {
	code int64
	info string
}

func (e Errno) Error() string {
	panic("implement me")
}

func main(){
	err := res()
	t := reflect.TypeOf(err)
	v := reflect.ValueOf(err)
	num := t.NumField()
	for i:=0;i<num;i++{
		f := v.Field(i)
		fmt.Println(f)
	}
}

func res() error {
	return Errno{
		code: 1,
		info: "test",
	}
}


//	t := new(test)
//	*t.A = 10
//	fmt.Println(t.A)
//
//
//
//	// 数组
//	array := new([5]int64)
//	fmt.Printf("array: %p %#v \n", &array, array)// array: 0xc0000ae018 &[5]int64{0, 0, 0, 0, 0}
//	(*array)[0] = 1
//	fmt.Printf("array: %p %#v \n", &array, array)// array: 0xc0000ae018 &[5]int64{1, 0, 0, 0, 0}
//
//	// 切片
//	slice := new([]int64)
//	fmt.Printf("slice: %p %#v \n", &slice, slice) // slice: 0xc0000ae028 &[]int64(nil)
//	(*slice)[0] = 1
//	fmt.Printf("slice: %p %#v \n", &slice, slice) // panic: runtime error: index out of safe [0] with length 0
//
//	// map
//	map1 := new(map[string]string)
//	fmt.Printf("map1: %p %#v \n", &map1, map1) // map1: 0xc00000e038 &map[string]string(nil)
//	(*map1)["key"] = "value"
//	fmt.Printf("map1: %p %#v \n", &map1, map1) // panic: assignment to entry in nil map
//
//	// channel
//	channel := new(chan string)
//	fmt.Printf("channel: %p %#v \n", &channel, channel) // channel: 0xc0000ae028 (*chan string)(0xc0000ae030)
//	//channel <- "123" // Invalid operation: channel <- "123" (send to non-chan type *chan string)
//}