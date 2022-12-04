package main

import (
	"testing"
)

func BenchmarkSumString(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		SumString()
	}
	b.StopTimer()
}

func BenchmarkSprintfString(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		SprintfString()
	}
	b.StopTimer()
}

func BenchmarkBuilderString(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		BuilderString()
	}
	b.StopTimer()
}

func BenchmarkBytesBufferString(b *testing.B) {
	b.ResetTimer()
	for i:=0; i < b.N; i++{
		bytesString()
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
		byteSliceString()
	}
	b.StopTimer()
}

