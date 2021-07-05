package main

import (
	"fmt"
)

func main()  {
	sl := make([]byte,0,2)
	sl = append(sl, 'A')
	sl = append(sl,'B')
	fmt.Println(sl)
}
