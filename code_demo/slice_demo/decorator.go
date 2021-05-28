package main

import (
	"fmt"
)

func main(){
	r := Log(DoSomething)
	r()
}

func Log(f func()) func() {
	wrapper := func() {
		fmt.Println("begin log")
		res := f
		res()
		fmt.Println("end log")
	}
	return wrapper
}


func DoSomething()  {
	fmt.Println("to do")
}