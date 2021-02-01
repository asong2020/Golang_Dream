package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"testing"
)

func BenchmarkMemoryStack(b *testing.B) {
	var s S

	f, err := os.Create("stack.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		s = byCopy()
	}

	trace.Stop()

	b.StopTimer()

	_ = fmt.Sprintf("%v", s.a)
}

func BenchmarkMemoryHeap(b *testing.B) {
	var s *S

	f, err := os.Create("heap.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		s = byPointer()
	}

	trace.Stop()

	b.StopTimer()

	_ = fmt.Sprintf("%v", s.a)
}

type S struct {
	a, b, c int64
	d, e, f string
	g, h, i float64
}

func byCopy() S {
	return S{
		a: 1, b: 1, c: 1,
		e: "foo", f: "foo",
		g: 1.0, h: 1.0, i: 1.0,
	}
}

func byPointer() *S {
	return &S{
		a: 1, b: 1, c: 1,
		e: "foo", f: "foo",
		g: 1.0, h: 1.0, i: 1.0,
	}
}