package main

import (
	"fmt"
)

func main()  {
	TestMulti(1,3,4,5,6)
}

func TestMulti(a int,b ...uint64)  {
	for i := range b{
		fmt.Println(i)
	}
}
