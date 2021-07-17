package main

import (
	"testing"
)

func BenchmarkMaxNoinline(b *testing.B) {
	x := []int{1,2,3,4}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		MaxNoinline(x)
	}
}

func BenchmarkMaxInline(b *testing.B) {
	x := []int{1,2,3,4}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		MaxInline(x)
	}
}