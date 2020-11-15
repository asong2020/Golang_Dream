package main

import (
	"fmt"
)

// golang没有enum 使用const代替
const (
	TYPE_BALANCE      = 1 // type = 1
	TYPE_INTEGRAL     = 2 // type = 2
	TYPE_COUPON       = 3 // type = 3
	TYPE_FREEPOSTAGE  = 4 // type = 4
)

// 是否使用有优惠卷
func IsUseDiscount(discountType , value uint32) bool {
	return (value & (1<< (discountType-1))) > 0
}


// 设置使用
func SetDiscountValue(discountType ,value uint32) uint32{
	return value | (1 << (discountType-1))
}

func main()  {
	// 测试1 不设置优惠类型
	var flag1 uint32 = 0
	fmt.Println(IsUseDiscount(TYPE_BALANCE,flag1))
	fmt.Println(IsUseDiscount(TYPE_INTEGRAL,flag1))
	fmt.Println(IsUseDiscount(TYPE_COUPON,flag1))
	fmt.Println(IsUseDiscount(TYPE_FREEPOSTAGE,flag1))


	// 测试2 只设置一个优惠类型
	var flag2 uint32 = 0
	flag2 = SetDiscountValue(TYPE_BALANCE,flag2)
	fmt.Println(IsUseDiscount(TYPE_BALANCE,flag2))
	fmt.Println(IsUseDiscount(TYPE_INTEGRAL,flag2))
	fmt.Println(IsUseDiscount(TYPE_COUPON,flag2))
	fmt.Println(IsUseDiscount(TYPE_FREEPOSTAGE,flag2))

	// 测试3 设置两个优惠类型
	var flag3 uint32 = 0
	flag3 = SetDiscountValue(TYPE_BALANCE,flag3)
	flag3 = SetDiscountValue(TYPE_INTEGRAL,flag3)
	fmt.Println(IsUseDiscount(TYPE_BALANCE,flag3))
	fmt.Println(IsUseDiscount(TYPE_INTEGRAL,flag3))
	fmt.Println(IsUseDiscount(TYPE_COUPON,flag3))
	fmt.Println(IsUseDiscount(TYPE_FREEPOSTAGE,flag3))

	// 测试4 设置三个优惠类型
	var flag4 uint32 = 0
	flag4 = SetDiscountValue(TYPE_BALANCE,flag4)
	flag4 = SetDiscountValue(TYPE_INTEGRAL,flag4)
	flag4 = SetDiscountValue(TYPE_COUPON,flag4)
	fmt.Println(IsUseDiscount(TYPE_BALANCE,flag4))
	fmt.Println(IsUseDiscount(TYPE_INTEGRAL,flag4))
	fmt.Println(IsUseDiscount(TYPE_COUPON,flag4))
	fmt.Println(IsUseDiscount(TYPE_FREEPOSTAGE,flag4))

	// 测试5 设置四个优惠类型
	var flag5 uint32 = 0
	flag5 = SetDiscountValue(TYPE_BALANCE,flag5)
	flag5 = SetDiscountValue(TYPE_INTEGRAL,flag5)
	flag5 = SetDiscountValue(TYPE_COUPON,flag5)
	flag5 = SetDiscountValue(TYPE_FREEPOSTAGE,flag5)
	fmt.Println(IsUseDiscount(TYPE_BALANCE,flag5))
	fmt.Println(IsUseDiscount(TYPE_INTEGRAL,flag5))
	fmt.Println(IsUseDiscount(TYPE_COUPON,flag5))
	fmt.Println(IsUseDiscount(TYPE_FREEPOSTAGE,flag5))
}
