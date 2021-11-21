package main

import (
	"testing"
)

func BenchmarkSumString(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		SumString(base)
	}
	b.StopTimer()
}

func BenchmarkSprintfString(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		SprintfString(base)
	}
	b.StopTimer()
}

func BenchmarkBuilderString(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		BuilderString(base)
	}
	b.StopTimer()
}

func BenchmarkBytesBuffString(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		bytesString(base)
	}
	b.StopTimer()
}

func BenchmarkJoinstring(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		Joinstring()
	}
	b.StopTimer()
}

func BenchmarkByteSliceString(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		byteSliceString(base)
	}
	b.StopTimer()
}


