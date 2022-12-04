package main

import (
	"math"
	"time"
)

func main()  {
	for i := 0; i < math.MaxInt64; i++ {
		go func(i int) {
			time.Sleep(5 * time.Second)
		}(i)
	}
}
