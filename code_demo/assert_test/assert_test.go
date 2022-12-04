package main

import (
"testing"
)

var dst int64

// 空接口类型直接类型断言具体的类型
func Benchmark_efaceToType(b *testing.B) {
	b.Run("efaceToType", func(b *testing.B) {
		var ebread interface{} = int64(666)
		for i := 0; i < b.N; i++ {
			dst = ebread.(int64)
		}
	})
}

// 空接口类型使用TypeSwitch 只有部分类型
func Benchmark_efaceWithSwitchOnlyIntType(b *testing.B) {
	b.Run("efaceWithSwitchOnlyIntType", func(b *testing.B) {
		var ebread interface{} = 666
		for i := 0; i < b.N; i++ {
			OnlyInt(ebread)
		}
	})
}

// 空接口类型使用TypeSwitch 所有类型
func Benchmark_efaceWithSwitchAllType(b *testing.B) {
	b.Run("efaceWithSwitchAllType", func(b *testing.B) {
		var ebread interface{} = 666
		for i := 0; i < b.N; i++ {
			Any(ebread)
		}
	})
}

//直接使用类型转换
func Benchmark_TypeConversion(b *testing.B) {
	b.Run("typeConversion", func(b *testing.B) {
		var ebread int32 = 666

		for i := 0; i < b.N; i++ {
			dst = int64(ebread)
		}
	})
}

// 非空接口类型判断一个类型是否实现了该接口 两个方法
func Benchmark_ifaceToType(b *testing.B) {
	b.Run("ifaceToType", func(b *testing.B) {
		var iface Basic = &User{}
		for i := 0; i < b.N; i++ {
			iface.GetName()
			iface.SetName("1")
		}
	})
}

// 非空接口类型判断一个类型是否实现了该接口 12个方法
func Benchmark_ifaceToTypeWithMoreMethod(b *testing.B) {
	b.Run("ifaceToTypeWithMoreMethod", func(b *testing.B) {
		var iface MoreMethod = &More{}
		for i := 0; i < b.N; i++ {
			iface.Get()
			iface.Set()
			iface.One()
			iface.Two()
			iface.Three()
			iface.Four()
			iface.Five()
			iface.Six()
			iface.Seven()
			iface.Eight()
			iface.Nine()
			iface.Ten()
		}
	})
}

// 直接调用方法
func Benchmark_DirectlyUseMethod(b *testing.B) {
	b.Run("directlyUseMethod", func(b *testing.B) {
		m := &More{
			Name: "asong",
		}
		m.Get()
	})
}

func OnlyInt(val interface{}) {
	switch val.(type) {
	case int:
	case []int:
	case int64:
	case []int64:
	case int32:
	case []int32:
	case int16:
	case []int16:
	case int8:
	case []int8:
	}
}

func Any(value interface{}) {
	switch value.(type) {
	case bool:
	case []bool:
	case complex128:
	case []complex128:
	case complex64:
	case []complex64:
	case float64:
	case []float64:
	case float32:
	case []float32:
	case int:
	case []int:
	case int64:
	case []int64:
	case int32:
	case []int32:
	case int16:
	case []int16:
	case int8:
	case []int8:
	case string:
	case []string:
	case uint:
	case []uint:
	case uint64:
	case []uint64:
	case uint32:
	case []uint32:
	case uint16:
	case []uint16:
	case uint8:
	case []byte:
	case uintptr:
	case []uintptr:
	case error:
	case []error:
	default:
	}
}
