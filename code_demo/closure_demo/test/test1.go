package main

import (
	"fmt"
)

func makeAverager() func(val float32) float32{
	series := make([]float32,0)
	return func(val float32) float32 {
		series = append(series, val)
		total := float32(0)
		for _,v:=range series{
			total +=v
		}
		return total/ float32(len(series))
	}
}

func main() {
	avg := makeAverager()
	fmt.Println(avg(10))
	fmt.Println(avg(30))
}
