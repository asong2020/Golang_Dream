package main

import (
	"testing"
)

//go:noinline
func AddNoinline(x,y,z int) int {
	return x+y+z
}

func AddInline(x,y,z int) int {
	return x+y+z
}
func BenchmarkAddNoinline(b *testing.B) {
	x,y,z :=1,2,3
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		AddNoinline(x,y,z)
	}
}

func BenchmarkAddInline(b *testing.B) {
	x,y,z :=1,2,3
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		AddInline(x,y,z)
	}
}