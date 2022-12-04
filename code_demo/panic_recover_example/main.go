package main

import (
	"runtime"
)

func main()  {
 maybeGoexit()
}
func maybeGoexit() {
	defer func() {
		recover()
	}()
	//defer panic("cancelled Goexit!")
	runtime.Goexit()
}
func Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}